package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"

	ev "github.com/mchmarny/gcputil/env"
	mt "github.com/mchmarny/gcputil/metric"
)

const (
	metricSourceID      = "ab-test-demo"
	revisionHeaderParam = "Revision-Version"
)

var (
	templates       *template.Template
	queryLimit      = ev.MustGetIntEnvVar("QUERY_LIMIT", 50)
	revisionVersion = ev.MustGetEnvVar("VERSION", "a")
	mtClient        *mt.Client
)

func initHandlers() {
	tmpls, err := template.ParseGlob("template/*.html")
	if err != nil {
		logger.Fatalf("Error while parsing templates: %v", err)
	}
	templates = tmpls

	c, err := mt.NewClientWithSource(context.Background(), metricSourceID)
	if err != nil {
		logger.Fatalf("Error while creating metrics client: %v", err)
	}
	mtClient = c

}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	meterAction(r, "visits")
	w.Header().Set(revisionHeaderParam, revisionVersion)
	if err := templates.ExecuteTemplate(w, "index", getData()); err != nil {
		logger.Printf("Error in view template: %s", err)
	}
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	meterAction(r, "clicks")
	w.Header().Set(revisionHeaderParam, revisionVersion)
	if err := templates.ExecuteTemplate(w, "form", getData()); err != nil {
		logger.Printf("Error in view template: %s", err)
	}
}

func meterAction(r *http.Request, measurement string) {
	if err := mtClient.PublishForSource(r.Context(), measurement, 1); err != nil {
		logger.Printf("Error publishing metrics: %s", err)
	}
}

func getData() map[string]interface{} {
	data := make(map[string]interface{})
	data["version"] = revisionVersion
	data["release"] = fmt.Sprintf("v0.0.1-%s", revisionVersion)
	return data
}
