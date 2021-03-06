package graphql

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/machinebox/graphql"
	"messari.io/uni-rest-api/common"
	"messari.io/uni-rest-api/config"
)

var graphqlClient = graphql.NewClient(config.UniswapV3Endpoint)

func QueryAssetPools(assetId string) ([]byte, error) {
	graphqlQuery := `{
		poolsWithToken0: pools(where: { token0: "` + assetId + `"}) {
			id
		}
		poolsWithToken1: pools(where: { token1: "` + assetId + `"}) {
			id
		}
	}`
	graphqlRequest := graphql.NewRequest(graphqlQuery)
	var allPools common.AllPools
	if err := graphqlClient.Run(context.Background(), graphqlRequest, &allPools); err != nil {
		return nil, err
	}
	combinedPools := append(allPools.PoolsWithToken0, allPools.PoolsWithToken1...)
	jsonResponse, jsonError := json.Marshal(combinedPools)
	return jsonResponse, jsonError
}

func QueryAssetVolume(assetId string, startTimeUnix uint64, endTimeUnix uint64) ([]byte, error) {
	totalVolume := 0.0

	for {
		startTimeQueryString := ""
		if startTimeUnix != common.MinUint {
			startTimeQueryString = " date_gt: " + strconv.FormatUint(startTimeUnix, 10) + " "
		}
		endTimeQueryString := ""
		if endTimeUnix != common.MaxUint {
			endTimeQueryString = " date_lt: " + strconv.FormatUint(endTimeUnix, 10) + " "
		}

		graphqlQuery := `{
			tokenDayDatas(first: ` + strconv.Itoa(config.GraphqlResultsPerPage) + `, orderBy: date, orderDirection: desc, where: {
				token: "` + assetId + `"
				` + startTimeQueryString + `
				` + endTimeQueryString + `
			}) {
				date
				volumeUSD
			}
		}`

		graphqlRequest := graphql.NewRequest(graphqlQuery)
		var tokenDayDataResult common.TokenDayDataResult
		if err := graphqlClient.Run(context.Background(), graphqlRequest, &tokenDayDataResult); err != nil {
			return nil, err
		}

		for i := 0; i < len(tokenDayDataResult.TokenDayData); i++ {
			totalVolume += tokenDayDataResult.TokenDayData[i].VolumeUSD
		}

		if len(tokenDayDataResult.TokenDayData) == config.GraphqlResultsPerPage {
			endTimeUnix = tokenDayDataResult.TokenDayData[config.GraphqlResultsPerPage-1].Date
		} else {
			break
		}
	}

	jsonResponse, jsonError := json.Marshal(map[string]float64{"TotalVolumeUSD": totalVolume})
	return jsonResponse, jsonError
}

func QueryBlockSwaps(blockNumber uint64) ([]byte, error) {
	lastTxQueryString := ""
	allSwaps := []string{}

	for {
		graphqlQuery := `{
			transactions(first: ` + strconv.Itoa(config.GraphqlResultsPerPage) + `, orderBy: id, orderDirection: desc, where: {
				blockNumber: ` + strconv.FormatUint(blockNumber, 10) + `
				` + lastTxQueryString + `
			}) {
				id
				swaps {
					id
				}
			}
		}`
		graphqlRequest := graphql.NewRequest(graphqlQuery)
		var transactionsResult common.TransactionsResult
		if err := graphqlClient.Run(context.Background(), graphqlRequest, &transactionsResult); err != nil {
			return nil, err
		}
		transactions := transactionsResult.Transactions

		for i := 0; i < len(transactions); i++ {
			tx := transactions[i]
			for j := 0; j < len(tx.Swaps); j++ {
				allSwaps = append(allSwaps, tx.Swaps[j].Id)
			}
		}

		if len(transactions) == config.GraphqlResultsPerPage {
			lastTxQueryString = " id_lt: \"" + transactions[len(transactions)-1].Id + "\" "
		} else {
			break
		}
	}
	jsonResponse, jsonError := json.Marshal(allSwaps)
	return jsonResponse, jsonError
}

func QueryBlockSwapsAssets(blockNumber uint64) ([]byte, error) {
	lastTxQueryString := ""
	allAssets := map[string]bool{}

	for {
		graphqlQuery := `{
			transactions(first: ` + strconv.Itoa(config.GraphqlResultsPerPage) + `, orderBy: id, orderDirection: desc, where: {
				blockNumber: ` + strconv.FormatUint(blockNumber, 10) + `
				` + lastTxQueryString + `
			}) {
				id
				swaps {
					token0 {
						id
					}
					token1 {
						id
					}
				}
			}
		}`
		graphqlRequest := graphql.NewRequest(graphqlQuery)
		var transactionsResult common.TransactionsResult
		if err := graphqlClient.Run(context.Background(), graphqlRequest, &transactionsResult); err != nil {
			return nil, err
		}
		transactions := transactionsResult.Transactions

		for i := 0; i < len(transactions); i++ {
			tx := transactions[i]
			for j := 0; j < len(tx.Swaps); j++ {
				allAssets[tx.Swaps[j].Token0.Id] = true
				allAssets[tx.Swaps[j].Token1.Id] = true
			}
		}

		if len(transactions) == config.GraphqlResultsPerPage {
			lastTxQueryString = " id_lt: \"" + transactions[len(transactions)-1].Id + "\" "
		} else {
			break
		}
	}

	jsonAssets := make([]string, len(allAssets))
	i := 0
	for k := range allAssets {
		jsonAssets[i] = k
		i++
	}
	jsonResponse, jsonError := json.Marshal(jsonAssets)
	return jsonResponse, jsonError
}
