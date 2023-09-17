package models

import (
	"fmt"
	"roulette-api-server/config"
	"roulette-api-server/models/schema"
	"roulette-api-server/types"
	"strings"
)

func QueryEventWalletConn(events *[]schema.EventwalletConnRow, promotionId string, addr string, startDate string, endDate string) error {
	var conditions []string
	if promotionId != "" {
		conditions = append(conditions, fmt.Sprintf("promotion_id = %s", promotionId))
	}
	if addr != "" {
		conditions = append(conditions, fmt.Sprintf("addr = '%s'", addr))
	}
	if startDate != "" && endDate != "" {
		conditions = append(conditions, fmt.Sprintf("(created_at BETWEEN '%s' AND '%s')", startDate, endDate))
	}

	sql := `
		SELECT *
		FROM event_wallet_conn
	`
	if len(conditions) != 0 {
		sql += "WHERE " + strings.Join(conditions, " AND ")
	}

	return config.DB.Raw(sql).Scan(events).Error
}

func QueryEventWalletConnCount(cnt *types.ResGetEventCount, promotionId string, addr string, startDate string, endDate string) error {
	var conditions []string
	if promotionId != "" {
		conditions = append(conditions, fmt.Sprintf("promotion_id = %s", promotionId))
	}
	if addr != "" {
		conditions = append(conditions, fmt.Sprintf("addr = '%s'", addr))
	}
	if startDate != "" && endDate != "" {
		conditions = append(conditions, fmt.Sprintf("(created_at BETWEEN '%s' AND '%s')", startDate, endDate))
	}

	sql := `
		SELECT COUNT(*) as cnt
		FROM event_wallet_conn
	`
	if len(conditions) != 0 {
		sql += "WHERE " + strings.Join(conditions, " AND ")
	}

	return config.DB.Raw(sql).Scan(cnt).Error
}

func CreateEventWalletConn(event *types.ReqPostEvent) error {
	err := config.DB.Table("event_wallet_conn").Create(&event).Error
	return err
}
