package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/go-chi/chi"
)

func validAddress(addr string) bool {
	match, _ := regexp.MatchString("^0x[a-fA-F0-9]{40}$", addr)
	return match
}

func assetPoolsRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	assetId := chi.URLParam(r, "assetId")

	if !validAddress(assetId) {
		http.Error(w, "Invalid asset id", http.StatusBadRequest)
		return
	}

	responseJson, err := queryAssetPools(assetId)

	if err != nil {
		log.Println("Error querying asset volume: " + err.Error())
		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
		return
	} else {
		w.Write(responseJson)
	}
}

func assetVolumeRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	assetId := chi.URLParam(r, "assetId")

	if !validAddress(assetId) {
		http.Error(w, "Invalid asset id", http.StatusBadRequest)
		return
	}

	startParam := r.URL.Query().Get("start")
	endParam := r.URL.Query().Get("end")

	var startTimeUnix uint64
	var endTimeUnix uint64
	var errStartTime error = nil
	var errEndTime error = nil

	if len(startParam) == 0 {
		startTimeUnix = MinUint
	} else {
		startTimeUnix, errStartTime = strconv.ParseUint(startParam, 10, 64)
	}

	if len(endParam) == 0 {
		endTimeUnix = MaxUint
	} else {
		endTimeUnix, errEndTime = strconv.ParseUint(endParam, 10, 64)
	}

	if errStartTime != nil {
		http.Error(w, "Incorrect start time: "+errStartTime.Error(), http.StatusBadRequest)
		return
	}

	if errEndTime != nil {
		http.Error(w, "Incorrect end time: "+errEndTime.Error(), http.StatusBadRequest)
		return
	}

	if endTimeUnix < startTimeUnix {
		http.Error(w, "Invalid parameters", http.StatusBadRequest)
		return
	}

	responseJson, err := queryAssetVolume(assetId, startTimeUnix, endTimeUnix)

	if err != nil {
		log.Println("Error querying asset volume: " + err.Error())
		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
		return
	} else {
		w.Write(responseJson)
	}
}

func blockSwapsRequestHandler(w http.ResponseWriter, r *http.Request) {
}

func blockSwappedAssetsRequestHandler(w http.ResponseWriter, r *http.Request) {
}

func main() {
	router := chi.NewRouter()
	router.Get("/assets/{assetId}/pools", assetPoolsRequestHandler)
	router.Get("/assets/{assetId}/volume", assetVolumeRequestHandler)
	router.Get("/blocks/{blockId}/swaps", blockSwapsRequestHandler)
	router.Get("/blocks/{blockId}/assets-swapped", blockSwappedAssetsRequestHandler)

	fmt.Println("Listening on port " + ServerPort + "...")
	log.Fatal(http.ListenAndServe(":"+ServerPort, router))
}
