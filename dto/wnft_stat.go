package dto

import (
	"github.com/gin-gonic/gin"
	"github.com/letterScape/backend/public"
)

type TradePageInput struct {
	PageSize int    `form:"pageSize" json:"pageSize" validate:"" example:"10"`
	Page     int    `form:"page" json:"page" validate:"required" example:"1"`
	WnftId   string `form:"wnftId" json:"wnftId" validate:"required"`
}

func (params *TradePageInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

type TransactionPageInput struct {
	PageSize int    `form:"pageSize" json:"pageSize" validate:"" example:"10"`
	Page     int    `form:"page" json:"page" validate:"required" example:"1"`
	WnftId   string `form:"wnftId" json:"wnftId" validate:"required"`
}

func (params *TransactionPageInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

type DealPriceStatInput struct {
	StartTime string `form:"startTime" json:"startTime" validate:"required"`
	EndTime   string `form:"endTime" json:"endTime" validate:"required"`
	WnftId    string `form:"wnftId" json:"wnftId" validate:"required"`
}

func (params *DealPriceStatInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}
