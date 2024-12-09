package vo

import (
	"time"
)

type WnftDetailVO struct {
	WnftId         string    `form:"wnftId" json:"wnftId" validate:"required"`
	TokenId        string    `form:"tokenId" json:"tokenId" validate:"required"`
	ChainId        string    `form:"chainId" json:"chainId" validate:"required"`
	Owner          string    `form:"owner" json:"owner" validate:"required"`
	Price          string    `form:"price" json:"price" validate:"required"`
	Interval       int       `form:"interval" json:"interval" validate:"required,gt=1"`
	Deadline       int       `form:"deadline" json:"deadline" validate:"required"`
	IsPaid         bool      `form:"isPaid" json:"isPaid" validate:""`
	IsListed       bool      `form:"isListed" json:"isListed" validate:""`
	IsExpired      bool      `form:"isExpired" json:"isExpired" validate:""`
	IsBurnt        bool      `form:"isBurnt" json:"isBurnt" validate:""`
	Title          string    `form:"title" json:"title" validate:"required"`
	Content        string    `form:"content" json:"content" validate:""`
	Hostname       string    `form:"hostname" json:"hostname" validate:"required"`
	OriginUri      string    `form:"originUri" json:"originUri" validate:"required"`
	IsTitlePicture bool      `form:"isTitlePicture" json:"isTitlePicture" validate:""`
	CreateTime     time.Time `form:"createTime" json:"createTime" validate:"required"`
	ModifyTime     time.Time `form:"modifyTime" json:"modifyTime" validate:""`
}

type WnftDetailVOList struct {
	List  *[]WnftDetailVO `form:"list" json:"list" validate:""`
	Total int64           `form:"page" json:"total" validate:"required"`
}
