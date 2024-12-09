package dao

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/letterScape/backend/constants/transaction"
	"github.com/letterScape/backend/dto"
	"gorm.io/gorm"
	"time"
)

type WnftBase struct {
	WnftId        string    `json:"wnftId" gorm:"column:wnft_id;primaryKey"`
	TokenId       string    `json:"tokenId" gorm:"column:token_id;index"`
	ChainId       string    `json:"chainId" gorm:"column:chain_id"`
	Owner         string    `json:"owner" gorm:"column:owner;index"`
	Price         string    `json:"price" gorm:"column:price"`
	LastDealPrice string    `json:"lastDealPrice" gorm:"column:last_deal_price"`
	Interval      int64     `json:"interval" gorm:"column:interval"`
	Deadline      int64     `json:"deadline" gorm:"column:deadline"`
	IsPaid        bool      `json:"isPaid" gorm:"column:is_paid"`
	IsListed      bool      `json:"isListed" gorm:"column:is_listed"`
	IsExpired     bool      `json:"isExpired" gorm:"column:is_expired"`
	IsBurnt       bool      `json:"isBurnt" gorm:"column:is_burnt"`
	CreateTime    time.Time `json:"createTime" gorm:"column:create_time"`
	ModifyTime    time.Time `json:"modifyTime" gorm:"column:modify_time"`
	Readonly      bool      `json:"readonly" gorm:"column:readonly"`
}

type WnftDetailDTO struct {
	WnftId     string    `json:"wnftId" gorm:"column:wnft_id;primaryKey" validate:"required"`
	TokenId    string    `form:"tokenId" json:"tokenId" gorm:"column:token_id" validate:"required"`
	ChainId    string    `json:"chainId" gorm:"column:chain_id" validate:"required"`
	Owner      string    `form:"owner" json:"owner" gorm:"column:owner;index" validate:"required"`
	Price      string    `form:"price" json:"price" gorm:"column:price" validate:"required"`
	Interval   int       `form:"interval" json:"interval" gorm:"column:interval" validate:"required,gt=1"`
	Deadline   int64     `form:"deadline" json:"deadline" gorm:"column:deadline" validate:"required"`
	IsPaid     bool      `form:"isPaid" json:"isPaid" gorm:"column:is_paid" validate:""`
	IsListed   bool      `form:"isListed" json:"isListed" gorm:"column:is_listed" validate:""`
	IsExpired  bool      `form:"isExpired" json:"isExpired" gorm:"column:is_expired" validate:""`
	IsBurnt    bool      `form:"isBurnt" json:"isBurnt" gorm:"column:is_burnt" validate:""`
	Title      string    `form:"title" json:"title" gorm:"column:title" validate:"required"`
	Content    string    `form:"content" json:"content" gorm:"column:content" validate:""`
	Hostname   string    `form:"hostname" json:"hostname" gorm:"column:hostname" validate:"required"`
	OriginUri  string    `form:"originUri" json:"originUri" gorm:"column:origin_uri" validate:"required"`
	CreateTime time.Time `form:"createTime" json:"createTime" gorm:"column:create_time" validate:"required"`
	ModifyTime time.Time `form:"modifyTime" json:"modifyTime" gorm:"column:modify_time" validate:""`
}

type WnftTxDTO struct {
	WnftId        string               `json:"wnftId" gorm:"column:wnft_id;primaryKey"`
	TokenId       string               `json:"tokenId" gorm:"column:token_id"`
	ChainId       string               `json:"chainId" gorm:"column:chain_id"`
	Owner         string               `json:"owner" gorm:"column:owner;index"`
	Price         string               `json:"price" gorm:"column:price"`
	LastDealPrice string               `json:"lastDealPrice" gorm:"column:last_deal_price"`
	Interval      int                  `json:"interval" gorm:"column:interval"`
	Deadline      int                  `json:"deadline" gorm:"column:deadline"`
	IsPaid        bool                 `json:"isPaid" gorm:"column:is_paid"`
	IsListed      bool                 `json:"isListed" gorm:"column:is_listed"`
	IsExpired     bool                 `json:"isExpired" gorm:"column:is_expired"`
	IsBurnt       bool                 `json:"isBurnt" gorm:"column:is_burnt"`
	CreateTime    time.Time            `json:"createTime" gorm:"column:create_time"`
	ModifyTime    time.Time            `json:"modifyTime" gorm:"column:modify_time"`
	Readonly      bool                 `json:"readonly" gorm:"column:readonly"`
	TxId          string               `json:"txId" gorm:"column:tx_id;primaryKey"`
	TxObject      transaction.TxObject `json:"txObject" gorm:"column:tx_object"`
	TxHash        string               `json:"txHash" gorm:"column:tx_hash"`
	TxStatus      transaction.TxStatus `json:"txStatus" gorm:"column:tx_status;type:enum('none','pending','mined','failed','success');default:'none'"`
	TxType        int                  `json:"txType" gorm:"column:tx_type"`
	Detail        string               `json:"detail" gorm:"column:detail"`
}

