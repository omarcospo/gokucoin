package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Api struct {
	Key    string
	Secret string
	Pass   string
}

func ReadApiJson() Api {
	path := os.Getenv("KCAPI")
	file, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
	}

	var Api Api
	json.Unmarshal([]byte(file), &Api)

	return Api
}

var (
	apiKey    = ReadApiJson().Key
	apiSecret = ReadApiJson().Secret
	apiPass   = ReadApiJson().Pass
)

var (
	timestamp = strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	preSign   = timestamp + "GET" + "/api/v1/accounts"
)

func MakePass() (signature string, passphrase string) {
	mac := hmac.New(sha256.New, []byte(apiSecret))
	mac.Write([]byte(preSign))
	signature = base64.StdEncoding.EncodeToString(mac.Sum(nil))

	macPassphrase := hmac.New(sha256.New, []byte(apiSecret))
	macPassphrase.Write([]byte(apiPass))
	passphrase = base64.StdEncoding.EncodeToString(macPassphrase.Sum(nil))

	return signature, passphrase
}

func RequestKuCoin(url string) {
	sign, _ := MakePass()
	_, pass := MakePass()

	// Create the request
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	// Add headers
	req.Header.Add("KC-API-KEY", apiKey)
	req.Header.Add("KC-API-SIGN", sign)
	req.Header.Add("KC-API-TIMESTAMP", timestamp)
	req.Header.Add("KC-API-PASSPHRASE", pass)
	req.Header.Add("KC-API-KEY-VERSION", "3")

	// Send the request
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer res.Body.Close()

	// Read and print the response
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	type Data struct {
		Symbol         string
		MarkPrice      float64
		LastTradePrice float64
		IndexPrice     float64
		VolumeOf24h    float64
		LowPrice       float64
		HighPrice      float64
	}

	type Response struct {
		Data Data
	}

	var Coin Response

	json.Unmarshal([]byte(body), &Coin)

	fmt.Println(Coin.Data)
}

func main() {
	var operation string = "contracts/BRETTUSDTM"
	preSign = timestamp + "GET" + "/api/v1/" + operation
	var url string = "https://api-futures.kucoin.com/api/v1/"
	RequestKuCoin(url + operation)
}
