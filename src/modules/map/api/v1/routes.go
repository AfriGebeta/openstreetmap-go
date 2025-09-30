package v1

import (
	"errors"
	"fmt"
	"html"
	"net/http"
	"openstreetmap-go/src/api/middlewares"
	gormModel "openstreetmap-go/src/db/generated/gorm"
	"openstreetmap-go/src/modules/current_nodes/respository"
	repository2 "openstreetmap-go/src/modules/current_relation/repository"
	repository3 "openstreetmap-go/src/modules/current_way_nodes/repository"
	"openstreetmap-go/src/modules/current_ways/repository"
	"openstreetmap-go/src/repository/model"
	"openstreetmap-go/src/utils/geo"
	"strings"
	"time"
)

func Setup() *http.ServeMux {
	route := http.NewServeMux()
	route.HandleFunc("/map", middlewares.HandleError(GetMap))
	return route
}

func GetMap(w http.ResponseWriter, r *http.Request) error {

	bboxParameter := r.URL.Query().Get("bbox")
	if bboxParameter == "" {
		return errors.New("invalid input")
	}
	bbox, err := FromBBoxParams(bboxParameter)
	if err != nil {
		return err
	}

	err = Checkboundaries(bbox)
	if err != nil {
		return err
	}

	err = CheckAreaSize(bbox)
	if err != nil {
		return err
	}

	tileWhereClause := geo.BuildTileAreaSQLClause(bbox, "current_nodes.")

	// Note: Step 1 - filter out current nodes by tile
	currentNodes, err := respository.GetCurrentNodeInBBox(bbox, tileWhereClause)
	if err != nil {
		return err
	}

	if len(currentNodes) != 0 {
		nodeMaps := map[int64]gormModel.CurrentNodes{}
		currentNodeIds := []int64{}
		for _, node := range currentNodes {
			currentNodeIds = append(currentNodeIds, node.ID)
			nodeMaps[node.ID] = node
		}
		// Note: step 2 Get the current way nodes from tha current nodes to get the current_ways
		currentWaysNodes, err := repository3.GetCurrentWayNodes(currentNodeIds)
		if err != nil {
			return err
		}

		currentWayIdsFromWaysNodes := []int64{}
		for _, way := range currentWaysNodes {
			currentWayIdsFromWaysNodes = append(currentWayIdsFromWaysNodes, way.WayId)
		}

		// Note: step 3 Get the current ways from the collected wayIds from currentWayNodes
		ways, err := repository.FetchCurrentWaysByIds(currentWayIdsFromWaysNodes)
		if err != nil {
			return err
		}
		wayIds := []int64{}
		wayMaps := map[int64]gormModel.CurrentWays{}
		for _, way := range ways {
			wayIds = append(wayIds, way.ID)
			wayMaps[way.ID] = way
		}
		// Note: step 4 Get the current Nodes using way ids
		currentWayNodesByWayIds, err := repository3.GetCurrentWayNodesByWaysIds(wayIds)
		if err != nil {
			return err
		}

		for _, currentWayNode := range currentWayNodesByWayIds {
			if existingNode, ok := nodeMaps[currentWayNode.Node.ID]; !ok {
				nodeMaps[currentWayNode.Node.ID] = currentWayNode.Node
			} else {
				existingNode.CurrentNodeTags = append(existingNode.CurrentNodeTags, currentWayNode.Node.CurrentNodeTags...)
				nodeMaps[existingNode.ID] = existingNode
			}
		}

		// Note; step 5 get the current relation
		visibleNodeIds := []int64{}
		for _, node := range currentNodes {
			if node.Visible {
				visibleNodeIds = append(visibleNodeIds, node.ID)
			}
		}

		visibleWayIds := []int64{}
		for _, way := range ways {
			if way.Visible {
				visibleWayIds = append(visibleWayIds, way.ID)
			}
		}
		//Note; step 6 get the visible relation by nodes
		relationByNodes, err := repository2.FetchRelationByVisibleNodes(visibleNodeIds)
		if err != nil {
			return err
		}

		//Note; step 7 get the visible way by ways
		relationByWay, err := repository2.FetchRelationByVisibleWays(visibleWayIds)
		if err != nil {
			return err
		}

		relationIds := []int64{}
		relationMaps := map[int64]gormModel.CurrentRelations{}

		for _, relation := range relationByWay {
			relationIds = append(relationIds, relation.ID)
			relationMaps[relation.ID] = relation
		}

		for _, relation := range relationByNodes {
			relationIds = append(relationIds, relation.ID)
			if existingRelation, ok := relationMaps[relation.ID]; !ok {
				relationMaps[relation.ID] = relation
				relationIds = append(relationIds, relation.ID)
			} else {
				existingRelation.CurrentRelationTags = append(existingRelation.CurrentRelationTags, relation.CurrentRelationTags...)
				existingRelation.CurrentRelationMembers = append(existingRelation.CurrentRelationMembers, relation.CurrentRelationMembers...)
				existingRelation.Changeset = relation.Changeset
				relationMaps[relation.ID] = existingRelation
			}
		}

		//Note; step 8 get
		visibleRelations, err := repository2.FetchRelationByVisibleRelations(relationIds)
		if err != nil {
			return err
		}
		relationIds = []int64{}
		for _, relation := range visibleRelations {
			relationIds = append(relationIds, relation.ID)
		}

		relationdata, err := repository2.FetchDetailedRelationsByIDs(relationIds)
		if err != nil {
			return err
		}
		for _, relation := range relationdata {
			if existingRelation, ok := relationMaps[relation.ID]; !ok {
				relationMaps[relation.ID] = relation
			} else {
				existingRelation.CurrentRelationTags = append(existingRelation.CurrentRelationTags, relation.CurrentRelationTags...)
				existingRelation.CurrentRelationMembers = append(existingRelation.CurrentRelationMembers, relation.CurrentRelationMembers...)
				existingRelation.Changeset = relation.Changeset
			}
		}
		result, err := ToOsmXml(bbox, nodeMaps, wayMaps, relationMaps)
		if err != nil {
			return err
		}
		w.Header().Set("Content-Disposition", `attachment; filename="map.osm"`)
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write([]byte(result))
	} else {
		http.Error(w, "Not found", http.StatusNotFound)

	}

	return nil
}

