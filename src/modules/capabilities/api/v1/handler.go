package v1

import (
	"encoding/xml"
	"net/http"
	"openstreetmap-go/src/config"
)

func getCapabilites(w http.ResponseWriter, r *http.Request) error {
	policy := []Policy{}
	for _, v := range config.ServerConfiguration.ImageryBlacklist {
		policy = append(policy, Policy{
			Imagery: v,
		})
	}

	osm := Osm{
		XMLName:     xml.Name{},
		Version:     config.ServerConfiguration.ApiVersion,
		Generator:   "Gebeta Maps",
		Copyright:   "Gebeta Maps",
		Attribution: "https://gebeta.app/copyright",
		License:     "http://opendatacommons.org/licenses/odbl/1-0/",
		API: API{
			Version: Version{
				Minimum: config.ServerConfiguration.ApiVersion,
				Maximum: config.ServerConfiguration.ApiVersion,
			},
			Area:            Area{Maximum: config.ServerConfiguration.MaxRequestArea},
			NoteArea:        NoteArea{Maximum: config.ServerConfiguration.MaxNoteRequestArea},
			Tracepoints:     Tracepoints{PerPage: config.ServerConfiguration.TracepointsPerPage},
			WayNodes:        WayNodes{Maximum: config.ServerConfiguration.MaxNumberOfWayNodes},
			RelationMembers: RelationMembers{Maximum: config.ServerConfiguration.MaxNumberOfRelationMembers},
			Changesets: Changesets{
				MaximumElements:   "",
				DefaultQueryLimit: config.ServerConfiguration.DefaultChangesetQueryLimit,
				MaximumQueryLimit: config.ServerConfiguration.MaxChangesetQueryLimit,
			},
			Notes: Notes{
				DefaultQueryLimit: config.ServerConfiguration.DefaultNoteQueryLimit,
				MaximumQueryLimit: config.ServerConfiguration.MaxNoteQueryLimit,
			},
			Timeout: Timeout{Seconds: config.ServerConfiguration.ApiTimeout},
			Status: Status{
				Database: "online",
				API:      "online",
				GPX:      "online",
			},
		},
		Policy: policy,
	}

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	xml.NewEncoder(w).Encode(osm)
	return nil
}
