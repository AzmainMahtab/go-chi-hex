package handlers // package handlers

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

// Health check handler func
func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	// Data payload for the response
	response := map[string]string{
		"status":      "UP",
		"service":     "DocPad API",
		"environment": "development",
	}

	if err := jsonutil.WriteJSON(w, http.StatusOK, response); err != nil {
		log.Printf("ERROR in helth check %v", err)
	}
}
