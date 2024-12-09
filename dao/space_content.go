package dao

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/letterScape/backend/dto"
	"gorm.io/gorm"
	"time"
)

type SpaceContent struct {
	ContentId    string    `json:"contentId" gorm:"column:content_id;primaryKey"`
	ChainId      string    `json:"chainId" gorm:"column:chain_id"`
	Author       string    `json:"author" gorm:"column:author"`
	Title        string    `json:"title" gorm:"column:title"`
	Resource     string    `json:"resource" gorm:"column:resource"`
	FavouriteNum int64     `json:"favouriteNum" gorm:"column:favourite_num"`
	Label        int64     `json:"label" gorm:"column:label"`
	IsShown      bool      `json:"isShown" gorm:"column:is_shown"`
	IsDeleted    bool      `json:"isDeleted" gorm:"column:is_deleted"`
	CreateTime   time.Time `json:"createTime" gorm:"column:create_time"`
	ModifyTime   time.Time `json:"modifyTime" gorm:"column:modify_time"`
	Readonly     bool      `json:"readonly" gorm:"column:readonly"`
}

func (content *SpaceContent) TableName() string {
	return "t_space_content"
}

func (content *SpaceContent) FindById(c *gin.Context, db *gorm.DB, contentId string) (*SpaceContent, error) {
	dbSession := db.WithContext(c)
	err := dbSession.Find(&content, "content_id = ?", contentId).Error
	if err != nil {
		return nil, err
	}
	return content, nil
}

func (content *SpaceContent) PageList(c *gin.Context, db *gorm.DB, params *dto.SpacePageInput) (*[]SpaceContent, int64, error) {
	var list []SpaceContent
	var count int64
	offset := (params.Page - 1) * params.PageSize
	dbSession := db.WithContext(c)
	dbSession = dbSession.Model(&SpaceContent{})
	if params.IsShown != nil {
		dbSession = dbSession.Where("is_shown = ?", params.IsShown)
	}
	if params.Author != nil && len(*(params.Author)) > 0 {
		dbSession = dbSession.Where("author = ? and chain_id = ?", params.Author, params.ChainId)
	}
	err := dbSession.Where("is_deleted = 0").Order("create_time desc").Limit(params.PageSize).Offset(offset).Find(&list).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, 0, err
	}
	err = dbSession.Table("t_space_content").Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	return &list, count, nil
}

func (content *SpaceContent) Save(c *gin.Context, db *gorm.DB) error {
	dbSession := db.WithContext(c)
	err := dbSession.Save(content).Error
	if err != nil {
		return err
	}
	return nil
}

func (content *SpaceContent) UpdateById(c *gin.Context, db *gorm.DB) error {
	dbSession := db.WithContext(c)
	err := dbSession.Model(&SpaceContent{}).Select("chain_id", "author", "title", "resource", "favourite_num", "label", "is_shown", "is_deleted", "readonly", "modify_time").Where("content_id = ?", content.ContentId).Updates(content).Error
	if err != nil {
		return err
	}
	return nil
}
