package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"messari.io/uni-rest-api/config"
	"messari.io/uni-rest-api/handlers"
)

func main() {
	router := chi.NewRouter()
	router.Get("/assets/{assetId}/pools", handlers.AssetPoolsRequestHandler)
	router.Get("/assets/{assetId}/volume", handlers.AssetVolumeRequestHandler)
	router.Get("/blocks/{blockNumber}/swaps", handlers.BlockSwapsRequestHandler)
	router.Get("/blocks/{blockNumber}/swapped-assets", handlers.BlockSwappedAssetsRequestHandler)

	fmt.Println("Listening on port " + config.ServerPort + "...")
	log.Fatal(http.ListenAndServe(":"+config.ServerPort, router))
}
