package middlewares

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"roulette-api-server/config"
	"roulette-api-server/models/schema"
	"roulette-api-server/types"
	"time"
)

var goldGameAPI string

func init() {
	goldGameAPI = fmt.Sprintf("http://%s:%d", config.Cfg.TPConf.Host, config.Cfg.TPConf.Port)
}

func StartGoldGame(order *schema.OrderRow) (*types.ResStartGame, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	url := fmt.Sprintf("%s/game-mgmt/start", goldGameAPI)

	reqBytes, err := json.Marshal(order)
	if err != nil {
		return nil, err
	}
	reqBody := bytes.NewBuffer(reqBytes)

	req, err := http.NewRequestWithContext(ctx, "POST", url, reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("failed to start gold game")
	}

	resBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var res types.ResStartGame
	err = json.Unmarshal(resBytes, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
