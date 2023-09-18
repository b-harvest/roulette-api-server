package services

import "strings"

func GetAddressType(addr string) string {
	if strings.HasPrefix(addr, "cre") || strings.HasPrefix(addr, "CRE") {
		return "CRE"
	}
	// 이더리움 주소 형식 파악 필요
	if strings.HasPrefix(addr, "0x") {
		return "ETH"
	}

	return ""
}