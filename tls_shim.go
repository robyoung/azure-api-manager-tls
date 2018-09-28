package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	listenPort := getEnv("LISTEN_PORT", "8080")

	handler := httpHandler{
		version: loadVersion(),
	}

	s := &http.Server{
		Addr:    fmt.Sprintf(":%s", listenPort),
		Handler: &handler,
	}
	s.ListenAndServe()
}

// Environment handling
func getEnv(key, def string) string {
	if res := os.Getenv(key); res == "" {
		return def
	} else {
		return res
	}
}

// Version handling
type Version struct {
	TLSShimVersion string `json:"tlsShimVersion"`
	AppVersion     string `json:"appVersion"`
}

func loadVersion() string {
	v := Version{
		TLSShimVersion: "no version set in tls-shim-version.txt",
	}
	if tlsShimVersion, err := ioutil.ReadFile("tls-shim-version.txt"); err != nil && string(tlsShimVersion) != "" {
		v.TLSShimVersion = string(tlsShimVersion)
	}
	if appVersion, err := ioutil.ReadFile("app-version.txt"); err != nil {
		v.AppVersion = string(appVersion)
	}
	result, err := json.Marshal(v)
	if err != nil {
		panic("Could marshal version payload")
	} else {
		return string(result)
	}
}

func writeVersion(w http.ResponseWriter, v string) {
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(200)
	fmt.Fprint(w, v)
}

// Proxy handling
const clientCertHeader string = "X-Arr-Clientcert"

func handleProxy(w http.ResponseWriter, r *http.Request) {
	clientCert := r.Header.Get(clientCertHeader)
	if clientCert == "" {
		respondUnauthorized(w, r, "Request did not contain a certificate")
	}
}

func doProxy(w http.ResponseWriter, r *http.Request) {

}

func respondUnauthorized(w http.ResponseWriter, r *http.Request, message string) {
	log.Printf("Responding with unauthorized - %s", message)

	w.WriteHeader(403)
}

// HTTP Server
type httpHandler struct {
	version string
}

func (h *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/version" {
		writeVersion(w, h.version)
	} else {
		fmt.Fprintf(w, "Other response %s", r.URL.Path)
	}
}
