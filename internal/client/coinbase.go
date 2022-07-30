package client

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type Rate struct {
	Data struct {
		Base     string `json:"base"`
		Currency string `json:"currency"`
		Amount   string `json:"amount"`
	} `json:"data"`
}

type CoinBaseClient struct {
	url string
}

func NewCoinBaseClient() *CoinBaseClient {
	return &CoinBaseClient{
		url: "https://api.coinbase.com/v2/prices/spot?currency=UAH",
	}
}

func (c CoinBaseClient) GetRate() (int, error) {
	response, err := http.Get(c.url)
	if err != nil {
		return 0, err
	}
	var rawBody Rate
	err = json.NewDecoder(response.Body).Decode(&rawBody)
	if err != nil {
		return 0, err
	}
	before, _, _ := strings.Cut(rawBody.Data.Amount, ".")
	price, err := strconv.Atoi(before)
	return price, nil
}
