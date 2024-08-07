package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"roulette-api-server/config"
)

var api string

type WinningRequst struct {
	Address string `json:"address"`
	Amount  int    `json:"amount"`
}

func init() {
	api = fmt.Sprintf("http://%s:%d", config.Cfg.TPConf.Host, config.Cfg.TPConf.Port)
}

func IsDelegated(address string) (bool, error) {
	url := fmt.Sprintf("%s/checkaddr/%s", api, address)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return false, nil
	}

	return true, nil
}

func SendToken(address string, amount int) error {
	url := fmt.Sprintf("%s/winning", api)

	reqBytes, err := json.Marshal(WinningRequst{
		Address: address,
		Amount:  amount,
	})
	if err != nil {
		return err
	}
	reqBody := bytes.NewBuffer(reqBytes)

	req, err := http.NewRequest("POST", url, reqBody)
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
