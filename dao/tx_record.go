package dao

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/letterScape/backend/constants/transaction"
	"github.com/letterScape/backend/dto"
	"gorm.io/gorm"
	"time"
)

type TxRecord struct {
	TxId       string               `json:"txId" gorm:"column:tx_id;primaryKey"`
	ChainId    string               `json:"chainId" gorm:"column:chain_id"`
	TxObject   transaction.TxObject `json:"txObject" gorm:"column:tx_object"`
	TxHash     string               `json:"txHash" gorm:"column:tx_hash"`
	TxStatus   transaction.TxStatus `json:"txStatus" gorm:"column:tx_status;type:enum('none','pending','mined','failed','success');default:'none'"`
	TxType     int                  `json:"txType" gorm:"column:tx_type"`
	DetailId   string               `json:"detailId" gorm:"column:detail_id"`
	CreateTime time.Time            `json:"createTime" gorm:"column:create_time"`
	ModifyTime time.Time            `json:"modifyTime" gorm:"column:modify_time"`
	IsLatest   bool                 `json:"isLatest" gorm:"column:is_latest"`
}

func (txRecord *TxRecord) TableName() string {
	return "t_tx_record"
}

func (txRecord *TxRecord) FindById(c *gin.Context, tx *gorm.DB, txId string) (*TxRecord, error) {
	var record *TxRecord
	dbSession := tx.WithContext(c)
	err := dbSession.Where("tx_id = ?", txId).First(&record).Error
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (txRecord *TxRecord) List(c *gin.Context, tx *gorm.DB, input *dto.TxRecordListInput) (*[]TxRecord, error) {
	var list []TxRecord
	dbSession := tx.WithContext(c)
	dbSession = dbSession.Model(&TxRecord{}).Where("tx_status = ?", input.TxStatus).Where("tx_object = ?", input.TxObject)
	if input.TxType != 0 {
		dbSession = dbSession.Where("tx_type = ?", input.TxType)
	}
	err := dbSession.Order("create_time").Limit(input.Size).Find(&list).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return &list, nil
}

func (txRecord *TxRecord) Save(c *gin.Context, tx *gorm.DB) error {
	dbSession := tx.WithContext(c)
	err := dbSession.Save(txRecord).Error
	if err != nil {
		return err
	}
	return nil
}

func (txRecord *TxRecord) UpdateById(c *gin.Context, tx *gorm.DB) error {
	dbSession := tx.WithContext(c)
	err := dbSession.Where("tx_id = ?", txRecord.TxId).Updates(txRecord).Error
	if err != nil {
		return err
	}
	return nil
}

func (txRecord *TxRecord) Page(c *gin.Context, db *gorm.DB, input *dto.TransactionPageInput) (*[]TxRecord, int64, error) {
	var records []TxRecord
	var total int64
	offset := (input.Page - 1) * input.PageSize
	dbSession := db.WithContext(c)
	err := dbSession.Model(&TxRecord{}).Where("detail_id = ?", input.WnftId).Order("create_time desc").Limit(input.PageSize).Offset(offset).Find(&records).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, 0, err
	}
	err = dbSession.Table("t_tx_record").Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	return &records, total, nil
}
