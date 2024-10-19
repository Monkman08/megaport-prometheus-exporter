package handlers

import "net/http"

// LandingPageHandler provides a simple landing page
func LandingPageHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`<html>
        <head><title>Megaport Exporter</title></head>
        <body>
        <h1>Megaport Exporter</h1>
        <p><a href="/metrics">Metrics</a></p>
        </body>
        </html>`))
}
