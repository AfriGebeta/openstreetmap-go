package respository

import (
	"database/sql"
	"gorm.io/gorm"
	"openstreetmap-go/src/config"
	dbClient "openstreetmap-go/src/db/client"
	gormModel "openstreetmap-go/src/db/generated/gorm"
	"openstreetmap-go/src/repository/model"
)

func CreateCurrentNode(currentNode *gormModel.CurrentNodes) error {
	return createCurrentNode(dbClient.DatabaseClients.GetMasterConnection(), currentNode)
}

func createCurrentNode(client *gorm.DB, currentNode *gormModel.CurrentNodes) error {
	err := client.Model(&gormModel.CurrentNodes{}).Create(&currentNode).Error
	if err != nil {
		return err
	}
	return nil
}

func FetchCurrentNodeInBBox(bbox model.BBox, tileWhereClause string) ([]gormModel.CurrentNodes, error) {
	return fetchCurrentNodeInBBox(dbClient.DatabaseClients.GetRandomConnection(), bbox, tileWhereClause)
}
func fetchCurrentNodeInBBox(db *gorm.DB, bbox model.BBox, tileWhereClause string) ([]gormModel.CurrentNodes, error) {
	var nodesMap = make(map[int64]*gormModel.CurrentNodes)
	var nodes []gormModel.CurrentNodes

	sqlQuery := `
		SELECT 
		 "current_nodes".id as current_node_id,
                          "current_nodes".latitude,
                          "current_nodes".longitude,
                          "current_nodes".changeset_id,
                          "current_nodes".visible,
                          "current_nodes".timestamp,
                          "current_nodes".tile, 
                          node_tags.node_id,
                          node_tags.k,
                          node_tags.v, 
                          changesets.id,
                          changesets.user_id,
                          changesets.created_at,
                          changesets.min_lat,
                          changesets.max_lat,
                          changesets.min_lon,
                          changesets.max_lon,
                          changesets.closed_at,
                          changesets.num_changes,
                          users.display_name,
						  users.id
 					from 
                        "current_nodes"
                        LEFT JOIN node_tags ON current_nodes.id = node_tags.node_id
                        LEFT JOIN changesets ON current_nodes.changeset_id = changesets.id
                        LEFT JOIN users ON changesets.user_id = users.id	
		WHERE
		` + tileWhereClause + `	
			AND current_nodes.latitude BETWEEN $1 AND $2 
			AND current_nodes.longitude BETWEEN $3 AND $4 
			AND current_nodes.visible = TRUE
		LIMIT $5
	`

	rows, err := db.Raw(sqlQuery,

		bbox.MinLat*1e7,
		bbox.MaxLat*1e7,
		bbox.MinLon*1e7,
		bbox.MaxLon*1e7,
		config.ServerConfigObject.OSMSettings.MaxNumberOfNodes,
	).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			node      gormModel.CurrentNodes
			changeset gormModel.Changesets
			user      gormModel.Users

			nullNodeId sql.NullInt64
			nullK      sql.NullString
			nullV      sql.NullString
		)

		err := rows.Scan(
			&node.ID,
			&node.Latitude,
			&node.Longitude,
			&node.ChangesetId,
			&node.Visible,
			&node.Timestamp,
			&node.Tile,
			&nullNodeId,
			&nullK,
			&nullV,
			&changeset.ID,
			&changeset.UserId,
			&changeset.CreatedAt,
			&changeset.MinLat,
			&changeset.MaxLat,
			&changeset.MinLon,
			&changeset.MaxLon,
			&changeset.ClosedAt,
			&changeset.NumChanges,
			&user.DisplayName,
			&user.ID,
		)
		if err != nil {
			return nil, err
		}

		_, exists := nodesMap[node.ID]
		if !exists {

			node.Changeset = changeset
			node.Changeset.User = user
			node.CurrentNodeTags = []gormModel.CurrentNodeTags{}
			nodesMap[node.ID] = &node
		}
		if nullK.Valid && nullV.Valid && nullNodeId.Valid {

			nodesMap[node.ID].CurrentNodeTags = append(nodesMap[node.ID].CurrentNodeTags, gormModel.CurrentNodeTags{
				NodeId: nullNodeId.Int64,
				K:      nullK.String,
				V:      nullV.String,
			})

		}

	}

	for _, n := range nodesMap {
		nodes = append(nodes, *n)
	}

	return nodes, nil
}

//func fetchCurrentNodeInBBox(client *gorm.DB, bbox model.BBox, tileWhereClause string) ([]gormModel.CurrentNodes, error) {
//	var nodes []gormModel.CurrentNodes
//	query := client.Model(&gormModel.CurrentNodes{})
//
//	query = query.
//		Where("current_nodes.latitude BETWEEN ? AND ?", bbox.MinLat*1e7, bbox.MaxLat*1e7).
//		Where("current_nodes.longitude BETWEEN ? AND ?", bbox.MinLon*1e7, bbox.MaxLon*1e7).
//		Where("current_nodes.visible = ?", true).
//		Preload("CurrentNodeTags").
//		Preload("Changeset").
//		Preload("Changeset.User")
//
//	if tileWhereClause != "" {
//		query = query.Where(tileWhereClause)
//	}
//
//	result := query.
//		Limit(config.ServerConfigObject.OSMSettings.MaxNumberOfNodes).
//		Find(&nodes)
//
//	if result.Error != nil {
//		return nil, result.Error
//	}
//	return nodes, nil
//}