func ToOsmXml(bounds model.BBox, nodes map[int64]gormModel.CurrentNodes, ways map[int64]gormModel.CurrentWays, relations map[int64]gormModel.CurrentRelations) (string, error) {
	header := `<?xml version="1.0" encoding="UTF-8"?>`
	osm := `<osm version="0.6" generator="MyNodeOsmServer"  copyright="OpenStreetMap and contributors" attribution="http://www.openstreetmap.org/copyright" license="http://opendatacommons.org/licenses/odbl/1-0/">`
	boundsXml := fmt.Sprintf(
		`  <bounds minlat="%.7f" minlon="%.7f" maxlat="%.7f" maxlon="%.7f"/>`,
		bounds.MinLat, bounds.MinLon, bounds.MaxLat, bounds.MaxLon,
	)
	var nodeXmls []string
	var wayXmls []string
	fmt.Println(len(nodes))
	for _, node := range nodes {
		tagsXml := buildNodeTagsXml(node.CurrentNodeTags)
		nodeXml := ""
		if len(tagsXml) > 0 {
			nodeXml = fmt.Sprintf(
				`  <node id="%d" visible="%t" version="%d" changeset="%d" timestamp="%s" user="%s" uid="%d" lat="%.7f" lon="%.7f">
%s
  </node>`,
				node.ID, node.Visible, node.Version, node.Changeset.ID,
				node.Timestamp.Format(time.RFC3339),
				node.Changeset.User.DisplayName, node.Changeset.UserId,
				float64(node.Latitude)/1e7, float64(node.Longitude)/1e7,
				tagsXml,
			)
		} else {
			nodeXml = fmt.Sprintf(
				`  <node id="%d" visible="%t" version="%d" changeset="%d" timestamp="%s" user="%s" uid="%d" lat="%.7f" lon="%.7f"/>`,
				node.ID, node.Visible, node.Version, node.Changeset.ID,
				node.Timestamp.Format(time.RFC3339),
				node.Changeset.User.DisplayName, node.Changeset.UserId,
				float64(node.Latitude)/1e7, float64(node.Longitude)/1e7,
			)
		}
		nodeXmls = append(nodeXmls, nodeXml)
	}

	for _, way := range ways {
		var nds []string
		filteredNds := func(arr []gormModel.CurrentWayNodes) []gormModel.CurrentWayNodes {
			if len(arr) == 0 {
				return []gormModel.CurrentWayNodes{}
			}

			result := []gormModel.CurrentWayNodes{arr[0]}
			for i := 1; i < len(arr); i++ {
				if arr[i].NodeId != arr[i-1].NodeId {
					result = append(result, arr[i])
				}
			}

			return result
		}
		for _, id := range filteredNds(way.CurrentWayNodes) {
			nds = append(nds, fmt.Sprintf(`    <nd ref="%d" />`, id.NodeId))
		}

		tagsXml := buildWayTagsXml(way.CurrentWayTags)
		wayXml := fmt.Sprintf(
			`  <way id="%d" visible="%t" version="%d" changeset="%d" timestamp="%s" user="%s" uid="%d">
%s
%s
  </way>`,
			way.ID, way.Visible, way.Version, way.Changeset.ID,
			way.Timestamp.Format(time.RFC3339),
			way.Changeset.User.DisplayName, way.Changeset.UserId,
			strings.Join(nds, "\n"), tagsXml,
		)
		wayXmls = append(wayXmls, wayXml)
	}
	var relationXmls []string
	for _, rel := range relations {
		var members []string

		for _, m := range rel.CurrentRelationMembers {
			members = append(members, fmt.Sprintf(`    <member type="%s" ref="%d" role="%s" />`, strings.ToLower(m.MemberType), m.MemberId, m.MemberRole))
		}

		tagsXml := buildRelationTagsXml(rel.CurrentRelationTags)

		relationXml := fmt.Sprintf(
			`  <relation id="%d" visible="true" version="%d" changeset="%d" timestamp="%s" user="Osmosis Anonymous" uid="%d">
%s
%s
  </relation>`,
			rel.ID, rel.Version, rel.Changeset.ID,
			rel.Timestamp.Format(time.RFC3339),
			rel.Changeset.UserId,
			strings.Join(members, "\n"),
			tagsXml,
		)

		relationXmls = append(relationXmls, relationXml)
	}
	footer := `</osm>`
	xmlStr := strings.Join([]string{
		header,
		osm,
		boundsXml,
		strings.Join(nodeXmls, "\n"),
		strings.Join(wayXmls, "\n"),
		strings.Join(relationXmls, "\n"),
		footer,
	}, "\n")

	return xmlStr, nil

}

