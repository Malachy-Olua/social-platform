package main

import (
	"log"
	"net/http"

	"github.com/Malachy-Olua/social-platform/helpers"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "application/json")
	// w.Write([]byte(`{"status": "ok"}`))
	data := map[string]string{
		"status":  "ok",
		"env":     app.config.env,
		"version": version,
	}

	if err := helpers.WriteJSON(w, http.StatusOK, data); err != nil {
		log.Println(err.Error())
	}
}
