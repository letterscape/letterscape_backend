package dto

import (
	"github.com/gin-gonic/gin"
	"github.com/letterScape/backend/public"
)

type SpacePageInput struct {
	PageSize int     `form:"pageSize" json:"pageSize" validate:"" example:"10"`
	Page     int     `form:"page" json:"page" validate:"required" example:"1"`
	ChainId  string  `form:"chainId" json:"chainId" validate:""`
	Author   *string `form:"author" json:"author" validate:""`
	Label    string  `form:"label" json:"label" validate:""`
	IsShown  *bool   `form:"isShown" json:"isShown" validate:""`
}

func (params *SpacePageInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

type SpaceContentSaveInput struct {
	ContentId    string `json:"contentId" validate:"required"`
	ChainId      string `json:"chainId" validate:"required"`
	Author       string `json:"author" validate:"required"`
	Title        string `json:"title" validate:"required"`
	Resource     string `json:"resource" validate:"required"`
	FavouriteNum int64  `json:"favouriteNum" validate:""`
	Label        int64  `json:"label" validate:""`
	IsShown      bool   `json:"isShown" validate:""`
	IsDeleted    bool   `json:"isDeleted" validate:""`
}

func (params *SpaceContentSaveInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

type ContentPublishInput struct {
	ContentId string `json:"contentId" validate:"required"`
}

func (params *ContentPublishInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

type UploadContentInput struct {
	ContentId string `json:"contentId" validate:"required"`
	Content   string `json:"content" validate:"required"`
}

func (params *UploadContentInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}
