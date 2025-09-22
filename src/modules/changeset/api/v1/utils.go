package v1

import (
	"encoding/xml"
	gormModel "openstreetmap-go/src/db/generated/gorm"
	repository3 "openstreetmap-go/src/modules/changeset/repository"
	"openstreetmap-go/src/utils/geo"
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
func ProcessCreateElement(changesetId int64, element Elements, uploadResult *UploadChangesetResult) error {
	if len(element.Nodes) == 0 {
		err := ProcessCreateNodeElement(changesetId, element.Nodes, uploadResult)
		if err != nil {
			return err
		}
	}

	if len(element.Ways) == 0 {
		return nil
	}

	if len(element.Relations) == 0 {
		return nil
	}

	return nil
}

func ProcessCreateNodeElement(changeSetId int64, element []Node, uploadChangesetResult *UploadChangesetResult) error {
	for _, node := range element {
		x := geo.Lon2x(node.Lon)
		y := geo.Lon2x(node.Lat)
		tile := int64(geo.XY2Tile(x, y))

		currentNode := gormModel.CurrentNodes{
			Latitude:    int(node.Lat * 1e7),
			Longitude:   int(node.Lon * 1e7),
			ChangesetId: changeSetId,
			Visible:     true,
			Timestamp:   time.Now(),
			Tile:        tile,
			Version:     int64(node.Version),
		}

		oldNode := gormModel.Nodes{
			Latitude:    int(node.Lat * 1e7),
			Longitude:   int(node.Lon * 1e7),
			ChangesetId: changeSetId,
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
