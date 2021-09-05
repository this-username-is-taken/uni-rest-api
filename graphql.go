package main

import (
	"context"
	"encoding/json"

	"github.com/machinebox/graphql"
)

func QueryAssetPools(assetId string) ([]byte, error) {
	graphqlClient := graphql.NewClient(UniswapV3Endpoint)
	graphqlRequest := graphql.NewRequest(`{
			poolsWithToken0: pools(where: { token0: "` + assetId + `"}) {
				id
			}
			poolsWithToken1: pools(where: { token1: "` + assetId + `"}) {
				id
			}
		}`)
	var allPools AllPools
	if err := graphqlClient.Run(context.Background(), graphqlRequest, &allPools); err != nil {
		return nil, err
	}
	combinedPools := append(allPools.PoolsWithToken0, allPools.PoolsWithToken1...)
	jsonResponse, _ := json.Marshal(combinedPools)
	return jsonResponse, nil
}
