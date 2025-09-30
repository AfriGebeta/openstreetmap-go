package v1

import (
	"encoding/xml"
	"errors"
	"fmt"
	gormModel "openstreetmap-go/src/db/generated/gorm"
	repository3 "openstreetmap-go/src/modules/changeset/repository"
	"openstreetmap-go/src/modules/current_nodes/respository"
	repository2 "openstreetmap-go/src/modules/current_relation/repository"
	"openstreetmap-go/src/modules/current_way_nodes/repository"
	"openstreetmap-go/src/utils/geo"
	"strconv"
	"time"
)

func findNodeIdFromUploadChangesetResult(id int64, Nodes []UploadChangesetNode) int64 {
	for _, uploadChangesetNode := range Nodes {
		if uploadChangesetNode.OldId == id {
			return uploadChangesetNode.Id
		}
	}
	return -1
}

func findWayIdFromUploadChangesetResult(id int64, Ways []UploadChangesetWay) int64 {
	for _, uploadChangesetway := range Ways {
		if uploadChangesetway.OldId == id {
			return uploadChangesetway.NewId
		}
	}
	return -1
}

func findRelationIdFromUploadChangesetResult(id int64, Relation []UploadChangesetRelation) int64 {
	for _, uploadChangesetRelation := range Relation {
		if uploadChangesetRelation.OldId == id {
			return uploadChangesetRelation.NewId
		}
	}
	return -1
}

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

		if node.NewVersion > 0 {
			nodeDiff.NewVersion = node.NewVersion
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

func ProcessCreateElement(changesetId int64, osm OsmChange, changeset gormModel.Changesets, userId int64, element Elements, uploadResult *UploadChangesetResult) error {
	if len(element.Nodes) != 0 {
		err := processCreateNodeElements(changesetId, element.Nodes, uploadResult)
		if err != nil {
			return err
		}
	}

	if len(element.Ways) != 0 {
		err := processCreateWayElements(changesetId, element.Ways, uploadResult)
		if err != nil {
			return err
		}
		return nil
	}

	if len(element.Relations) != 0 {
		err := processCreateRelationElements(changesetId, element.Relations, uploadResult)
		if err != nil {
			return err
		}
	}

	return nil
}
func ProcessModifyElement(changesetId int64, element Elements, uploadResult *UploadChangesetResult) error {
	if len(element.Nodes) != 0 {
		err := processModifyNodeElements(changesetId, element.Nodes, uploadResult)
		if err != nil {
			return err
		}
	}

	if len(element.Ways) != 0 {
		err := processModifyWayElements(changesetId, element.Ways, uploadResult)
		if err != nil {
			return err
		}
	}

	if len(element.Relations) != 0 {
		err := processModifyRelationElements(changesetId, element.Relations, uploadResult)
		if err != nil {
			return err
		}
	}

	return nil
}
func ProcessDeleteElement(changesetId int64, element Elements, uploadResult *UploadChangesetResult) error {
	if len(element.Nodes) != 0 {
		err := processDeleteNodeElements(changesetId, element.Nodes, uploadResult)
		if err != nil {
			return err
		}
	}

	if len(element.Ways) != 0 {
		err := processDeleteWayElements(changesetId, element.Ways, uploadResult)
		if err != nil {
			return err
		}
	}

	if len(element.Relations) != 0 {
		err := processDeleteRelationElements(changesetId, element.Relations, uploadResult)
		if err != nil {
			return err
		}
	}

	return nil
}

// process node elements
func processCreateNodeElements(changesetId int64, element []Node, uploadChangesetResult *UploadChangesetResult) error {

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
			Version:     int64(node.Version) + 1,
		}

		oldNode := gormModel.Nodes{
			Latitude:    int(node.Lat * 1e7),
			Longitude:   int(node.Lon * 1e7),
			ChangesetId: changesetId,
			Visible:     true,
			Timestamp:   time.Now(),
			Tile:        tile,
			Version:     int64(node.Version) + 1,
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
				Version: int64(node.Version) + 1,
			})
		}

		err := repository3.InsertNode(&currentNode, &oldNode, currentNodeTags, oldNodeTags)
		if err != nil {
			return err
		}

		uploadChangesetResult.Node = append(uploadChangesetResult.Node, UploadChangesetNode{
			Id:         currentNode.ID,
			OldId:      node.ID,
			NewVersion: node.Version + 1,
			Lat:        int(node.Lat * 1e7),
			Lon:        int(node.Lon * 1e7),
			Action:     "create",
		})

	}

	return nil
}

