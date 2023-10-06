package main

import (
	"net/http"
	"time"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	// Declare an envelope map containing the data for the response. Notice that the way
	// we've constructed this means the environment and version data will now be nested
	// under a system_info key in the JSON response.
	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}
	time.Sleep(4 * time.Second)

	// writeJSON Creating a writeJSON helper method with helpers.go
	err := app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		// Use the new serverErrorResponse() helper
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) healthCheckHandler2(w http.ResponseWriter, r *http.Request) {
	env := envelope{
		"status": "notAvailable",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}

	err := app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) healthCheckHandler3(w http.ResponseWriter, r *http.Request) {

	env := envelope{
		"status": "availableIfYouCanCatch",
		"system_info": map[string]string{
			"environment": app.config.env,
		},
	}
	err := app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
