package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

type PageData struct {
	BackgroundColor string
	SecretSource    string
	AppVersion      string
	Hostname        string
	Timestamp       string
	VaultPath       string
	Namespace       string
}

var tmpl *template.Template

func init() {
	var err error
	tmpl, err = template.ParseFiles("templates/index.html")
	if err != nil {
		// Fallback: try embedded template if file not found
		tmpl = template.Must(template.New("index").Parse(fallbackTemplate))
	}
}

func main() {
	port := getEnv("PORT", "8080")

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/healthz", handleHealth)
	http.HandleFunc("/readyz", handleReady)

	log.Printf("Starting vault-color-demo on :%s", port)
	log.Printf("Background color: %s", getEnv("BG_COLOR", "#1a1a2e"))
	log.Printf("Secret source: %s", getEnv("SECRET_SOURCE", "default"))

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname()

	data := PageData{
		BackgroundColor: getEnv("BG_COLOR", "#1a1a2e"),
		SecretSource:    getEnv("SECRET_SOURCE", "default (no secret mounted)"),
		AppVersion:      getEnv("APP_VERSION", "dev"),
		Hostname:        hostname,
		Timestamp:       time.Now().Format(time.RFC3339),
		VaultPath:       getEnv("VAULT_SECRET_PATH", "N/A"),
		Namespace:       getEnv("POD_NAMESPACE", "unknown"),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, fmt.Sprintf("Template error: %v", err), http.StatusInternalServerError)
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")
}

func handleReady(w http.ResponseWriter, r *http.Request) {
	// Check if we have a secret mounted
	if getEnv("BG_COLOR", "") == "" {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprint(w, "waiting for secret")
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ready")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

const fallbackTemplate = `<!DOCTYPE html>
<html>
<head><title>Vault Color Demo</title></head>
<body style="background:{{.BackgroundColor}};color:#fff;font-family:sans-serif;display:flex;align-items:center;justify-content:center;height:100vh;margin:0">
<div style="text-align:center">
<h1>HashiCorp Vault + Kubernetes Demo</h1>
<p>Background color: {{.BackgroundColor}}</p>
<p>Source: {{.SecretSource}}</p>
</div>
</body>
</html>`
