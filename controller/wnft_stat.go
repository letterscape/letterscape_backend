package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/letterScape/backend/dto"
	"github.com/letterScape/backend/middleware"
	"github.com/letterScape/backend/services"
)

type WnftStatController struct{}

func WnftStatRegister(router *gin.RouterGroup) {
	wnftInfo := WnftStatController{}
	router.GET("/price", wnftInfo.Price)
	router.GET("/trade", wnftInfo.Trade)
	router.GET("/transaction", wnftInfo.Transaction)
}

func (stat *WnftStatController) Price(c *gin.Context) {
	input := &dto.DealPriceStatInput{}
	if err := input.BindingValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	service := &services.WNFTStatService{}
	priceStat, err := service.DealPriceStat(c, input)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	middleware.ResponseSuccess(c, priceStat)
	return
}

func (stat *WnftStatController) Trade(c *gin.Context) {
	input := &dto.TradePageInput{}
	if err := input.BindingValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	service := &services.WNFTStatService{}
	page, err := service.TradePage(c, input)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	middleware.ResponseSuccess(c, page)
	return
}

func (stat *WnftStatController) Transaction(c *gin.Context) {
	input := &dto.TransactionPageInput{}
	if err := input.BindingValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	service := &services.WNFTStatService{}
	page, err := service.TransactionPage(c, input)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	middleware.ResponseSuccess(c, page)
	return
}
