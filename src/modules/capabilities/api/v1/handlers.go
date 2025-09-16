package v1

import (
	"encoding/xml"
	"net/http"
)

func capability(w http.ResponseWriter, r *http.Request) error {

	capabilities := CapabilityResponse{
		Version:     "0.6",
		Generator:   "OpenStreetMap server",
		Copyright:   "OpenStreetMap and contributors",
		Attribution: "http://www.openstreetmap.org/copyright",
		License:     "http://opendatacommons.org/licenses/odbl/1-0/",
		API: API{
			Version:         Version{Minimum: "0.6", Maximum: "0.6"},
			Area:            Area{Maximum: "0.25"},
			NoteArea:        NoteArea{Maximum: "25"},
			Tracepoints:     Tracepoints{PerPage: "5000"},
			Waynodes:        Waynodes{Maximum: "2000"},
			Relationmembers: Relationmembers{Maximum: "32000"},
			Changesets: Changesets{
				MaximumElements:   "10000",
				DefaultQueryLimit: "100",
				MaximumQueryLimit: "100",
			},
			Notes: Notes{
				DefaultQueryLimit: "100",
				MaximumQueryLimit: "10000",
			},
			Timeout: Timeout{Seconds: "300"},
			Status:  Status{Database: "online", API: "online", GPX: "online"},
		},
		Policy: Policy{Imagery: ""},
	}
	output, err := xml.MarshalIndent(capabilities, "", "  ")
	if err != nil {
		http.Error(w, "Failed to generate XML", http.StatusInternalServerError)
		return nil
	}

	w.Header().Set("Content-Type", "application/xml")
	w.Write([]byte(xml.Header))
	w.Write(output)

	return nil
}