func processModifyNodeElements(changesetId int64, element []Node, uploadResult *UploadChangesetResult) error {
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

		if int64(node.Version) != currentNode.Version {
			return errors.New("invalid version")
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

		oldNode := gormModel.Nodes{
			Latitude:    int(node.Lat * 1e7),
			Longitude:   int(node.Lon * 1e7),
			ChangesetId: changesetId,
			Visible:     true,
			Timestamp:   time.Now(),
			Tile:        tile,
			Version:     int64(node.Version) + 1,
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
				Version: int64(node.Version) + 1,
			})
		}
		err = repository3.ModifyNode(strconv.FormatInt(currentNode.ID, 10), data, &oldNode, currentNodeTags, oldNodeTags)
		if err != nil {
			return err
		}

		uploadResult.Node = append(uploadResult.Node, UploadChangesetNode{
			Id:         currentNode.ID,
			OldId:      currentNode.ID,
			NewVersion: node.Version + 1,
			Lat:        int(node.Lat * 1e7),
			Lon:        int(node.Lon * 1e7),
			Action:     "modify",
		})

	}

	return nil
}

func processDeleteNodeElements(changesetId int64, element []Node, uploadChangesetResult *UploadChangesetResult) error {
	for _, node := range element {

		if node.Changeset != changesetId {
			return errors.New("invalid id")
		}
		currentNode, err := respository.GetCurrentNodeById(node.ID)
		if err != nil {
			return err
		}

		// check if the node exists
		var data map[string]interface{}
		data["changeset_id"] = changesetId
		data["visible"] = false
		data["version"] = currentNode.Version + 1
		data["timestamp"] = time.Now()

		oldNode := gormModel.Nodes{
			Latitude:    int(currentNode.Latitude * 1e7),
			Longitude:   int(currentNode.Longitude * 1e7),
			ChangesetId: changesetId,
			Visible:     true,
			Timestamp:   time.Now(),
			Tile:        currentNode.Tile,
			Version:     int64(currentNode.Version) + 1,
		}

		err = repository3.DeleteNode(node.ID, data, &oldNode)
		if err != nil {
			return err
		}

		uploadChangesetResult.Node = append(uploadChangesetResult.Node, UploadChangesetNode{
			OldId:  node.ID,
			Action: "delete",
		})

	}

	return nil
}

// process way elements
func processCreateWayElements(changesetId int64, element []Way, uploadChangesetResult *UploadChangesetResult) error {
	for _, way := range element {

		if way.Changeset != changesetId {
			return errors.New("invalid id")
		}

		currentWay := gormModel.CurrentWays{
			ChangesetId: changesetId,
			Visible:     true,
			Timestamp:   time.Now(),
			Version:     int64(way.Version) + 1,
		}

		oldway := gormModel.Ways{
			ChangesetId: changesetId,
			Visible:     true,
			Timestamp:   time.Now(),
			Version:     int64(way.Version) + 1,
		}

		currentWayTags := make([]gormModel.CurrentWayTags, 0)
		oldWayTags := make([]gormModel.WayTags, 0)

		for _, tag := range way.Tags {
			currentWayTags = append(currentWayTags, gormModel.CurrentWayTags{
				K: tag.Key,
				V: tag.Value,
			})
			oldWayTags = append(oldWayTags, gormModel.WayTags{
				K:       tag.Key,
				V:       tag.Value,
				Version: int64(way.Version) + 1,
			})
		}

		currentWayNodes := make([]gormModel.CurrentWayNodes, 0)
		oldWayNodes := make([]gormModel.WayNodes, 0)

		for index, nd := range way.NDs {
			nodeId := findNodeIdFromUploadChangesetResult(nd.Ref, uploadChangesetResult.Node)
			if nodeId < 0 {
				return errors.New("invalid id")
			}

			if nd.Ref < 0 {
				currentWayNodes = append(currentWayNodes, gormModel.CurrentWayNodes{
					SequenceId: int64(index),
					NodeId:     nodeId,
				})
				oldWayNodes = append(oldWayNodes, gormModel.WayNodes{
					SequenceId: int64(index),
					NodeId:     nodeId,
				})
			} else {
				currentWayNodes = append(currentWayNodes, gormModel.CurrentWayNodes{
					SequenceId: int64(index),
					NodeId:     nodeId,
				})
				oldWayNodes = append(oldWayNodes, gormModel.WayNodes{
					SequenceId: int64(index),
					NodeId:     nodeId,
				})
			}
		}

		err := repository3.InsertWays(&currentWay, currentWayNodes, currentWayTags, oldway, oldWayNodes, oldWayTags)
		if err != nil {
			return err
		}

		uploadChangesetResult.Way = append(uploadChangesetResult.Way, UploadChangesetWay{
			NewId:      currentWay.ID,
			OldId:      way.ID,
			NewVersion: int64(way.Version + 1),
		})

	}

	return nil
}

