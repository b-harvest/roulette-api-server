package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"roulette-api-server/models"
	"roulette-api-server/models/schema"
	"roulette-api-server/services"
	"roulette-api-server/types"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// accounts 조회
func GetAccounts(c *gin.Context) {
	accs := make([]schema.AccountRow, 0, 100)
	err := models.QueryAccounts(&accs)
	if err != nil {
		fmt.Printf("%+v\n",err.Error())
		services.NotAcceptable(c, "fail " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		return
	}

	services.Success(c, nil, accs)
}


// account 생성
func CreateAccount(c *gin.Context) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		services.BadRequest(c, "Bad Request " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		return
	}
	var req types.ReqTbCreateAccount
	if err = json.Unmarshal(jsonData, &req); err != nil {
		fmt.Println(string(jsonData))
		fmt.Println(err.Error())
		services.BadRequest(c, "Bad Request Unmarshal error: " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		return
	}

	// data handling
	// acc := schema.AccountRow{
	// 	Addr:       req.Addr,
	// 	TicketAmount:      req.TicketAmount,
	// 	AdminMemo:     req.AdminMemo,
	// 	Type:             req.Type,
	// 	IsBlacklisted:      req.IsBlacklisted,
	// 	// LastLoginAt: time.Time{},
	// 	LastLoginAt: time.Now(),
	// }
	err = models.CreateAccount(&req)

	// result
	if err != nil {
		fmt.Printf("%+v\n",err.Error())
		if strings.Contains(err.Error(),"1062") {
			services.NotAcceptable(c, "data already exists", err)
		} else {
			services.NotAcceptable(c, "fail " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		}
	} else {
		services.Success(c, nil, req)
	}
}

// 특정 Account 조회
func GetAccount(c *gin.Context) {
	// 파라미터 조회
	strId := c.Param("addr")
	acc := schema.AccountRow{
		Addr: strId,
	}
	err := models.QueryAccount(&acc)

	// result
	if err != nil {
		//if err.Error() == "record not found" {
		services.NotAcceptable(c, "fail " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
	} else {
		services.Success(c, nil, acc)
	}
}

// Account 정보 수정
func UpdateAccount(c *gin.Context) {
	// 파라미터 조회 -> body 조회 -> 언마샬
	strId := c.Param("addr")
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
			services.BadRequest(c, "Bad Body request " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
			return
	}
	var req types.ReqTbUpdateAccount
	if err = json.Unmarshal(jsonData, &req); err != nil {
		services.BadRequest(c, "Bad Request Unmarshal error", err)
		return
	}

	// handler data
	acc := schema.AccountRow{
		Addr: strId,
		TicketAmount: req.TicketAmount,
		AdminMemo: req.AdminMemo,
		Type: req.Type,
		IsBlacklisted: req.IsBlacklisted,
		UpdatedAt: time.Now(),
	}
	err = models.UpdateAccount(&acc)

	// result
	if err != nil {
		fmt.Printf("%+v\n",err.Error())
		if strings.Contains(err.Error(),"1062") {
			services.NotAcceptable(c, "something duplicated. already exists. fail " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		} else {
			services.NotAcceptable(c, "fail " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		}
	} else {
		services.Success(c, nil, acc)
	}
}


// Account 삭제
func DeleteAccount(c *gin.Context) {
	// 파라미터 조회
	strId := c.Param("addr")

	// handler data
	acc := schema.AccountRow{
		Addr: strId,
	}
	err := models.DeleteAccount(&acc)

	// result
	if err != nil {
		services.NotAcceptable(c, "failed " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
	} else {
		services.Success(c, nil, acc)
	}
}




