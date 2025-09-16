package v1

import "encoding/xml"

type CapabilityResponse struct {
	XMLName     xml.Name `xml:"osm"`
	Generator   string   `xml:"generator,attr"`
	Version     string   `xml:"version,attr"`
	Copyright   string   `xml:"copyright,attr"`
	Attribution string   `xml:"attribution,attr"`
	License     string   `xml:"license,attr"`
	API         API      `xml:"api"`
	Policy      Policy   `xml:"policy"`
}

type API struct {
	Version         Version         `xml:"version"`
	Area            Area            `xml:"area"`
	NoteArea        NoteArea        `xml:"note_area"`
	Tracepoints     Tracepoints     `xml:"tracepoints"`
	Waynodes        Waynodes        `xml:"waynodes"`
	Relationmembers Relationmembers `xml:"relationmembers"`
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
	Maximum string `xml:"maximum,attr"`
}

type NoteArea struct {
	Maximum string `xml:"maximum,attr"`
}

type Tracepoints struct {
	PerPage string `xml:"per_page,attr"`
}

type Waynodes struct {
	Maximum string `xml:"maximum,attr"`
}

type Relationmembers struct {
	Maximum string `xml:"maximum,attr"`
}

type Changesets struct {
	MaximumElements   string `xml:"maximum_elements,attr"`
	DefaultQueryLimit string `xml:"default_query_limit,attr"`
	MaximumQueryLimit string `xml:"maximum_query_limit,attr"`
}

type Notes struct {
	DefaultQueryLimit string `xml:"default_query_limit,attr"`
	MaximumQueryLimit string `xml:"maximum_query_limit,attr"`
}

type Timeout struct {
	Seconds string `xml:"seconds,attr"`
}

type Status struct {
	Database string `xml:"database,attr"`
	API      string `xml:"api,attr"`
	GPX      string `xml:"gpx,attr"`
}

type Policy struct {
	Imagery string `xml:"imagery"`
}
