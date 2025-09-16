package repository

import (
	"database/sql"
	"github.com/lib/pq"
	"gorm.io/gorm"
	dbClient "openstreetmap-go/src/db/client"
	gormModel "openstreetmap-go/src/db/generated/gorm"
)

func FetchCurrentWaysByIds(ids []int64) ([]gormModel.CurrentWays, error) {
	return fetchCurrentWaysIds(dbClient.DatabaseClients.GetRandomConnection(), ids)
}
func fetchCurrentWaysIds(db *gorm.DB, ids []int64) ([]gormModel.CurrentWays, error) {
	var waysMap = make(map[int64]*gormModel.CurrentWays)
	var ways []gormModel.CurrentWays

	sqlQuery := `
        SELECT 
            w.id AS way_id, 
            w.changeset_id AS way_changeset_id, 
            w.timestamp AS way_timestamp,
            w.visible AS way_visible,
            w.version AS way_version,

            wn.node_id AS current_way_node_id, 
            wn.way_id AS current_way_way_id,
            wn.sequence_id AS current_way_sequence_id,

            wt.way_id AS current_way_tags_way_id,
            wt.k  AS current_way_tags_k,
            wt.v  AS current_way_tags_v,

            c.id AS changeset_id,
            c.created_at AS changeset_created_at,
            c.min_lat AS changeset_min_lat,
            c.max_lat AS changeset_max_lat,
            c.min_lon AS changeset_min_lon,
            c.max_lon AS changeset_max_lon,
            c.closed_at AS changeset_closed_at,
            c.num_changes AS changeset_num_changes,
            c.user_id AS changeset_user_id,

            u.display_name AS user_display_name,
            u.id AS user_id
        FROM current_ways w
        LEFT JOIN current_way_nodes wn ON wn.way_id = w.id
        LEFT JOIN current_way_tags wt ON wt.way_id = w.id
        LEFT JOIN changesets c ON c.id = w.changeset_id
        LEFT JOIN users u ON c.user_id = u.id
        WHERE w.id = ANY(?)
    `

	rows, err := db.Raw(sqlQuery, pq.Array(ids)).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			way       gormModel.CurrentWays
			changeset gormModel.Changesets
			user      gormModel.Users

			nullNodeID    sql.NullInt64
			nullWayNodeID sql.NullInt64
			nullSeqID     sql.NullInt64
			nullTagWayID  sql.NullInt64
			nullK         sql.NullString
			nullV         sql.NullString
		)

		err := rows.Scan(
			&way.ID,
			&way.ChangesetId,
			&way.Timestamp,
			&way.Visible,
			&way.Version,

			&nullNodeID,
			&nullWayNodeID,
			&nullSeqID,

			&nullTagWayID,
			&nullK,
			&nullV,

			&changeset.ID,
			&changeset.CreatedAt,
			&changeset.MinLat,
			&changeset.MaxLat,
			&changeset.MinLon,
			&changeset.MaxLon,
			&changeset.ClosedAt,
			&changeset.NumChanges,
			&changeset.UserId,

			&user.DisplayName,
			&user.ID,
		)
		if err != nil {
			return nil, err
		}

		_, exists := waysMap[way.ID]
		if !exists {
			way.Changeset = changeset
			way.Changeset.User = user
			way.CurrentWayTags = []gormModel.CurrentWayTags{}
			way.CurrentWayNodes = []gormModel.CurrentWayNodes{}
			waysMap[way.ID] = &way
		}

		if nullNodeID.Valid && nullWayNodeID.Valid && nullSeqID.Valid {
			waysMap[way.ID].CurrentWayNodes = append(waysMap[way.ID].CurrentWayNodes, gormModel.CurrentWayNodes{
				NodeId:     nullNodeID.Int64,
				WayId:      nullWayNodeID.Int64,
				SequenceId: int64(int(nullSeqID.Int64)),
			})
		}

		if nullTagWayID.Valid && nullK.Valid && nullV.Valid {
			waysMap[way.ID].CurrentWayTags = append(waysMap[way.ID].CurrentWayTags, gormModel.CurrentWayTags{
				WayId: nullTagWayID.Int64,
				K:     nullK.String,
				V:     nullV.String,
			})
		}
	}

	for _, way := range waysMap {
		ways = append(ways, *way)
	}

	return ways, nil
}

//func fetchCurrentWaysIds(client *gorm.DB, ids []int64) ([]gormModel.CurrentWays, error) {
//	var currentModels []gormModel.CurrentWays
//	query := client.Model(&gormModel.CurrentWays{})
//
//	query = query.
//		Joins("LEFT JOIN current_way_nodes ON current_way_nodes.way_id = current_ways.id").
//		Joins("LEFT JOIN current_way_tags ON current_way_tags.way_id = current_ways.id").
//		Joins("LEFT JOIN changesets ON changesets.id = current_ways.changeset_id").
//		Joins("LEFT JOIN users ON changesets.user_id = users.id").
//		Where("current_ways.id IN (?)", ids)
//
//	err := query.Find(&currentModels).Error
//	if err != nil {
//		return nil, err
//	}
//	return currentModels, err
//}
