package main

import (
	"fmt"
	"html/template"
	"net/http"

	ev "github.com/mchmarny/gcputil/env"
)

var (
	templates  *template.Template
	queryLimit = ev.MustGetIntEnvVar("QUERY_LIMIT", 50)
)

func initHandlers() {
	tmpls, err := template.ParseGlob("template/*.html")
	if err != nil {
		logger.Fatalf("Error while parsing templates: %v", err)
	}
	templates = tmpls
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	version := ev.MustGetEnvVar("VERSION", "A")
	data["version"] = version
	data["release"] = ev.MustGetEnvVar("RELEASE",
		fmt.Sprintf("v0.0.1-%s", version))

	if err := templates.ExecuteTemplate(w, "index", data); err != nil {
		logger.Printf("Error in view template: %s", err)
	}
}
