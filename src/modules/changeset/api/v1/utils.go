package v1

import (
	"encoding/xml"
	"errors"
	gormModel "openstreetmap-go/src/db/generated/gorm"
	repository3 "openstreetmap-go/src/modules/changeset/repository"
	"openstreetmap-go/src/modules/current_node_tags/repository"
	"openstreetmap-go/src/modules/current_nodes/respository"
	"openstreetmap-go/src/utils/geo"
	"strconv"
	"time"
)

func buildDiffResultXml(results UploadChangesetResult) (string, error) {
	diffResult := DiffResult{
		Version:     "0.6",
		Generator:   "OpenStreetMap server",
		Copyright:   "OpenStreetMap and contributors",
		Attribution: "http://www.openstreetmap.org/copyright",
		License:     "http://opendatacommons.org/licenses/odbl/1-0/",
	}

	for _, node := range results.Node {
		// Validate new_id and version
		if node.Id <= 0 {
			continue
		}

		nodeDiff := NodeDiff{
			NewID:      int(node.Id),
			NewVersion: 1,
		}

		//nodeDiff.OldID = utils.PtrOf(node.Id)

		if node.Version > 0 {
			nodeDiff.NewVersion = node.Version
		}

		diffResult.Nodes = append(diffResult.Nodes, nodeDiff)
	}

	output, err := xml.MarshalIndent(diffResult, "", "  ")
	if err != nil {
		return "", err
	}

	xmlHeader := `<?xml version="1.0" encoding="UTF-8"?>` + "\n"
	return xmlHeader + string(output), nil
}

func ProcessDeleteElement(changesetId int64, element Elements, uploadResult *UploadChangesetResult) error {
	if len(element.Nodes) != 0 {
		return nil
	}

	if len(element.Ways) != 0 {
		return nil
	}

	if len(element.Relations) != 0 {
		return nil
	}

	return nil
}
func ProcessModifyElement(changesetId int64, element Elements, uploadResult *UploadChangesetResult) error {
	if len(element.Nodes) != 0 {
		err := processModifyNodeElement(changesetId, element.Nodes, uploadResult)
		if err != nil {
			return err
		}
	}

	if len(element.Ways) != 0 {
		return nil
	}

	if len(element.Relations) != 0 {
		return nil
	}

	return nil
}
func ProcessCreateElement(changesetId int64, element Elements, uploadResult *UploadChangesetResult) error {

	if len(element.Nodes) != 0 {
		err := processCreateNodeElement(changesetId, element.Nodes, uploadResult)
		if err != nil {
			return err
		}
	}

	if len(element.Ways) != 0 {
		return nil
	}

	if len(element.Relations) != 0 {
		return nil
	}

	return nil
}

