package common

const MinUint uint64 = 0
const MaxUint uint64 = ^uint64(0)

type Pool struct {
	Id string `json:"id"`
}

// Pools with asset as token0 and pools with asset as token1
type AllPools struct {
	PoolsWithToken0 []Pool `json:"poolsWithToken0"`
	PoolsWithToken1 []Pool `json:"poolsWithToken1"`
}

type TokenDayData struct {
	Date      uint64  `json:"date"`
	VolumeUSD float64 `json:"volumeUSD,string"`
}

type TokenDayDataResult struct {
	TokenDayData []TokenDayData `json:"tokenDayDatas"`
}

type Token struct {
	Id string `json:"id"`
}

type Swap struct {
	Id     string `json:"id"`
	Token0 Token  `json:"token0"`
	Token1 Token  `json:"token1"`
}

type Transaction struct {
	Id    string `json:"id"`
	Swaps []Swap `json:"swaps"`
}

type TransactionsResult struct {
	Transactions []Transaction `json:"transactions"`
}
