package repository

import (
	"database/sql"
	"github.com/lib/pq"
	"gorm.io/gorm"
	dbClient "openstreetmap-go/src/db/client"
	gormModel "openstreetmap-go/src/db/generated/gorm"
)

func GetRelationByVisibleNodes(nodeIds []int64) ([]gormModel.CurrentRelations, error) {
	return getRelationByVisibleNodes(dbClient.DatabaseClients.GetMasterConnection(), nodeIds)
}

func GetRelationByVisibleWays(nodeIds []int64) ([]gormModel.CurrentRelations, error) {
	return getRelationByVisibleWays(dbClient.DatabaseClients.GetMasterConnection(), nodeIds)
}

func GetRelationByVisibleRelations(relationIds []int64) ([]gormModel.CurrentRelations, error) {
	return getRelationByVisibleRelations(dbClient.DatabaseClients.GetMasterConnection(), relationIds)
}

func GetDetailedRelationsByIDs(relationIDs []int64) ([]gormModel.CurrentRelations, error) {
	return getDetailedRelationsByIDs(dbClient.DatabaseClients.GetMasterConnection(), relationIDs)
}

func GetCurrentRelationById(id int64) (gormModel.CurrentRelations, error) {
	return getCurrentRelationById(dbClient.DatabaseClients.GetMasterConnection(), id)
}

func getRelationByVisibleNodes(db *gorm.DB, nodeIds []int64) ([]gormModel.CurrentRelations, error) {
	var relationsMap = make(map[int64]*gormModel.CurrentRelations)
	var results []gormModel.CurrentRelations

	query := `
		SELECT
			-- current_relations
			current_relations.id AS relation_id,
			current_relations.changeset_id,
			current_relations.timestamp,
			current_relations.visible,
			current_relations.version,

			-- changesets
			changesets.id AS changeset_id,
			changesets.user_id,
			changesets.created_at,
			changesets.min_lat,
			changesets.max_lat,
			changesets.min_lon,
			changesets.max_lon,
			changesets.closed_at,
			changesets.num_changes,

			-- current_relation_members
			current_relation_members.relation_id AS member_relation_id,
			current_relation_members.member_type,
			current_relation_members.member_id,
			current_relation_members.member_role,
			current_relation_members.sequence_id,

			-- current_relation_tags
			current_relation_tags.relation_id AS tag_relation_id,
			current_relation_tags.k,
			current_relation_tags.v,

			-- user
			users.display_name,
			users.id

		FROM current_relations
		INNER JOIN current_relation_members ON current_relation_members.relation_id = current_relations.id
		LEFT JOIN changesets ON changesets.id = current_relations.changeset_id
		LEFT JOIN current_relation_tags ON current_relation_tags.relation_id = current_relations.id
		LEFT JOIN users ON changesets.user_id = users.id
		WHERE current_relation_members.member_type = 'Node'
			AND current_relation_members.member_id = ANY(?)
	`

	rows, err := db.Raw(query, pq.Array(nodeIds)).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			relation  gormModel.CurrentRelations
			changeset gormModel.Changesets
			user      gormModel.Users
			member    gormModel.CurrentRelationMembers
			tag       gormModel.CurrentRelationTags

			nullK           sql.NullString
			nullV           sql.NullString
			nullDisplayName sql.NullString
			nullUserId      sql.NullInt64
		)

		err := rows.Scan(
			&relation.ID,
			&relation.ChangesetId,
			&relation.Timestamp,
			&relation.Visible,
			&relation.Version,

			&changeset.ID,
			&changeset.UserId,
			&changeset.CreatedAt,
			&changeset.MinLat,
			&changeset.MaxLat,
			&changeset.MinLon,
			&changeset.MaxLon,
			&changeset.ClosedAt,
			&changeset.NumChanges,

			&member.RelationId,
			&member.MemberType,
			&member.MemberId,
			&member.MemberRole,
			&member.SequenceId,

			&tag.RelationId,
			&nullK,
			&nullV,

			&nullDisplayName,
			&nullUserId,
		)
		if err != nil {
			return nil, err
		}

		if nullDisplayName.Valid {
			user.DisplayName = nullDisplayName.String
			user.ID = nullUserId.Int64
			changeset.User = user
		}

		r, exists := relationsMap[relation.ID]
		if !exists {
			relation.Changeset = changeset
			relation.CurrentRelationMembers = []gormModel.CurrentRelationMembers{}
			relation.CurrentRelationTags = []gormModel.CurrentRelationTags{}
			relationsMap[relation.ID] = &relation
			r = &relation
		}

		r.CurrentRelationMembers = append(r.CurrentRelationMembers, member)

		if nullK.Valid && nullV.Valid {
			tag.K = nullK.String
			tag.V = nullV.String
			r.CurrentRelationTags = append(r.CurrentRelationTags, tag)
		}
	}

	for _, rel := range relationsMap {
		results = append(results, *rel)
	}

	return results, nil
}

