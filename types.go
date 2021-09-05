package main

type Pool struct {
	Id string `json:"id"`
}

// Pools with asset as token0 and pools with asset as token1
type AllPools struct {
	PoolsWithToken0 []Pool `json:"poolsWithToken0"`
	PoolsWithToken1 []Pool `json:"poolsWithToken1"`
}
