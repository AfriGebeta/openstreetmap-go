package v1

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"openstreetmap-go/src/api/utils"
	apiValues "openstreetmap-go/src/api/values"
	gormModel "openstreetmap-go/src/db/generated/gorm"
	repository4 "openstreetmap-go/src/modules/apisize/repository"
	repository3 "openstreetmap-go/src/modules/changeset/repository"
	"openstreetmap-go/src/modules/oauth2/repository"
	repository2 "openstreetmap-go/src/modules/user_blocks/repository"
	"strconv"
	"strings"
	"time"
)

func create(w http.ResponseWriter, r *http.Request) error {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	var osm Osm
	err = xml.Unmarshal(body, &osm)
	if err != nil {
		return err
	}

	if len(osm.Changeset.Tags) == 0 {
		return fmt.Errorf("missing <changeset> or no <tag> elements")
	}

	authHeader := r.Header.Get("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return fmt.Errorf("Authorization required")
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	oat, err := repository.GetOAuthAccessToken(token)
	if err != nil {
		return err
	}

	if oat.ResourceOwnerId != nil {
		return errors.New("invalid authorization")
	}

	block, err := repository2.GetBlockedUser(*oat.ResourceOwnerId, time.Now())
	if err != nil {
		return err
	}
	if block.CreatedAt == nil {
		return errors.New("invalid authorization")
	}

	changeset := gormModel.Changesets{
		UserId:               *oat.ResourceOwnerId,
		CreatedAt:            time.Time{},
		ClosedAt:             time.Now().Add(time.Hour),
		NumChanges:           0,
		NumCreatedNodes:      0,
		NumModifiedNodes:     0,
		NumDeletedNodes:      0,
		NumCreatedWays:       0,
		NumModifiedWays:      0,
		NumDeletedWays:       0,
		NumCreatedRelations:  0,
		NumModifiedRelations: 0,
		NumDeletedRelations:  0,
	}

	changesetTag := []gormModel.ChangesetTags{}
	for _, tag := range osm.Changeset.Tags {
		changesetTag = append(changesetTag, gormModel.ChangesetTags{
			K: tag.Key,
			V: tag.Value,
		})
	}

	changesetSubscribers := gormModel.ChangesetsSubscribers{
		SubscriberId: *oat.ResourceOwnerId,
	}

	changeset, err = repository3.InsertIntoChangeset(&changeset, changesetTag, changesetSubscribers)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	_, err = w.Write([]byte(fmt.Sprint(changeset.ID)))
	if err != nil {
		return err
	}
	return nil
}

func closeChangeset(w http.ResponseWriter, r *http.Request) error {
	changesetId := r.PathValue("id")
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return fmt.Errorf("Authorization required")
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")

	oat, err := repository.GetOAuthAccessToken(token)
	if err != nil {
		return err
	}

	if oat.ResourceOwnerId != nil {
		return errors.New("invalid authorization")
	}

	block, err := repository2.GetBlockedUser(*oat.ResourceOwnerId, time.Now())
	if err != nil {
		return err
	}
	if block.CreatedAt == nil {
		return errors.New("invalid authorization")
	}

	changeSetIdInt, err := strconv.ParseInt(changesetId, 10, 64)
	if err != nil {
		return err
	}

	changeSet, err := repository3.GetChangeSetById(changeSetIdInt)
	if err != nil {
		return err
	}
	if changeSet.UserId != *oat.ResourceOwnerId {
		return errors.New("invalid authorization")
	}

	err = repository3.CloseChangeSetById(changeSetIdInt)
	if err != nil {
		return err
	}

	if err := utils.EncodeBody(w, "Changeset closed successfully", http.StatusOK); err != nil {
		return apiValues.ApiErrorInternalServerError
	}

	return nil
}
func upload(w http.ResponseWriter, r *http.Request) error {
	changesetId := r.PathValue("id")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	var osm OsmChange
	err = xml.Unmarshal(body, &osm)
	if err != nil {
		return err
	}

	authHeader := r.Header.Get("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return fmt.Errorf("Authorization required")
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")

	oat, err := repository.GetOAuthAccessToken(token)
	if err != nil {
		return err
	}

	if oat.ResourceOwnerId != nil {
		return errors.New("invalid authorization")
	}

	block, err := repository2.GetBlockedUser(*oat.ResourceOwnerId, time.Now())
	if err != nil {
		return err
	}
	if block.CreatedAt == nil {
		return errors.New("invalid authorization")
	}

	//roleResults, err := respository.GetUserRoleByUserId(*oat.ResourceOwnerId)
	//if err != nil {
	//	return err
	//}
	// GetChangeSetById
	changeSetIdint, err := strconv.ParseInt(changesetId, 10, 64)
	if err != nil {
		return err
	}
	changeSet, err := repository3.GetChangeSetById(changeSetIdint)
	if err != nil {
		return err
	}
	if changeSet.UserId == *oat.ResourceOwnerId {
		return errors.New("invalid authorization")
	}

	if !changeSet.ClosedAt.IsZero() && changeSet.ClosedAt.Before(time.Now()) {
		return errors.New("changeset is already closed")
	}

	size, err := repository4.GetApiSizeLimitByUserId(*oat.ResourceOwnerId)
	if err != nil {
		return err
	}
	if size == 0 {
		return errors.New("invalid authorization")
	}

	changesetIdInt64, err := strconv.ParseInt(changesetId, 10, 64)
	if err != nil {
		return err
	}

	var uploadResult UploadChangesetResult
	err = ProcessCreateElement(changesetIdInt64, osm, changeSet, *oat.ResourceOwnerId, osm.Create, &uploadResult)
	if err != nil {
		return err
	}

	err = ProcessModifyElement(changesetIdInt64, osm.Modify, &uploadResult)
	if err != nil {
		return err
	}

	err = ProcessDeleteElement(changesetIdInt64, osm.Modify, &uploadResult)
	if err != nil {
		return err
	}

	minLat := uploadResult.Node[0].Lat
	minLon := uploadResult.Node[0].Lon
	maxLat := uploadResult.Node[0].Lat
	maxLon := uploadResult.Node[0].Lon
	createdCount := 0

	for _, node := range uploadResult.Node {
		if node.Action == "create" {
			createdCount++
		}
		if node.Lat > maxLat {
			maxLat = node.Lat
		}

		if node.Lon > maxLon {
			maxLon = node.Lon
		}

		if node.Lat < minLat {
			minLat = node.Lat
		}

		if node.Lon < minLon {
			minLon = node.Lon
		}
	}

	// close the changeset
	data := make(map[interface{}]interface{})
	data["min_lat"] = minLat
	data["max_lat"] = maxLat
	data["min_lon"] = minLon
	data["max_lon"] = maxLon
	data["num_changes"] = "num_changes + " + strconv.Itoa(createdCount)
	data["num_created_nodes"] = "num_created_nodes + " + strconv.Itoa(len(uploadResult.Node))
	data["closed_at"] = time.Now().UTC()

	err = repository3.UpdateChangeset(changesetId, data)
	if err != nil {
		return err
	}

	result, err := buildDiffResultXml(uploadResult)
	if err != nil {
		return err
	}

	if err := utils.EncodeBody(w, result, http.StatusOK); err != nil {
		return apiValues.ApiErrorInternalServerError
	}
	return nil
}
