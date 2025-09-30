package v1

import "encoding/xml"

type Osm struct {
	XMLName   xml.Name  `xml:"osm"`
	Changeset Changeset `xml:"changeset"`
}

type Changeset struct {
	ID   int   `xml:"id,attr,omitempty"`
	Open bool  `xml:"open,attr,omitempty"`
	Tags []Tag `xml:"tag"`
}

type OsmChange struct {
	XMLName xml.Name `xml:"osmChange"`
	Version string   `xml:"version,attr"`
	// Optional metadata
	Generator string   `xml:"generator,attr,omitempty"`
	Create    Elements `xml:"create"`
	Modify    Elements `xml:"modify"`
	Delete    Elements `xml:"delete"`
}

type Elements struct {
	Nodes     []Node     `xml:"node"`
	Ways      []Way      `xml:"way"`
	Relations []Relation `xml:"relation"`
}

type Node struct {
	ID        int64   `xml:"id,attr"`
	Changeset int64   `xml:"changeset,attr"`
	Version   int     `xml:"version,attr,omitempty"`
	Lat       float64 `xml:"lat,attr,omitempty"`
	Lon       float64 `xml:"lon,attr,omitempty"`
	Tags      []Tag   `xml:"tag"`
}

type Way struct {
	ID        int64 `xml:"id,attr"`
	Changeset int64 `xml:"changeset,attr"`
	Version   int   `xml:"version,attr,omitempty"`
	NDs       []ND  `xml:"nd"`
	Tags      []Tag `xml:"tag"`
}

type ND struct {
	Ref int64 `xml:"ref,attr"`
}

type Relation struct {
	ID        int64    `xml:"id,attr"`
	Changeset int64    `xml:"changeset,attr"`
	Version   int      `xml:"version,attr,omitempty"`
	Members   []Member `xml:"member"`
	Tags      []Tag    `xml:"tag"`
}

type Member struct {
	Type string `xml:"type,attr"` // node, way, or relation
	Ref  int64  `xml:"ref,attr"`
	Role string `xml:"role,attr"`
}

type Tag struct {
	Key   string `xml:"k,attr"`
	Value string `xml:"v,attr"`
}

type UploadChangesetResult struct {
	Node     []UploadChangesetNode     `xml:"node"`
	Way      []UploadChangesetWay      `xml:"way"`
	Relation []UploadChangesetRelation `xml:"relation"`
}

type UploadChangesetRelation struct {
	OldId      int64  `xml:"old_id,attr"`
	NewId      int64  `xml:"new_id,attr"`
	NewVersion int64  `xml:"new_version,attr"`
	Action     string `xml:"action,attr"`
}
type UploadChangesetWay struct {
	OldId      int64 `xml:"old_id,attr"`
	NewId      int64 `xml:"new_id,attr"`
	NewVersion int64 `xml:"new_version,attr"`
}

type UploadChangesetNode struct {
	Id         int64  `xml:"new_id,attr"`
	OldId      int64  `xml:"old_id,attr"`
	NewVersion int    `xml:"new_version,attr,omitempty"`
	Lat        int    `xml:"lat,attr,omitempty"`
	Lon        int    `xml:"long,attr,omitempty"`
	Action     string `xml:"action,attr,omitempty"`
}

type DiffResult struct {
	XMLName     xml.Name   `xml:"diffResult"`
	Version     string     `xml:"version,attr"`
	Generator   string     `xml:"generator,attr"`
	Copyright   string     `xml:"copyright,attr"`
	Attribution string     `xml:"attribution,attr"`
	License     string     `xml:"license,attr"`
	Nodes       []NodeDiff `xml:"node"`
}

type NodeDiff struct {
	OldID      *int `xml:"old_id,attr,omitempty"` // use pointer to omit if nil
	NewID      int  `xml:"new_id,attr"`
	NewVersion int  `xml:"new_version,attr"`
}
