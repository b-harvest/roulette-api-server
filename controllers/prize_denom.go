package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"roulette-api-server/models"
	"roulette-api-server/models/schema"
	"roulette-api-server/services"
	"roulette-api-server/types"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// 당첨품 종류 조회
func GetPrizeDenoms(c *gin.Context) {
	denoms := make([]schema.PrizeDenomRow, 0, 100)
	err := models.QueryPrizeDenoms(&denoms)
	if err != nil {
		fmt.Printf("%+v\n",err.Error())
		services.NotAcceptable(c, "fail " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		return
	}

	services.Success(c, nil, denoms)
}


// 당첨풍 종류 생성
func CreatePrizeDenom(c *gin.Context) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		services.BadRequest(c, "Bad Request " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		return
	}
	var req types.ReqTbCreatePrizeDenom
	if err = json.Unmarshal(jsonData, &req); err != nil {
		fmt.Println(err.Error())
		services.BadRequest(c, "Bad Request Unmarshal error: " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		return
	}

	// data handling
	denom := schema.PrizeDenomRow{
		Name:req.Name,
		Type:req.Type,
		IsActive:req.IsActive,
	}
	err = models.CreatePrizeDenom(&denom)

	// result
	if err != nil {
		fmt.Printf("%+v\n",err.Error())
		if strings.Contains(err.Error(),"1062") {
			services.NotAcceptable(c, "data already exists", err)
		} else {
			services.NotAcceptable(c, "fail " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		}
	} else {
		services.Success(c, nil, denom)
	}
}

// 특정 담청품 종류 조회
func GetPrizeDenom(c *gin.Context) {
	// 파라미터 조회
	strId := c.Param("id")
	reqId, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		services.BadRequest(c, "Bad Request id path parameter " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		return
	}
	denom := schema.PrizeDenomRow{
		PrizeDenomId: reqId,
	}
	err = models.QueryPrizeDenom(&denom)

	// result
	if err != nil {
		//if err.Error() == "record not found" {
		services.NotAcceptable(c, "fail " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
	} else {
		services.Success(c, nil, denom)
	}
}

// 데놈 정보 수정
func UpdatePrizeDenom(c *gin.Context) {
	// 파라미터 조회 -> body 조회 -> 언마샬
	strId := c.Param("id")
	reqId, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		services.BadRequest(c, "Bad Request id path parameter " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		return
	}
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
			services.BadRequest(c, "Bad Body request " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
			return
	}
	var req types.ReqTbUpdatePrizeDenom
	if err = json.Unmarshal(jsonData, &req); err != nil {
		services.BadRequest(c, "Bad Request Unmarshal error", err)
		return
	}

	// handler data
	denom := schema.PrizeDenomRow{
		PrizeDenomId: reqId,
		Name: req.Name,
		Type: req.Type,
		IsActive: req.IsActive,
		UpdatedAt: time.Now(),
	}
	err = models.UpdatePrizeDenomn(&denom)

	// result
	if err != nil {
		fmt.Printf("%+v\n",err.Error())
		if strings.Contains(err.Error(),"1062") {
			services.NotAcceptable(c, "something duplicated. already exists. fail " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		} else {
			services.NotAcceptable(c, "fail " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		}
	} else {
		services.Success(c, nil, denom)
	}
}


// 데놈 삭제
func DeletePrizeDenom(c *gin.Context) {
	// 파라미터 조회
	strId := c.Param("id")
	reqId, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		services.NotAcceptable(c, "Bad Request Id path parameter " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
		return
	}

	// handler data
	denom := schema.PrizeDenomRow{
		PrizeDenomId: reqId,
	}
	err = models.DeletePrizeDenom(&denom)

	// result
	if err != nil {
		services.NotAcceptable(c, "failed " + c.Request.Method + " " + c.Request.RequestURI + " : " + err.Error(), err)
	} else {
		services.Success(c, nil, denom)
	}
}




