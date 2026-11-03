package main

import (
	"net/http"
)

// HealthCheck godoc
//
//	@Summary		Health check
//	@Description	Returns the API status, environment and version.
//	@Tags			ops
//	@Produce		json
//	@Success		200	{object}	map[string]string
//	@Failure		500	{object}	error	"Failed to write the response"
//
//	@Security		ApiKeyAuth
//
//	@Router			/health [get]
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"env":     app.config.env,
		"version": version,
	}

	if err := app.jsonResponse(w, http.StatusOK, data); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
