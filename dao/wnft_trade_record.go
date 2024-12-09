package dao

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/letterScape/backend/constants/sql"
	"github.com/letterScape/backend/dto"
	"gorm.io/gorm"
	"time"
)

type WnftTradeRecord struct {
	TradeId    string    `json:"tradeId" gorm:"column:trade_id;primaryKey"`
	WnftId     string    `json:"wnftId" gorm:"column:wnft_id"`
	ChainId    string    `json:"chainId" gorm:"column:chain_id"`
	Seller     string    `json:"seller" gorm:"column:seller"`
	Buyer      string    `json:"buyer" gorm:"column:buyer"`
	DealPrice  string    `json:"dealPrice" gorm:"column:deal_price"`
	CreateTime time.Time `json:"createTime" gorm:"column:create_time"`
}

type QueryListInput struct {
	TradeId     string      `json:"tradeId"`
	WnftId      string      `json:"wnftId"`
	ChainId     string      `json:"chainId"`
	Seller      string      `json:"seller"`
	Buyer       string      `json:"buyer"`
	StartTime   string      `json:"startTime"`
	EndTime     string      `json:"endTime"`
	OrderByTime sql.OrderBy `json:"orderByTime"`
}

func (record *WnftTradeRecord) TableName() string {
	return "t_wnft_trade_record"
}

func (record *WnftTradeRecord) FindById(c *gin.Context, db *gorm.DB, wnftId string) (*WnftTradeRecord, error) {
	var tradeRecord *WnftTradeRecord
	dbSession := db.WithContext(c)
	err := dbSession.Find(tradeRecord, "wnft_id = ?", wnftId).Error
	if err != nil {
		return nil, err
	}
	return tradeRecord, nil
}

func (record *WnftTradeRecord) List(c *gin.Context, db *gorm.DB, input *QueryListInput) (*[]WnftTradeRecord, error) {
	var list []WnftTradeRecord
	dbSession := db.WithContext(c).Model(&WnftTradeRecord{})

	if input.TradeId != "" {
		dbSession = dbSession.Where("trade_id = ?", input.TradeId)
	}
	if input.WnftId != "" {
		dbSession = dbSession.Where("wnft_id = ?", input.WnftId)
	}
	if input.ChainId != "" {
		dbSession = dbSession.Where("chain_id = ?", input.ChainId)
	}
	if input.Seller != "" {
		dbSession = dbSession.Where("seller = ?", input.Seller)
	}
	if input.Buyer != "" {
		dbSession = dbSession.Where("buyer = ?", input.Buyer)
	}
	if input.StartTime != "" {
		dbSession = dbSession.Where("create_time >= ?", input.StartTime)
	}
	if input.EndTime != "" {
		dbSession = dbSession.Where("create_time < ?", input.EndTime)
	}
	if input.OrderByTime != "" {
		dbSession = dbSession.Order("create_time " + input.OrderByTime)
	}
	err := dbSession.Find(&list).Error
	if err != nil {
		return nil, err
	}
	return &list, nil
}

func (record *WnftTradeRecord) Save(c *gin.Context, db *gorm.DB) error {
	dbSession := db.WithContext(c)
	err := dbSession.Save(record).Error
	if err != nil {
		return err
	}
	return nil
}

func (record *WnftTradeRecord) Page(c *gin.Context, db *gorm.DB, input *dto.TradePageInput) (*[]WnftTradeRecord, int64, error) {
	var records []WnftTradeRecord
	var total int64
	offset := (input.Page - 1) * input.PageSize
	dbSession := db.WithContext(c)
	err := dbSession.Model(&WnftTradeRecord{}).Where("wnft_id = ?", input.WnftId).Order("create_time desc").Limit(input.PageSize).Offset(offset).Find(&records).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, 0, err
	}
	err = dbSession.Table("t_wnft_trade_record").Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	return &records, total, nil
}
