package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"messari.io/uni-rest-api/common"
	"messari.io/uni-rest-api/graphql"
)

func AssetPoolsRequestHandler(w http.ResponseWriter, r *http.Request) {
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

func AssetVolumeRequestHandler(w http.ResponseWriter, r *http.Request) {
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

func BlockSwapsRequestHandler(w http.ResponseWriter, r *http.Request) {
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

func BlockSwappedAssetsRequestHandler(w http.ResponseWriter, r *http.Request) {
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
