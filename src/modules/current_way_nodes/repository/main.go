package repository

import (
	"database/sql"
	"github.com/lib/pq"
	"gorm.io/gorm"
	dbClient "openstreetmap-go/src/db/client"
	gormModel "openstreetmap-go/src/db/generated/gorm"
)

func GetCurrentWayNodes(ids []int64) ([]gormModel.CurrentWayNodes, error) {
	return getCurrentWayNodes(dbClient.DatabaseClients.GetRandomConnection(), ids)
}
func GetCurrentWayNodesByWaysIds(ids []int64) ([]gormModel.CurrentWayNodes, error) {
	return getCurrentWatNodesByWayIds(dbClient.DatabaseClients.GetMasterConnection(), ids)
}
func GetCurrentWayNodesByWayId(id int64) ([]gormModel.CurrentWayNodes, error) {
	return getCurrentWayNodeByWayId(dbClient.DatabaseClients.GetMasterConnection(), id)
}

func getCurrentWayNodes(client *gorm.DB, nodeIds []int64) ([]gormModel.CurrentWayNodes, error) {
	var currentModels []gormModel.CurrentWayNodes
	query := client.Model(&gormModel.CurrentWayNodes{})

	err := query.Where("current_way_nodes.node_id IN (?)", nodeIds).Find(&currentModels).Error
	if err != nil {
		return nil, err
	}
	return currentModels, err
}

func getCurrentWatNodesByWayIds(client *gorm.DB, ids []int64) ([]gormModel.CurrentWayNodes, error) {
	var wayNodesMap = make(map[int64]*gormModel.CurrentWayNodes) // key by node ID (or combine wayID+nodeID)
	var result []gormModel.CurrentWayNodes

	sqlQuery := `
        SELECT
            current_nodes.id AS current_node_id,
            current_nodes.latitude,
            current_nodes.longitude,
            current_nodes.changeset_id,
            current_nodes.visible,
            current_nodes.timestamp,
            current_nodes.tile,
            current_nodes.version AS current_node_version,

            current_way_nodes.way_id AS current_way_node_way_id,
            current_way_nodes.node_id AS current_way_node_node_id,

            node_tags.node_id,
            node_tags.k,
            node_tags.v,
            node_tags.version,

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
        FROM current_way_nodes
        JOIN current_nodes ON current_nodes.id = current_way_nodes.node_id
        LEFT JOIN node_tags ON current_nodes.id = node_tags.node_id
        LEFT JOIN changesets ON current_nodes.changeset_id = changesets.id
        LEFT JOIN users ON changesets.user_id = users.id
        WHERE current_way_nodes.way_id = ANY(?)
    `

	rows, err := client.Raw(sqlQuery, pq.Array(ids)).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			node  gormModel.CurrentNodes
			tag   gormModel.CurrentNodeTags
			cs    gormModel.Changesets
			user  gormModel.Users
			wayID int64

			nullTagNodeID sql.NullInt64
			nullK         sql.NullString
			nullV         sql.NullString
			nullTagVer    sql.NullInt64
			nullUserName  sql.NullString
			nullUserId    sql.NullInt64
		)

		err := rows.Scan(
			&node.ID,
			&node.Latitude,
			&node.Longitude,
			&node.ChangesetId,
			&node.Visible,
			&node.Timestamp,
			&node.Tile,
			&node.Version,

			&wayID,
			&node.ID,

			&nullTagNodeID,
			&nullK,
			&nullV,
			&nullTagVer,

			&cs.ID,
			&cs.UserId,
			&cs.CreatedAt,
			&cs.MinLat,
			&cs.MaxLat,
			&cs.MinLon,
			&cs.MaxLon,
			&cs.ClosedAt,
			&cs.NumChanges,

			&nullUserName,
			&nullUserId,
		)
		if err != nil {
			return nil, err
		}

		if nullUserName.Valid {
			user.DisplayName = nullUserName.String
			user.ID = nullUserId.Int64
			cs.User = user
		}

		mapKey := wayID*1e10 + node.ID
		_, exists := wayNodesMap[mapKey]
		if !exists {
			wayNode := gormModel.CurrentWayNodes{
				WayId:  wayID,
				NodeId: node.ID,
				Node: gormModel.CurrentNodes{
					ID:              node.ID,
					Latitude:        node.Latitude,
					Longitude:       node.Longitude,
					ChangesetId:     node.ChangesetId,
					Visible:         node.Visible,
					Timestamp:       node.Timestamp,
					Tile:            node.Tile,
					Version:         node.Version,
					Changeset:       cs,
					CurrentNodeTags: []gormModel.CurrentNodeTags{},
				},
			}
			wayNodesMap[mapKey] = &wayNode
		}

		if nullTagNodeID.Valid && nullK.Valid && nullV.Valid {
			tag = gormModel.CurrentNodeTags{
				NodeId: nullTagNodeID.Int64,
				K:      nullK.String,
				V:      nullV.String,
			}

			wayNodesMap[mapKey].Node.CurrentNodeTags = append(wayNodesMap[mapKey].Node.CurrentNodeTags, tag)
		}
	}

	for _, wn := range wayNodesMap {
		result = append(result, *wn)
	}

	return result, nil
}

