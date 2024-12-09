package vo

import "time"

type SpaceContentVO struct {
	ContentId    string    `json:"contentId" gorm:"column:content_id;primaryKey"`
	ChainId      string    `json:"chainId" gorm:"column:chain_id"`
	Author       string    `json:"author" gorm:"column:author"`
	Title        string    `json:"title" gorm:"column:title"`
	Resource     string    `json:"resource" gorm:"column:resource"`
	FavouriteNum string    `json:"favouriteNum" gorm:"column:favourite_num"`
	Label        int64     `json:"label" gorm:"column:label"`
	IsShown      bool      `json:"isShown" gorm:"column:is_shown"`
	IsDeleted    bool      `json:"isDeleted" gorm:"column:is_deleted"`
	CreateTime   time.Time `json:"createTime" gorm:"column:create_time"`
	ModifyTime   time.Time `json:"modifyTime" gorm:"column:modify_time"`
}

type SpaceContentVOList struct {
	List  *[]SpaceContentVO `form:"list" json:"list" validate:""`
	Total int64             `form:"page" json:"total" validate:"required"`
}
