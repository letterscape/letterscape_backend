package dao

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

type WnftFeeRecord struct {
	FeeId      string    `json:"feeId" gorm:"column:fee_id;primaryKey"`
	WnftId     string    `json:"wnftId" gorm:"column:wnft_id"`
	FeeType    string    `json:"wnftFeeType" gorm:"column:wnft_fee_type"`
	Payer      string    `json:"payer" gorm:"column:payer"`
	Amt        uint64    `json:"amt" gorm:"column:amt"`
	Unit       string    `json:"unit" gorm:"column:unit"`
	CreateTime time.Time `json:"createTime" gorm:"column:create_time"`
}

func (record *WnftFeeRecord) TableName() string {
	return "t_wnft_fee_record"
}

func (record *WnftFeeRecord) FindById(c *gin.Context, tx *gorm.DB, wnftId string) (*WnftFeeRecord, error) {
	var feeRecord *WnftFeeRecord
	dbSession := tx.WithContext(c)
	err := dbSession.Find(feeRecord, "wnft_id = ?", wnftId).Error
	if err != nil {
		return nil, err
	}
	return feeRecord, nil
}

func (record *WnftFeeRecord) Save(c *gin.Context, tx *gorm.DB) error {
	dbSession := tx.WithContext(c)
	err := dbSession.Save(record).Error
	if err != nil {
		return err
	}
	return nil
}
