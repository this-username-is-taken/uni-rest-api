package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"regexp"

	"github.com/go-chi/chi"
)

// Volume given time range

// Block number
// What swaps
// List all assets

func validAddress(addr string) bool {
	match, _ := regexp.MatchString("^0x[a-fA-F0-9]{40}$", addr)
	return match
}

func assetPoolsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	assetId := chi.URLParam(r, "assetId")

	if !validAddress(assetId) {
		http.Error(w, "Invalid asset id", http.StatusBadRequest)
		return
	}

	responseJson, err := QueryAssetPools(assetId)

	if err == nil {
		w.Write(responseJson)
	} else {
		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
	}
}

func blockHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET params were:", r.URL.Query())

	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func main() {
	router := chi.NewRouter()
	router.Get("/asset/{assetId}/pools", assetPoolsHandler)

	// http.HandleFunc("/asset/", assetHandler)
	// http.HandleFunc("/block", blockHandler)

	fmt.Println("Listening on port " + ServerPort + "...")
	log.Fatal(http.ListenAndServe(":"+ServerPort, router))
}
