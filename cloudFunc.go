package umlgen

import (
	"encoding/json"
	"github.com/cass-dlcm/umlgen/lib"
	"net/http"
	"strings"
)

type writer struct {
	text strings.Builder
}

func (w writer) Write(p []byte) (n int, err error) {
	return w.text.Write(p)
}

func GenDiagram(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "405 - Method Not Allowed (use POST instead)", http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "400 - Bad Request (use \"application/json\" instead)", http.StatusBadRequest)
	}
	var diagram lib.Diagram
	writer := writer{ strings.Builder{} }
	if err := json.NewDecoder(r.Body).Decode(&diagram); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	lib.Generate(writer, diagram)
	w.Header().Add("Content-Type", "image/svg+xml")
	if _, err := w.Write([]byte(writer.text.String())); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
