package dto

import (
	"github.com/gin-gonic/gin"
	"github.com/letterScape/backend/public"
)

type ResourceInput struct {
	ResourceId string `json:"resourceId" validate:"required"`
	TypeId     string `json:"typeId" validate:"required"`
	Url        string `json:"url" validate:""`
	Text       string `json:"text" validate:""`
}

func (params *ResourceInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

type FindResourceInput struct {
	Fp      string `form:"fp" json:"fp" validate:"required"`
	ChainId string `form:"chainId" json:"chainId" validate:"required"`
}

func (params *FindResourceInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}