func getRelationByVisibleWays(db *gorm.DB, wayIds []int64) ([]gormModel.CurrentRelations, error) {
	var relationsMap = make(map[int64]*gormModel.CurrentRelations)
	var results []gormModel.CurrentRelations

	query := `
		SELECT
			-- current_relations
			current_relations.id AS relation_id,
			current_relations.changeset_id,
			current_relations.timestamp,
			current_relations.visible,
			current_relations.version,

			-- changesets
			changesets.id AS changeset_id,
			changesets.user_id,
			changesets.created_at,
			changesets.min_lat,
			changesets.max_lat,
			changesets.min_lon,
			changesets.max_lon,
			changesets.closed_at,
			changesets.num_changes,

			-- current_relation_members
			current_relation_members.relation_id AS member_relation_id,
			current_relation_members.member_type,
			current_relation_members.member_id,
			current_relation_members.member_role,
			current_relation_members.sequence_id,

			-- current_relation_tags
			current_relation_tags.relation_id AS tag_relation_id,
			current_relation_tags.k,
			current_relation_tags.v,

			-- users
			users.display_name,
			users.id AS user_id

		FROM current_relations
		INNER JOIN current_relation_members ON current_relation_members.relation_id = current_relations.id
		LEFT JOIN changesets ON changesets.id = current_relations.changeset_id
		LEFT JOIN current_relation_tags ON current_relation_tags.relation_id = current_relations.id
		LEFT JOIN users ON changesets.user_id = users.id
		WHERE current_relation_members.member_type = 'Way'
			AND current_relation_members.member_id = ANY(?)
	`

	rows, err := db.Raw(query, pq.Array(wayIds)).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			relation  gormModel.CurrentRelations
			changeset gormModel.Changesets
			user      gormModel.Users
			member    gormModel.CurrentRelationMembers
			tag       gormModel.CurrentRelationTags

			nullK           sql.NullString
			nullV           sql.NullString
			nullDisplayName sql.NullString
			nullUserId      sql.NullInt64
		)

		err := rows.Scan(
			&relation.ID,
			&relation.ChangesetId,
			&relation.Timestamp,
			&relation.Visible,
			&relation.Version,

			&changeset.ID,
			&changeset.UserId,
			&changeset.CreatedAt,
			&changeset.MinLat,
			&changeset.MaxLat,
			&changeset.MinLon,
			&changeset.MaxLon,
			&changeset.ClosedAt,
			&changeset.NumChanges,

			&member.RelationId,
			&member.MemberType,
			&member.MemberId,
			&member.MemberRole,
			&member.SequenceId,

			&tag.RelationId,
			&nullK,
			&nullV,

			&nullDisplayName,
			&nullUserId,
		)
		if err != nil {
			return nil, err
		}

		if nullDisplayName.Valid {
			user.DisplayName = nullDisplayName.String
			user.ID = nullUserId.Int64
			changeset.User = user
		}

		r, exists := relationsMap[relation.ID]
		if !exists {
			relation.Changeset = changeset
			relation.CurrentRelationMembers = []gormModel.CurrentRelationMembers{}
			relation.CurrentRelationTags = []gormModel.CurrentRelationTags{}
			relationsMap[relation.ID] = &relation
			r = &relation
		}

		r.CurrentRelationMembers = append(r.CurrentRelationMembers, member)

		if nullK.Valid && nullV.Valid {
			tag.K = nullK.String
			tag.V = nullV.String
			r.CurrentRelationTags = append(r.CurrentRelationTags, tag)
		}
	}

	for _, rel := range relationsMap {
		results = append(results, *rel)
	}

	return results, nil
}

