package middlewares

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"roulette-api-server/config"
	"time"
)

var thirdPartyAPI string

type WinningRequst struct {
	Address string `json:"address"`
	Amount  int    `json:"amount"`
}

type (
	IsSomethingResponse struct {
		Status  string      `json:"status"`
		Address string      `json:"address"`
		Amount  json.Number `json:"amount"`
		}	
 	
	IsSomethingReturnType struct {
		Status  string `json:"status"`
		Address string `json:"address"`
		Amount  float64  `json:"amount"`
	}
)

func init() {
	thirdPartyAPI = fmt.Sprintf("http://%s:%d", config.Cfg.TPConf.Host, config.Cfg.TPConf.Port)
}

// if return is nil then not delegated
func IsDelegated(address string) (*IsSomethingReturnType, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	url := fmt.Sprintf("%s/checkaddr/%s", thirdPartyAPI, address)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, nil
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var res IsSomethingResponse
	err = json.Unmarshal(bodyBytes, &res)
	if err != nil {
		return nil, err
	}

	amount, err := res.Amount.Float64()
	if err != nil {
		return nil, err
	}

	output := IsSomethingReturnType{
		Status:  res.Status,
		Address: res.Address,
		Amount:  amount,
	}

	return &output, nil
}

// if return is nil then not delegated
func IsYeetardHave(address string) (*IsSomethingReturnType, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	url := fmt.Sprintf("%s/haveYeetard/%s", thirdPartyAPI, address)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, nil
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var res IsSomethingResponse
	err = json.Unmarshal(bodyBytes, &res)
	if err != nil {
		return nil, err
	}

	amount, err := res.Amount.Float64()
	if err != nil {
		return nil, err
	}

	output := IsSomethingReturnType{
		Status:  res.Status,
		Address: res.Address,
		Amount:  amount,
	}

	return &output, nil
}

func SendToken(address string, amount int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	url := fmt.Sprintf("%s/winning", thirdPartyAPI)

	reqBytes, err := json.Marshal(WinningRequst{
		Address: address,
		Amount:  amount,
	})
	if err != nil {
		return err
	}
	reqBody := bytes.NewBuffer(reqBytes)

	req, err := http.NewRequestWithContext(ctx, "POST", url, reqBody)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf(string(bodyBytes))
	}

	return nil
}