func getCurrentWayNodeByWayId(client *gorm.DB, id int64) ([]gormModel.CurrentWayNodes, error) {
	var wayNodesMap = make(map[int64]*gormModel.CurrentWayNodes)
	var result []gormModel.CurrentWayNodes

	sqlQuery := `
        SELECT
            current_nodes.id AS current_node_id,
            current_nodes.latitude,
            current_nodes.longitude,
            current_nodes.changeset_id,
            current_nodes.visible,
            current_nodes.timestamp,
            current_nodes.tile,
            current_nodes.version AS current_node_version,

            current_way_nodes.way_id AS current_way_node_way_id,
            current_way_nodes.node_id AS current_way_node_node_id,

            node_tags.node_id,
            node_tags.k,
            node_tags.v,
            node_tags.version,

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
        FROM current_way_nodes
        JOIN current_nodes ON current_nodes.id = current_way_nodes.node_id
        LEFT JOIN node_tags ON current_nodes.id = node_tags.node_id
        LEFT JOIN changesets ON current_nodes.changeset_id = changesets.id
        LEFT JOIN users ON changesets.user_id = users.id
        WHERE current_way_nodes.way_id = ? 
    `

	rows, err := client.Raw(sqlQuery, id).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			node  gormModel.CurrentNodes
			tag   gormModel.CurrentNodeTags
			cs    gormModel.Changesets
			user  gormModel.Users
			wayID int64

			nullTagNodeID sql.NullInt64
			nullK         sql.NullString
			nullV         sql.NullString
			nullTagVer    sql.NullInt64
			nullUserName  sql.NullString
			nullUserId    sql.NullInt64
		)

		err := rows.Scan(
			&node.ID,
			&node.Latitude,
			&node.Longitude,
			&node.ChangesetId,
			&node.Visible,
			&node.Timestamp,
			&node.Tile,
			&node.Version,

			&wayID,
			&node.ID,

			&nullTagNodeID,
			&nullK,
			&nullV,
			&nullTagVer,

			&cs.ID,
			&cs.UserId,
			&cs.CreatedAt,
			&cs.MinLat,
			&cs.MaxLat,
			&cs.MinLon,
			&cs.MaxLon,
			&cs.ClosedAt,
			&cs.NumChanges,

			&nullUserName,
			&nullUserId,
		)
		if err != nil {
			return nil, err
		}

		if nullUserName.Valid {
			user.DisplayName = nullUserName.String
			user.ID = nullUserId.Int64
			cs.User = user
		}

		mapKey := wayID*1e10 + node.ID
		_, exists := wayNodesMap[mapKey]
		if !exists {
			wayNode := gormModel.CurrentWayNodes{
				WayId:  wayID,
				NodeId: node.ID,
				Node: gormModel.CurrentNodes{
					ID:              node.ID,
					Latitude:        node.Latitude,
					Longitude:       node.Longitude,
					ChangesetId:     node.ChangesetId,
					Visible:         node.Visible,
					Timestamp:       node.Timestamp,
					Tile:            node.Tile,
					Version:         node.Version,
					Changeset:       cs,
					CurrentNodeTags: []gormModel.CurrentNodeTags{},
				},
			}
			wayNodesMap[mapKey] = &wayNode
		}

		if nullTagNodeID.Valid && nullK.Valid && nullV.Valid {
			tag = gormModel.CurrentNodeTags{
				NodeId: nullTagNodeID.Int64,
				K:      nullK.String,
				V:      nullV.String,
			}

			wayNodesMap[mapKey].Node.CurrentNodeTags = append(wayNodesMap[mapKey].Node.CurrentNodeTags, tag)
		}
	}

	for _, wn := range wayNodesMap {
		result = append(result, *wn)
	}

	return result, nil
}