func getRelationByVisibleRelations(db *gorm.DB, relationIds []int64) ([]gormModel.CurrentRelations, error) {
	var relationsMap = make(map[int64]*gormModel.CurrentRelations)
	var results []gormModel.CurrentRelations

	query := `
		SELECT
			-- current_relations
			current_relations.id AS relation_id,
			current_relations.changeset_id,
			current_relations.timestamp,
			current_relations.visible,
			current_relations.version,

			-- changesets
			changesets.id AS changeset_id,
			changesets.user_id,
			changesets.created_at,
			changesets.min_lat,
			changesets.max_lat,
			changesets.min_lon,
			changesets.max_lon,
			changesets.closed_at,
			changesets.num_changes,

			-- current_relation_members
			current_relation_members.relation_id AS member_relation_id,
			current_relation_members.member_type,
			current_relation_members.member_id,
			current_relation_members.member_role,
			current_relation_members.sequence_id,

			-- current_relation_tags
			current_relation_tags.relation_id AS tag_relation_id,
			current_relation_tags.k,
			current_relation_tags.v,

			-- users
			users.display_name,
			users.id AS user_id

		FROM current_relations
		INNER JOIN current_relation_members ON current_relation_members.relation_id = current_relations.id
		LEFT JOIN changesets ON changesets.id = current_relations.changeset_id
		LEFT JOIN current_relation_tags ON current_relation_tags.relation_id = current_relations.id
		LEFT JOIN users ON changesets.user_id = users.id
		WHERE current_relation_members.member_type = 'Relation'
			AND current_relation_members.member_id = ANY(?)
			AND current_relations.visible = TRUE
	`

	rows, err := db.Raw(query, pq.Array(relationIds)).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			relation  gormModel.CurrentRelations
			changeset gormModel.Changesets
			user      gormModel.Users
			member    gormModel.CurrentRelationMembers
			tag       gormModel.CurrentRelationTags

			nullK           sql.NullString
			nullV           sql.NullString
			nullDisplayName sql.NullString
			nullUserId      sql.NullInt64
		)

		err := rows.Scan(
			&relation.ID,
			&relation.ChangesetId,
			&relation.Timestamp,
			&relation.Visible,
			&relation.Version,

			&changeset.ID,
			&changeset.UserId,
			&changeset.CreatedAt,
			&changeset.MinLat,
			&changeset.MaxLat,
			&changeset.MinLon,
			&changeset.MaxLon,
			&changeset.ClosedAt,
			&changeset.NumChanges,

			&member.RelationId,
			&member.MemberType,
			&member.MemberId,
			&member.MemberRole,
			&member.SequenceId,

			&tag.RelationId,
			&nullK,
			&nullV,

			&nullDisplayName,
			&nullUserId,
		)
		if err != nil {
			return nil, err
		}

		if nullDisplayName.Valid {
			user.DisplayName = nullDisplayName.String
			user.ID = nullUserId.Int64
			changeset.User = user
		}

		r, exists := relationsMap[relation.ID]
		if !exists {
			relation.Changeset = changeset
			relation.CurrentRelationMembers = []gormModel.CurrentRelationMembers{}
			relation.CurrentRelationTags = []gormModel.CurrentRelationTags{}
			relationsMap[relation.ID] = &relation
			r = &relation
		}

		r.CurrentRelationMembers = append(r.CurrentRelationMembers, member)

		if nullK.Valid && nullV.Valid {
			tag.K = nullK.String
			tag.V = nullV.String
			r.CurrentRelationTags = append(r.CurrentRelationTags, tag)
		}
	}

	for _, r := range relationsMap {
		results = append(results, *r)
	}

	return results, nil
}

