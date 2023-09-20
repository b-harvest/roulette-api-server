package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"roulette-api-server/config"
	"roulette-api-server/models"
	"roulette-api-server/models/schema"
	"time"
)


func PriceGetterHandler() {
	tick := time.NewTicker(10 * time.Minute) // block?
	// tick := time.NewTicker(5 * time.Second) // block?
	for {
		select {
		case <-tick.C:
			go GetPrices()
		}
	}
}

type DenomPriceResp struct {
	Result           string  `json:"result"`
	Data             []DenomPrices `json:"data"`
}

type DenomPrices struct {
	Denom           string  `json:"denom"`
	PriceOracle     float64 `json:"priceOracle"`
	UpdateTimestamp int64   `json:"updateTimestamp"`
}

func GetPrices() {
	if config.PriceUrl == "" {
		fmt.Println("Price 가져오는 URL 이 config.toml 에 추가되지 않음")
		return
	}

	resp, err := http.Get(config.PriceUrl)
	if err != nil {
		fmt.Printf("GetPrices err: %+v", err.Error())
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Body err: %+v", err.Error())
	}
	// fmt.Printf("%s\n", string(data))

	var denomPriceResp DenomPriceResp
	err = json.Unmarshal(data, &denomPriceResp)
	if err != nil {
		fmt.Printf("Unmarshal err: %+v", err.Error())
	}
	
	for _, v := range denomPriceResp.Data {
		fmt.Printf("%+v %+v\n",v.Denom,v.PriceOracle)
		// atom: ibc/C4CFF46FD6DE35CA4CF4CE031E643C8FDC9BA4B99AE598E9B0ED98FE3A2319F9
		// eth.axl: ibc/F1806958CA98757B91C3FA1573ECECD24F6FA3804F074A6977658914A49E65A3
		var row schema.PrizeDenomRow
		row.UsdPrice = v.PriceOracle
		switch v.Denom {
		case "ucre":
			row.Name = "cre"
			models.UpdatePrizeDenomByName(&row)
			row.Name = "CRE"
			models.UpdatePrizeDenomByName(&row)
		case "ubcre":
			row.Name = "bcre"
			models.UpdatePrizeDenomByName(&row)
			row.Name = "BCRE"
			models.UpdatePrizeDenomByName(&row)
		case "ibc/C4CFF46FD6DE35CA4CF4CE031E643C8FDC9BA4B99AE598E9B0ED98FE3A2319F9":
			row.Name = "atom"
			models.UpdatePrizeDenomByName(&row)
			row.Name = "ATOM"
			models.UpdatePrizeDenomByName(&row)
		case "ibc/F1806958CA98757B91C3FA1573ECECD24F6FA3804F074A6977658914A49E65A3":
			row.Name = "eth"
			models.UpdatePrizeDenomByName(&row)
			row.Name = "ETH"
			models.UpdatePrizeDenomByName(&row)
		}
	}

	return
}