func processModifyNodeElement(changesetId int64, element []Node, uploadResult *UploadChangesetResult) error {
	for _, node := range element {

		x := geo.Lon2x(node.Lon)
		y := geo.Lon2x(node.Lat)
		tile := int64(geo.XY2Tile(x, y))

		//OSM::APIChangesetMismatchError
		if node.Changeset != changesetId {
			return errors.New("invalid id")
		}

		currentNode, err := respository.GetCurrentNodeById(node.ID)
		if err != nil {
			return err
		}

		data := make(map[interface{}]interface{})
		data["changeset_id"] = changesetId
		data["version"] = node.Version + 1
		if int(node.Lat*1e7) != currentNode.Latitude {
			data["latitude"] = int(node.Lat * 1e7)
			data["tile"] = tile
		}

		if int(node.Lon*1e7) != currentNode.Longitude {
			data["longitude"] = int(node.Lon * 1e7)
			data["tile"] = tile
		}

		err = respository.UpdateCurrentNode(strconv.FormatInt(currentNode.ID, 10), data)
		if err != nil {
			return err
		}

		err = repository.DeleteCurrentNodeTags(strconv.FormatInt(currentNode.ID, 10))
		if err != nil {
			return err
		}

		oldNode := gormModel.Nodes{
			Latitude:    int(node.Lat * 1e7),
			Longitude:   int(node.Lon * 1e7),
			ChangesetId: changesetId,
			Visible:     true,
			Timestamp:   time.Now(),
			Tile:        tile,
			Version:     int64(node.Version),
		}

		oldNodeTags := make([]gormModel.NodeTags, 0)
		currentNodeTags := make([]gormModel.CurrentNodeTags, 0)
		for _, tag := range node.Tags {
			currentNodeTags = append(currentNodeTags, gormModel.CurrentNodeTags{
				K: tag.Key,
				V: tag.Value,
			})

		}

		for _, tag := range node.Tags {
			oldNodeTags = append(oldNodeTags, gormModel.NodeTags{
				K:       tag.Key,
				V:       tag.Value,
				Version: int64(node.Version),
			})
		}
		err = repository3.ModifyNode(strconv.FormatInt(currentNode.ID, 10), data, &oldNode, currentNodeTags, oldNodeTags)
		if err != nil {
			return err
		}

		uploadResult.Node = append(uploadResult.Node, UploadChangesetNode{
			Id:      currentNode.ID,
			OldId:   currentNode.ID,
			Version: node.Version,
			Lat:     int(node.Lat * 1e7),
			Lon:     int(node.Lon * 1e7),
			Action:  "modify",
		})

	}

	return nil
}
func processCreateNodeElement(changesetId int64, element []Node, uploadChangesetResult *UploadChangesetResult) error {
	for _, node := range element {

		if node.Changeset != changesetId {
			return errors.New("invalid id")
		}

		x := geo.Lon2x(node.Lon)
		y := geo.Lon2x(node.Lat)
		tile := int64(geo.XY2Tile(x, y))

		currentNode := gormModel.CurrentNodes{
			Latitude:    int(node.Lat * 1e7),
			Longitude:   int(node.Lon * 1e7),
			ChangesetId: changesetId,
			Visible:     true,
			Timestamp:   time.Now(),
			Tile:        tile,
			Version:     int64(node.Version),
		}

		oldNode := gormModel.Nodes{
			Latitude:    int(node.Lat * 1e7),
			Longitude:   int(node.Lon * 1e7),
			ChangesetId: changesetId,
			Visible:     true,
			Timestamp:   time.Now(),
			Tile:        tile,
			Version:     int64(node.Version),
		}

		oldNodeTags := make([]gormModel.NodeTags, 0)
		currentNodeTags := make([]gormModel.CurrentNodeTags, 0)
		for _, tag := range node.Tags {
			currentNodeTags = append(currentNodeTags, gormModel.CurrentNodeTags{
				K: tag.Key,
				V: tag.Value,
			})

		}

		for _, tag := range node.Tags {
			oldNodeTags = append(oldNodeTags, gormModel.NodeTags{
				K:       tag.Key,
				V:       tag.Value,
				Version: int64(node.Version),
			})
		}

		err := repository3.CreateNewNode(&currentNode, &oldNode, currentNodeTags, oldNodeTags)
		if err != nil {
			return err
		}

		uploadChangesetResult.Node = append(uploadChangesetResult.Node, UploadChangesetNode{
			Id:      currentNode.ID,
			OldId:   node.ID,
			Version: node.Version,
			Lat:     int(node.Lat * 1e7),
			Lon:     int(node.Lon * 1e7),
			Action:  "create",
		})

	}

	return nil
}
func processDeleteNodeElement(changesetId int64, element []Node, uploadChangesetResult *UploadChangesetResult) error {
	for _, node := range element {

		if node.Changeset != changesetId {
			return errors.New("invalid id")
		}

		x := geo.Lon2x(node.Lon)
		y := geo.Lon2x(node.Lat)
		tile := int64(geo.XY2Tile(x, y))

		currentNode := gormModel.CurrentNodes{
			Latitude:    int(node.Lat * 1e7),
			Longitude:   int(node.Lon * 1e7),
			ChangesetId: changesetId,
			Visible:     true,
			Timestamp:   time.Now(),
			Tile:        tile,
			Version:     int64(node.Version),
		}

		oldNode := gormModel.Nodes{
			Latitude:    int(node.Lat * 1e7),
			Longitude:   int(node.Lon * 1e7),
			ChangesetId: changesetId,
			Visible:     true,
			Timestamp:   time.Now(),
			Tile:        tile,
			Version:     int64(node.Version),
		}

		oldNodeTags := make([]gormModel.NodeTags, 0)
		currentNodeTags := make([]gormModel.CurrentNodeTags, 0)
		for _, tag := range node.Tags {
			currentNodeTags = append(currentNodeTags, gormModel.CurrentNodeTags{
				K: tag.Key,
				V: tag.Value,
			})

		}

		for _, tag := range node.Tags {
			oldNodeTags = append(oldNodeTags, gormModel.NodeTags{
				K:       tag.Key,
				V:       tag.Value,
				Version: int64(node.Version),
			})
		}

		err := repository3.CreateNewNode(&currentNode, &oldNode, currentNodeTags, oldNodeTags)
		if err != nil {
			return err
		}

		uploadChangesetResult.Node = append(uploadChangesetResult.Node, UploadChangesetNode{
			Id:      currentNode.ID,
			OldId:   node.ID,
			Version: node.Version,
			Lat:     int(node.Lat * 1e7),
			Lon:     int(node.Lon * 1e7),
			Action:  "create",
		})

	}

	return nil
}