func (wnft *WnftBase) TableName() string {
	return "t_wnft_base"
}

func (wnft *WnftBase) FindById(c *gin.Context, db *gorm.DB, wnftId string) (*WnftBase, error) {
	var info *WnftBase
	dbSession := db.WithContext(c)
	err := dbSession.Find(&info, "wnft_id = ?", wnftId).Error
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (wnft *WnftBase) GetDetailById(c *gin.Context, tx *gorm.DB, wnftId string) (detail *WnftDetailDTO, err error) {
	dbSession := tx.WithContext(c)
	err = dbSession.Model(&WnftBase{}).Select("t_wnft_base.*, t_wnft_detail.title, t_wnft_detail.content, t_wnft_detail.hostname, t_wnft_detail.origin_uri, t_wnft_detail.is_title_picture").Joins("left join t_wnft_detail on t_wnft_base.wnft_id = t_wnft_detail.wnft_id").Where("t_wnft_base.wnft_id = ?", wnftId).Scan(&detail).Error
	if err != nil {
		return nil, err
	}
	return detail, nil
}

func (wnft *WnftBase) PageList(c *gin.Context, db *gorm.DB, params *dto.WnftPageInput) (*[]WnftDetailDTO, int64, error) {
	var list []WnftDetailDTO
	var count int64
	offset := (params.Page - 1) * params.PageSize
	dbSession := db.WithContext(c)
	dbSession = dbSession.Model(&WnftBase{}).Select("t_wnft_base.*, t_wnft_detail.title, t_wnft_detail.content, t_wnft_detail.hostname, t_wnft_detail.origin_uri").Joins("left join t_wnft_detail on t_wnft_base.wnft_id = t_wnft_detail.wnft_id")
	if params.IsListed != nil {
		dbSession = dbSession.Where("is_listed = ?", params.IsListed)
	}
	if params.IsBurnt != nil {
		dbSession = dbSession.Where("is_burnt = ?", params.IsBurnt)
	}
	if params.Owner != nil && len(*(params.Owner)) > 0 {
		dbSession = dbSession.Where("owner = ? and chain_id = ?", params.Owner, params.ChainId)
	}
	err := dbSession.Order("t_wnft_base.create_time desc").Limit(params.PageSize).Offset(offset).Find(&list).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, 0, err
	}
	err = dbSession.Table("t_wnft_base").Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	return &list, count, nil
}

func (wnft *WnftBase) Save(c *gin.Context, db *gorm.DB) error {
	dbSession := db.WithContext(c)
	err := dbSession.Save(wnft).Error
	if err != nil {
		return err
	}
	return nil
}

func (wnft *WnftBase) UpdateByIdWithWriteable(c *gin.Context, db *gorm.DB) error {
	dbSession := db.WithContext(c)
	var info *WnftBase
	err := dbSession.Find(&info, "wnft_id = ? and owner = ?", wnft.WnftId, wnft.Owner).Error
	if err != nil {
		return err
	}
	if info == nil {
		return errors.New("data doesn't exist")
	}
	if info.Readonly == true {
		return errors.New("readonly data cannot be modified")
	}
	err = dbSession.Where("wnft_id = ?", wnft.WnftId).Updates(wnft).Error
	if err != nil {
		return err
	}
	return nil
}

func (wnft *WnftBase) UpdateById(c *gin.Context, db *gorm.DB) error {
	dbSession := db.WithContext(c)
	err := dbSession.Model(&WnftBase{}).Select("owner", "price", "last_deal_price", "interval", "deadline", "is_paid", "is_listed", "is_expired", "is_burnt", "readonly", "modify_time").Where("wnft_id = ?", wnft.WnftId).Updates(wnft).Error
	//err := dbSession.Model(&WnftBase{}).Select("t_wnft_base.*").Where("wnft_id = ?", wnft.WnftId).Updates(wnft).Error
	if err != nil {
		return err
	}
	return nil
}
