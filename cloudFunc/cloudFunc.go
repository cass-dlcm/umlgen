package cloudFunc

import (
	"encoding/json"
	"github.com/cass-dlcm/umlgen/lib"
	"net/http"
)

func GenDiagram(w http.ResponseWriter, r *http.Request) {
	var diagram lib.Diagram
	err := json.NewDecoder(r.Body).Decode(&diagram)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	lib.Generate(w, diagram)
}
