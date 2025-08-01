package v1

import "encoding/xml"

type Osm struct {
	XMLName     xml.Name `xml:"osm"`
	Version     string   `xml:"version,attr"`
	Generator   string   `xml:"generator,attr"`
	Copyright   string   `xml:"copyright,attr"`
	Attribution string   `xml:"attribution,attr"`
	License     string   `xml:"license,attr"`
	API         API      `xml:"api"`
	Policy      []Policy `xml:"policy"`
}

type API struct {
	Version         Version         `xml:"version"`
	Area            Area            `xml:"area"`
	NoteArea        NoteArea        `xml:"note_area"`
	Tracepoints     Tracepoints     `xml:"tracepoints"`
	WayNodes        WayNodes        `xml:"waynodes"`
	RelationMembers RelationMembers `xml:"relationmembers"`
	Changesets      Changesets      `xml:"changesets"`
	Notes           Notes           `xml:"notes"`
	Timeout         Timeout         `xml:"timeout"`
	Status          Status          `xml:"status"`
}

type Version struct {
	Minimum string `xml:"minimum,attr"`
	Maximum string `xml:"maximum,attr"`
}

type Area struct {
	Maximum float32 `xml:"maximum,attr"`
}

type NoteArea struct {
	Maximum float32 `xml:"maximum,attr"`
}

type Tracepoints struct {
	PerPage int32 `xml:"per_page,attr"`
}

type WayNodes struct {
	Maximum uint32 `xml:"maximum,attr"`
}

type RelationMembers struct {
	Maximum uint32 `xml:"maximum,attr"`
}

type Changesets struct {
	MaximumElements   string `xml:"maximum_elements,attr"`
	DefaultQueryLimit int32  `xml:"default_query_limit,attr"`
	MaximumQueryLimit int32  `xml:"maximum_query_limit,attr"`
}

type Notes struct {
	DefaultQueryLimit int32 `xml:"default_query_limit,attr"`
	MaximumQueryLimit int32 `xml:"maximum_query_limit,attr"`
}

type Timeout struct {
	Seconds int32 `xml:"seconds,attr"`
}

type Status struct {
	Database string `xml:"database,attr"`
	API      string `xml:"api,attr"`
	GPX      string `xml:"gpx,attr"`
}

type Policy struct {
	Imagery string `xml:"imagery"`
}