func buildNodeTagsXml(tags []gormModel.CurrentNodeTags) string {
	if len(tags) == 0 {
		return ""
	}
	var tagLines []string
	for _, tags := range tags {
		tagLines = append(tagLines, fmt.Sprintf(`    <tag k="%s" v="%s" />`, html.EscapeString(tags.K), html.EscapeString(tags.V)))
	}
	return strings.Join(tagLines, "\n")
}

func buildWayTagsXml(wayTags []gormModel.CurrentWayTags) string {
	if len(wayTags) == 0 {
		return ""
	}
	var wayTagLines []string
	for _, wayTag := range wayTags {
		wayTagLines = append(wayTagLines, fmt.Sprintf(`    <tag k="%s" v="%s" />`, html.EscapeString(wayTag.K), html.EscapeString(wayTag.V)))
	}
	return strings.Join(wayTagLines, "\n")
}

func buildRelationTagsXml(currentRelationTags []gormModel.CurrentRelationTags) string {
	if len(currentRelationTags) == 0 {
		return ""
	}
	var tagLines []string
	for _, tag := range currentRelationTags {
		tagLines = append(tagLines, fmt.Sprintf(`    <tag k="%s" v="%s" />`, html.EscapeString(tag.K), html.EscapeString(tag.V)))
	}
	return strings.Join(tagLines, "\n")
}
