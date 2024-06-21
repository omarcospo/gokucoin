package main

import (
	"main/integration"
)

func main() {
	var operation string = "contracts/BRETTUSDTM"
	integration.PreSign = integration.Timestamp + "GET" + "/api/v1/" + operation
	var url string = "https://api-futures.kucoin.com/api/v1/"
	integration.RequestKuCoin(url + operation)
}