func processModifyWayElements(changesetId int64, element []Way, uploadChangesetResult *UploadChangesetResult) error {

	for _, way := range element {
		if way.Changeset != changesetId {
			return errors.New("invalid id")
		}
		wayModify := map[string]interface{}{}
		wayModify["changeset_id"] = changesetId
		wayModify["timestamp"] = time.Now()
		wayModify["version"] = int64(way.Version) + 1

		oldway := gormModel.Ways{
			ChangesetId: changesetId,
			Visible:     true,
			Timestamp:   time.Now(),
			Version:     int64(way.Version) + 1,
		}

		currentWayTags := make([]gormModel.CurrentWayTags, 0)
		oldWayTags := make([]gormModel.WayTags, 0)

		for _, tag := range way.Tags {
			currentWayTags = append(currentWayTags, gormModel.CurrentWayTags{
				K: tag.Key,
				V: tag.Value,
			})
			oldWayTags = append(oldWayTags, gormModel.WayTags{
				K:       tag.Key,
				V:       tag.Value,
				Version: int64(way.Version) + 1,
			})
		}

		// get the way nodes
		wayNodes, err := repository.GetCurrentWayNodesByWayId(way.ID)
		if err != nil {
			return err
		}

		currentWayNodes := make([]gormModel.CurrentWayNodes, 0)
		oldWayNodes := make([]gormModel.WayNodes, 0)

		for index, nd := range way.NDs {
			for _, wayNode := range wayNodes {
				if wayNode.NodeId == nd.Ref {
					if nd.Ref < 0 {
						nodeId := findNodeIdFromUploadChangesetResult(nd.Ref, uploadChangesetResult.Node)
						if nodeId < 0 {
							return errors.New("invalid id")
						}
						currentWayNodes = append(currentWayNodes, gormModel.CurrentWayNodes{
							SequenceId: int64(index),
							WayId:      way.ID,
							NodeId:     nodeId,
						})
						oldWayNodes = append(oldWayNodes, gormModel.WayNodes{
							SequenceId: int64(index),
							WayId:      way.ID,
							NodeId:     nodeId,
							Version:    int64(way.Version) + 1,
						})
					} else {
						currentWayNodes = append(currentWayNodes, gormModel.CurrentWayNodes{
							SequenceId: int64(index),
							WayId:      way.ID,
							NodeId:     wayNode.NodeId,
						})
						oldWayNodes = append(oldWayNodes, gormModel.WayNodes{
							SequenceId: int64(index),
							WayId:      way.ID,
							NodeId:     wayNode.NodeId,
							Version:    int64(way.Version) + 1,
						})
					}

				}
			}

		}

		err = repository3.ModifyWay(way.ID, wayModify, currentWayNodes, currentWayTags, oldway, oldWayNodes, oldWayTags)
		if err != nil {
			return err
		}

		uploadChangesetResult.Way = append(uploadChangesetResult.Way, UploadChangesetWay{
			NewId:      way.ID,
			OldId:      way.ID,
			NewVersion: int64(way.Version + 1),
		})

	}
	return nil
}

