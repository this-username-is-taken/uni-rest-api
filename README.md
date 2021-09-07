# Uniswap V3 REST API

## Run

1. Go to project root directory
2. Start server with `go run .`

## Test

1. Go to project root directory
2. Start server with `go run .`
3. Run test with `go test -v`

## Usage

### Retrieve pools that include a particular asset

`/assets/{assetId}/pools`

Description: given an asset ID, this endpoint will return a list of pool IDs that contains this particular asset, either as token0 or token1.

Example Request:

`/assets/0xd533a949740bb3306d119cc777fa900ba034cd52/pools`

Example Response:

`[{"id":"0x07b1c12be0d62fe548a2b4b025ab7a5ca8def21e"},{"id":"0x8e26e2fc8140280fba3e34bfdca7fc1102c1ae04"},{"id":"0x4c83a7f819a5c37d64b4c5a2f8238ea082fa1f4e"},{"id":"0x620cd19eae24fb8a02df908bb71b81b6e3aa1ccc"},{"id":"0x645c3a387b8633df1d4d71ca4b50d27233bcb887"},{"id":"0x919fa96e88d67499339577fa202345436bcdaf79"},{"id":"0x9445bd19767f73dcae6f2de90e6cd31192f62589"},{"id":"0xcbeb7da1ec121fc37dde2bc9010f3a4001e1ebcb"}]`

### Retrieve total volume of a particular asset given a time period

Description: given an asset ID, this endpoint will return the total volume (in USD) of a particular asset within startTime and endTime. Note that the volume is aggregated on a per-day basis.

`/assets/{assetId}/volume?start={startTime}&end={endTime}`

Parameters:

`startTime`: (Optional) start of the period in Unix timestamp
`endTime`: (Optional) end of the period in Unix timestamp

Example Request:

`/assets/0xd533a949740bb3306d119cc777fa900ba034cd52/volume?start=1627000000&end=1628000000`

Example Response:

`{"TotalVolumeUSD":3106715.5626834165}`

### Retrieve all swaps in a block

Description: given a block number (note that this is a block number not a block hash), this endpoint will return the swap ID of all swaps occurred in this block.

`/blocks/{blockNumber}/swaps`

Example Request:

`/blocks/12738079/swaps`

Example Response:

`["0x8d5a3c7a2293aac0520b497cd6036ba5b225f7430411bf8b803cbf3a591cd8d0#22695","0x0bc17005eaa7084146ac1c002332072d6c4b650fd3ec0728089954c87b38e440#26623"]`

### Retrieve all assets swapped in a block

Description: given a block number (note that this is a block number not a block hash), this endpoint will return all asset IDs of that were swapped in this block.

`/blocks/{blockNumber}/swapped-assets`

Example Request:

`/blocks/12738079/swapped-assets`

Example Response:

`["0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2","0x7d1afa7b718fb893db30a3abc0cfc608aacfebb0","0x27c70cd1946795b66be9d954418546998b546634"]`

## Improvements

Given more time, here are the things I would work on to improvement:

- [ ] The total volume is currently aggregated by day and may not be precise enough for many applications. This could be improved to be aggregated on a per-block basis with more complicated GraphQL queries.
- [ ] Current tests only test a small number of scenarios and do not cover many edge cases.
- [ ] Better design on the returned fields. Currently only the bare minimal data fields (just the IDs) are returned for most endpoints. This probably can be designed in a much better way with customer/PM input.
- [ ] Logging and error handling.
