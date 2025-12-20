// Package handlers stores all the handler
// handlers are HERE!
package handlers

import (
	"log"
	"net/http"

	"github.com/AzmainMahtab/docpad/pkg/jsonutil"
)

type HealthHandler struct {
	// Keep DB connection interface herre
}

func NewHealthHandleer() *HealthHandler {
	return &HealthHandler{}
}

// HealthCheck handles the GET /api/v1/health request.
// @Summary Check the status and uptime of the API server.
// @Description Provides a simple UP/DOWN status and service identification.
// @Tags System
// @Accept json
// @Produce json
// @Success 200 {object} object{status=string,service=string,environment=string} "Successful response indicating API status is UP."
// @Router /health [get]
func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	// Data payload for the response
	response := map[string]string{
		"status":      "UP",
		"service":     "DocPad API",
		"environment": "development",
	}

	if err := jsonutil.WriteJSON(w, http.StatusOK, response, nil); err != nil {
		log.Printf("ERROR in helth check %v", err)
	}
}
