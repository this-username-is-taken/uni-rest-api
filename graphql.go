package main

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/machinebox/graphql"
)

func queryAssetPools(assetId string) ([]byte, error) {
	graphqlClient := graphql.NewClient(UniswapV3Endpoint)
	graphqlQuery := `{
		poolsWithToken0: pools(where: { token0: "` + assetId + `"}) {
			id
		}
		poolsWithToken1: pools(where: { token1: "` + assetId + `"}) {
			id
		}
	}`
	graphqlRequest := graphql.NewRequest(graphqlQuery)
	var allPools AllPools
	if err := graphqlClient.Run(context.Background(), graphqlRequest, &allPools); err != nil {
		return nil, err
	}
	combinedPools := append(allPools.PoolsWithToken0, allPools.PoolsWithToken1...)
	jsonResponse, jsonError := json.Marshal(combinedPools)
	return jsonResponse, jsonError
}

func queryAssetVolume(assetId string, startTimeUnix uint64, endTimeUnix uint64) ([]byte, error) {
	const resultsPerPage = 100
	totalVolume := 0.0
	graphqlClient := graphql.NewClient(UniswapV3Endpoint)

	for {
		startTimeQueryString := ""
		if startTimeUnix != MinUint {
			startTimeQueryString = " date_gt: " + strconv.FormatUint(startTimeUnix, 10) + " "
		}
		endTimeQueryString := ""
		if endTimeUnix != MaxUint {
			endTimeQueryString = " date_lt: " + strconv.FormatUint(endTimeUnix, 10) + " "
		}

		graphqlQuery := `{
			tokenDayDatas(first: ` + strconv.Itoa(resultsPerPage) + `, orderBy: date, orderDirection: desc, where: {
				token: "` + assetId + `"
				` + startTimeQueryString + `
				` + endTimeQueryString + `
			}) {
				date
				volumeUSD
			}
		}`

		graphqlRequest := graphql.NewRequest(graphqlQuery)
		var tokenDayDataResult TokenDayDataResult
		if err := graphqlClient.Run(context.Background(), graphqlRequest, &tokenDayDataResult); err != nil {
			return nil, err
		}

		for i := 0; i < len(tokenDayDataResult.TokenDayData); i++ {
			totalVolume += tokenDayDataResult.TokenDayData[i].VolumeUSD
		}

		if len(tokenDayDataResult.TokenDayData) == resultsPerPage {
			endTimeUnix = tokenDayDataResult.TokenDayData[resultsPerPage-1].Date
		} else {
			break
		}
	}

	jsonResponse, jsonError := json.Marshal(map[string]float64{"TotalVolumeUSD": totalVolume})
	return jsonResponse, jsonError
}
