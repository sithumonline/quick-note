package handler

import (
	"net/http"

	"github.com/sithumonline/quick-note/config"
)

// GetHealth godoc
// @Summary Returns health of the service
// @Router /healthz [get]
func GetHealth(w http.ResponseWriter, r *http.Request) {
	RespondWithJSON(w, http.StatusOK, "Go Note up and running version: "+config.GetEnv("quick-note.VERSION"))
}
