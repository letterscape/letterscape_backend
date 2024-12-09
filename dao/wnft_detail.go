package dao

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

type WnftDetail struct {
	WnftId         string    `json:"wnftId" gorm:"column:wnft_id;primaryKey"`
	Title          string    `json:"title" gorm:"column:title"`
	Content        string    `json:"content" gorm:"column:content"`
	Hostname       string    `json:"hostname" gorm:"column:hostname"`
	IsTitlePicture bool      `json:"isTitlePicture" gorm:"column:is_title_picture"`
	OriginUri      string    `json:"originUri" gorm:"column:origin_uri"`
	CreateTime     time.Time `json:"createTime" gorm:"column:create_time"`
	ModifyTime     time.Time `json:"modifyTime" gorm:"column:modify_time"`
}

func (detail *WnftDetail) TableName() string {
	return "t_wnft_detail"
}

func (detail *WnftDetail) FindById(c *gin.Context, db *gorm.DB, wnftId string) (*WnftDetail, error) {
	var info *WnftDetail
	dbSession := db.WithContext(c)
	err := dbSession.Find(&info, "wnft_id = ?", wnftId).Error
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (detail *WnftDetail) Save(c *gin.Context, db *gorm.DB) error {
	dbSession := db.WithContext(c)
	err := dbSession.Save(detail).Error
	if err != nil {
		return err
	}
	return nil
}

func (detail *WnftDetail) UpdateById(c *gin.Context, db *gorm.DB) error {
	dbSession := db.WithContext(c)
	err := dbSession.Where("wnft_id = ?", detail.WnftId).Updates(detail).Error
	if err != nil {
		return err
	}
	return nil
}
