package dto

import (
	"github.com/letterScape/backend/constants/transaction"
)

type TxRecordListInput struct {
	TxObject transaction.TxObject `form:"txObject" json:"txObject" validate:"required"`
	TxStatus transaction.TxStatus `form:"txStatus" json:"txStatus" validate:"required"`
	TxType   int8                 `form:"txType" json:"txType" validate:"required" default:"0"`
	Size     int                  `form:"size" json:"size" validate:"required" default:"100"`
}