func getDetailedRelationsByIDs(db *gorm.DB, relationIDs []int64) ([]gormModel.CurrentRelations, error) {
	var relationsMap = make(map[int64]*gormModel.CurrentRelations)
	var results []gormModel.CurrentRelations

	query := `
		SELECT
			-- current_relations
			current_relations.id AS relation_id,
			current_relations.changeset_id,
			current_relations.timestamp,
			current_relations.visible,
			current_relations.version,

			-- changesets
			changesets.id AS changeset_id,
			changesets.user_id,
			changesets.created_at,
			changesets.min_lat,
			changesets.max_lat,
			changesets.min_lon,
			changesets.max_lon,
			changesets.closed_at,
			changesets.num_changes,

			-- current_relation_members
			current_relation_members.relation_id AS member_relation_id,
			current_relation_members.member_type,
			current_relation_members.member_id,
			current_relation_members.member_role,
			current_relation_members.sequence_id,

			-- current_relation_tags
			current_relation_tags.relation_id AS tag_relation_id,
			current_relation_tags.k,
			current_relation_tags.v,

			-- users
			users.display_name,
			users.id AS user_id

		FROM current_relations
		LEFT JOIN current_relation_members ON current_relation_members.relation_id = current_relations.id
		LEFT JOIN current_relation_tags ON current_relation_tags.relation_id = current_relations.id
		LEFT JOIN changesets ON changesets.id = current_relations.changeset_id
		LEFT JOIN users ON users.id = changesets.user_id
		WHERE current_relations.id = ANY(?)
			AND current_relations.visible = TRUE
	`

	rows, err := db.Raw(query, pq.Array(relationIDs)).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			relation  gormModel.CurrentRelations
			changeset gormModel.Changesets
			user      gormModel.Users
			member    gormModel.CurrentRelationMembers
			tag       gormModel.CurrentRelationTags

			nullK           sql.NullString
			nullV           sql.NullString
			nullDisplayName sql.NullString
			nullUserId      sql.NullInt64
		)

		err := rows.Scan(
			&relation.ID,
			&relation.ChangesetId,
			&relation.Timestamp,
			&relation.Visible,
			&relation.Version,

			&changeset.ID,
			&changeset.UserId,
			&changeset.CreatedAt,
			&changeset.MinLat,
			&changeset.MaxLat,
			&changeset.MinLon,
			&changeset.MaxLon,
			&changeset.ClosedAt,
			&changeset.NumChanges,

			&member.RelationId,
			&member.MemberType,
			&member.MemberId,
			&member.MemberRole,
			&member.SequenceId,

			&tag.RelationId,
			&nullK,
			&nullV,

			&nullDisplayName,
			&nullUserId,
		)
		if err != nil {
			return nil, err
		}

		if nullDisplayName.Valid {
			user.DisplayName = nullDisplayName.String
			user.ID = nullUserId.Int64
			changeset.User = user
		}

		r, exists := relationsMap[relation.ID]
		if !exists {
			relation.Changeset = changeset
			relation.CurrentRelationMembers = []gormModel.CurrentRelationMembers{}
			relation.CurrentRelationTags = []gormModel.CurrentRelationTags{}
			relationsMap[relation.ID] = &relation
			r = &relation
		}

		if member.RelationId != 0 {
			r.CurrentRelationMembers = append(r.CurrentRelationMembers, member)
		}

		if nullK.Valid && nullV.Valid {
			tag.K = nullK.String
			tag.V = nullV.String
			r.CurrentRelationTags = append(r.CurrentRelationTags, tag)
		}
	}

	for _, r := range relationsMap {
		results = append(results, *r)
	}

	return results, nil
}

func getCurrentRelationById(client *gorm.DB, id int64) (gormModel.CurrentRelations, error) {

	var currentRelations gormModel.CurrentRelations

	err := client.Model(&currentRelations).Where("id = ?", id).First(&currentRelations).Error
	if err != nil {
		return gormModel.CurrentRelations{}, err
	}
	return currentRelations, nil
}