func processDeleteWayElements(changesetId int64, element []Way, uploadChangesetResult *UploadChangesetResult) error {
	for _, way := range element {

		if way.Changeset != changesetId {
			return errors.New("invalid id")
		}
		currentWay, err := respository.GetCurrentNodeById(way.ID)
		if err != nil {
			return err
		}

		// check if the node exists
		var data map[string]interface{}
		data["changeset_id"] = changesetId
		data["visible"] = false
		data["version"] = currentWay.Version + 1
		data["timestamp"] = time.Now()

		oldway := gormModel.Ways{
			ChangesetId: changesetId,
			Visible:     false,
			Timestamp:   time.Now(),
			Version:     int64(way.Version) + 1,
		}

		err = repository3.DeleteWay(way.ID, data, &oldway)
		if err != nil {
			return err
		}

		uploadChangesetResult.Node = append(uploadChangesetResult.Node, UploadChangesetNode{
			OldId:  way.ID,
			Action: "delete",
		})

	}

	return nil
}

// process relation elements
func processCreateRelationElements(changesetId int64, element []Relation, uploadChangesetResult *UploadChangesetResult) error {
	for _, relation := range element {
		if relation.Changeset != changesetId {
			return errors.New("invalid id")
		}

		currentRelation := gormModel.CurrentRelations{
			ChangesetId: changesetId,
			Visible:     true,
			Timestamp:   time.Now(),
			Version:     int64(relation.Version + 1),
		}

		oldRelation := gormModel.Relations{
			ChangesetId: changesetId,
			Visible:     true,
			Timestamp:   time.Now(),
			Version:     int64(relation.Version) + 1,
		}

		currentRelationTags := make([]gormModel.CurrentRelationTags, 0)
		oldRelationTags := make([]gormModel.RelationTags, 0)

		for _, relationTag := range relation.Tags {
			currentRelationTags = append(currentRelationTags, gormModel.CurrentRelationTags{
				K: relationTag.Key,
				V: relationTag.Value,
			})
			oldRelationTags = append(oldRelationTags, gormModel.RelationTags{
				K:       relationTag.Key,
				V:       relationTag.Value,
				Version: int64(relation.Version) + 1,
			})
		}

		currentRelationMembers := make([]gormModel.CurrentRelationMembers, 0)
		relationMembers := make([]gormModel.RelationMembers, 0)

		for index, relationMember := range relation.Members {
			if relationMember.Ref < 0 {
				var memberId int64
				if relationMember.Type == "Way" {
					memberId = findWayIdFromUploadChangesetResult(relationMember.Ref, uploadChangesetResult.Way)
				} else if relationMember.Type == "Relation" {
					memberId = findRelationIdFromUploadChangesetResult(relationMember.Ref, uploadChangesetResult.Relation)
				} else {
					memberId = findNodeIdFromUploadChangesetResult(relationMember.Ref, uploadChangesetResult.Node)
				}

				if memberId < 0 {
					return fmt.Errorf("invalid member id")
				}
				currentRelationMembers = append(currentRelationMembers, gormModel.CurrentRelationMembers{
					MemberType: relationMember.Type,
					SequenceId: index,
					MemberRole: relationMember.Role,
					MemberId:   memberId,
				})
				relationMembers = append(relationMembers, gormModel.RelationMembers{
					MemberType: relationMember.Type,
					SequenceId: index,
					MemberRole: relationMember.Role,
					MemberId:   relationMember.Ref,
				})
			} else {
				currentRelationMembers = append(currentRelationMembers, gormModel.CurrentRelationMembers{
					MemberType: relationMember.Type,
					SequenceId: index,
					MemberRole: relationMember.Role,
					MemberId:   relationMember.Ref,
				})
				relationMembers = append(relationMembers, gormModel.RelationMembers{
					MemberType: relationMember.Type,
					SequenceId: index,
					MemberRole: relationMember.Role,
					MemberId:   relationMember.Ref,
				})

			}
		}

		err := repository3.CreateRelation(&currentRelation, oldRelation, currentRelationTags, currentRelationMembers, relationMembers, oldRelationTags)
		if err != nil {
			return err
		}

		uploadChangesetResult.Relation = append(uploadChangesetResult.Relation, UploadChangesetRelation{
			OldId:  relation.ID,
			NewId:  currentRelation.ID,
			Action: "delete",
		})

	}

	return nil
}

