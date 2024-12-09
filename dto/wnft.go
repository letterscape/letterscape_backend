package dto

import (
	"github.com/gin-gonic/gin"
	"github.com/letterScape/backend/public"
)

type MintWnftInput struct {
	TokenId   string `form:"tokenId" json:"tokenId" validate:"required"`
	ChainId   string `form:"chainId" json:"chainId" validate:"required"`
	Owner     string `form:"owner" json:"owner" validate:"required"`
	Price     string `form:"price" json:"price" validate:"required"`
	Interval  int64  `form:"interval" json:"interval" validate:"required,gt=0"`
	Deadline  int64  `form:"deadline" json:"deadline" validate:""`
	Title     string `form:"title" json:"title" validate:"required"`
	Content   string `form:"content" json:"content" validate:""`
	Hostname  string `form:"hostname" json:"hostname" validate:""`
	OriginUri string `form:"originUri" json:"originUri" validate:""`
	TxHash    string `form:"txHash" json:"txHash" validate:"required"`
}

func (params *MintWnftInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

type WnftPageInput struct {
	PageSize int     `form:"pageSize" json:"pageSize" validate:"" example:"10"`
	Page     int     `form:"page" json:"page" validate:"required" example:"1"`
	IsListed *bool   `form:"isListed" json:"isListed" validate:"" example:"true"`
	IsBurnt  *bool   `form:"isBurnt" json:"isBurnt" validate:"" example:"false"`
	ChainId  string  `form:"chainId" json:"chainId" validate:"required"`
	Owner    *string `form:"owner" json:"owner" validate:""`
}

func (params *WnftPageInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

type WnftDetailInput struct {
	WnftId         string `form:"wnftId" json:"wnftId" validate:"required"`
	Title          string `form:"title" json:"title" validate:"required"`
	Content        string `form:"content" json:"content" validate:""`
	Hostname       string `form:"hostname" json:"hostname" validate:"required"`
	OriginUri      string `form:"originUri" json:"originUri" validate:"required"`
	IsTitlePicture bool   `form:"isTitlePicture" json:"isTitlePicture" validate:""`
}

func (params *WnftDetailInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

type WnftListInput struct {
	WnftId string `form:"wnftId" json:"wnftId" validate:"required"`
	Owner  string `form:"owner" json:"owner" validate:"required"`
	TxHash string `form:"txHash" json:"txHash" validate:"required"`
}

func (params *WnftListInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

type BuyWnftInput struct {
	WnftId string `form:"wnftId" json:"wnftId" validate:"required"`
	TxHash string `form:"txHash" json:"txHash" validate:"required"`
}

func (params *BuyWnftInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

type PayHoldfeeInput struct {
	WnftId string `form:"wnftId" json:"wnftId" validate:"required"`
	TxHash string `form:"txHash" json:"txHash" validate:"required"`
}

func (params *PayHoldfeeInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

type BurnWnftInput struct {
	WnftId string `form:"wnftId" json:"wnftId" validate:"required"`
	TxHash string `form:"txHash" json:"txHash" validate:"required"`
}

func (params *BurnWnftInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

type FetchResourceInput struct {
	Fp      string `form:"fp" json:"fp" validate:"required"`
	ChainId string `form:"chainId" json:"chainId" validate:"required"`
}

func (params *FetchResourceInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}
