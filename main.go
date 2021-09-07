package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"messari.io/uni-rest-api/common"
	"messari.io/uni-rest-api/config"
	"messari.io/uni-rest-api/graphql"
)

func assetPoolsRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	assetId := chi.URLParam(r, "assetId")

	if !common.ValidAddress(assetId) {
		http.Error(w, "Invalid asset id", http.StatusBadRequest)
		return
	}

	responseJson, err := graphql.QueryAssetPools(assetId)

	if err != nil {
		log.Println("Error querying asset pools: " + err.Error())
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

	if !common.ValidAddress(assetId) {
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
		startTimeUnix = common.MinUint
	} else {
		startTimeUnix, errStartTime = strconv.ParseUint(startParam, 10, 64)
	}

	if len(endParam) == 0 {
		endTimeUnix = common.MaxUint
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

	responseJson, err := graphql.QueryAssetVolume(assetId, startTimeUnix, endTimeUnix)

	if err != nil {
		log.Println("Error querying asset volume: " + err.Error())
		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
		return
	} else {
		w.Write(responseJson)
	}
}

func blockSwapsRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	blockNumberParam := chi.URLParam(r, "blockNumber")

	blockNumber, blockNumberErr := strconv.ParseUint(blockNumberParam, 10, 64)

	if blockNumberErr != nil {
		http.Error(w, "Invalid block id "+blockNumberParam, http.StatusBadRequest)
		return
	}

	responseJson, err := graphql.QueryBlockSwaps(blockNumber)

	if err != nil {
		log.Println("Error querying block swaps: " + err.Error())
		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
		return
	} else {
		w.Write(responseJson)
	}
}

func blockSwappedAssetsRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	blockNumberParam := chi.URLParam(r, "blockNumber")

	blockNumber, blockNumberErr := strconv.ParseUint(blockNumberParam, 10, 64)

	if blockNumberErr != nil {
		http.Error(w, "Invalid block id "+blockNumberParam, http.StatusBadRequest)
		return
	}

	responseJson, err := graphql.QueryBlockSwapsAssets(blockNumber)

	if err != nil {
		log.Println("Error querying block swap assets: " + err.Error())
		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
		return
	} else {
		w.Write(responseJson)
	}
}

func main() {
	router := chi.NewRouter()
	router.Get("/assets/{assetId}/pools", assetPoolsRequestHandler)
	router.Get("/assets/{assetId}/volume", assetVolumeRequestHandler)
	router.Get("/blocks/{blockNumber}/swaps", blockSwapsRequestHandler)
	router.Get("/blocks/{blockNumber}/swapped-assets", blockSwappedAssetsRequestHandler)

	fmt.Println("Listening on port " + config.ServerPort + "...")
	log.Fatal(http.ListenAndServe(":"+config.ServerPort, router))
}