func processModifyRelationElements(changesetId int64, element []Relation, uploadChangesetResult *UploadChangesetResult) error {

	for _, relation := range element {
		if relation.Changeset != changesetId {
			return errors.New("invalid id")
		}
		relationModify := map[string]interface{}{}
		relationModify["changeset_id"] = changesetId
		relationModify["timestamp"] = time.Now()
		relationModify["version"] = int64(relation.Version) + 1

		oldRelation := gormModel.Relations{
			ChangesetId: changesetId,
			Visible:     true,
			Timestamp:   time.Now(),
			Version:     int64(relation.Version) + 1,
			RelationId:  relation.ID,
		}

		currentRelationTags := make([]gormModel.CurrentRelationTags, 0)
		oldRelationTags := make([]gormModel.RelationTags, 0)

		for _, relationTag := range relation.Tags {
			currentRelationTags = append(currentRelationTags, gormModel.CurrentRelationTags{
				K:          relationTag.Key,
				V:          relationTag.Value,
				RelationId: relation.ID,
			})
			oldRelationTags = append(oldRelationTags, gormModel.RelationTags{
				K:          relationTag.Key,
				V:          relationTag.Value,
				Version:    int64(relation.Version) + 1,
				RelationId: relation.ID,
			})
		}
		currentRelationMembers := make([]gormModel.CurrentRelationMembers, 0)
		relationMembers := make([]gormModel.RelationMembers, 0)
		for index, relationMember := range relation.Members {
			if relationMember.Ref < 0 {
				var memberId int64
				if relationMember.Type == "Way" {
					memberId = findWayIdFromUploadChangesetResult(relationMember.Ref, uploadChangesetResult.Way)
				} else if relationMember.Type == "Relation" {
					memberId = findRelationIdFromUploadChangesetResult(relationMember.Ref, uploadChangesetResult.Relation)
				} else {
					memberId = findNodeIdFromUploadChangesetResult(relationMember.Ref, uploadChangesetResult.Node)
				}

				if memberId < 0 {
					return fmt.Errorf("invalid member id")
				}
				currentRelationMembers = append(currentRelationMembers, gormModel.CurrentRelationMembers{
					MemberType: relationMember.Type,
					SequenceId: index,
					MemberRole: relationMember.Role,
					MemberId:   memberId,
					RelationId: relation.ID,
				})
				relationMembers = append(relationMembers, gormModel.RelationMembers{
					MemberType: relationMember.Type,
					SequenceId: index,
					MemberRole: relationMember.Role,
					MemberId:   relationMember.Ref,
					RelationId: relation.ID,
				})
			} else {
				currentRelationMembers = append(currentRelationMembers, gormModel.CurrentRelationMembers{
					MemberType: relationMember.Type,
					SequenceId: index,
					MemberRole: relationMember.Role,
					MemberId:   relationMember.Ref,
					RelationId: relation.ID,
				})
				relationMembers = append(relationMembers, gormModel.RelationMembers{
					MemberType: relationMember.Type,
					SequenceId: index,
					MemberRole: relationMember.Role,
					MemberId:   relationMember.Ref,
					RelationId: relation.ID,
				})

			}
		}

		err := repository3.ModifyRelation(relation.ID, relationModify, oldRelation, currentRelationTags, oldRelationTags, currentRelationMembers, relationMembers)
		if err != nil {
			return err
		}

		uploadChangesetResult.Relation = append(uploadChangesetResult.Relation, UploadChangesetRelation{
			NewId:      relation.ID,
			OldId:      relation.ID,
			NewVersion: int64(relation.Version + 1),
			Action:     "modify",
		})

	}

	return nil
}

func processDeleteRelationElements(changesetId int64, element []Relation, uploadChangesetResult *UploadChangesetResult) error {
	for _, relation := range element {

		if relation.Changeset != changesetId {
			return errors.New("invalid id")
		}

		currentRelation, err := repository2.GetCurrentRelationById(relation.ID)
		if err != nil {
			return err
		}

		// check if the node exists
		var data map[string]interface{}
		data["changeset_id"] = changesetId
		data["visible"] = false
		data["version"] = currentRelation.Version + 1
		data["timestamp"] = time.Now()

		oldRelation := gormModel.Relations{
			ChangesetId: changesetId,
			Visible:     false,
			Timestamp:   time.Now(),
			Version:     int64(currentRelation.Version) + 1,
			RelationId:  currentRelation.ID,
		}

		err = repository3.DeleteRelation(currentRelation.ID, data, &oldRelation)
		if err != nil {
			return err
		}

		uploadChangesetResult.Relation = append(uploadChangesetResult.Relation, UploadChangesetRelation{
			OldId:  currentRelation.ID,
			Action: "delete",
		})

	}

	return nil
}
