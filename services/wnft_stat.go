package services

import (
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"github.com/letterScape/backend/constants/sql"
	"github.com/letterScape/backend/constants/transaction"
	"github.com/letterScape/backend/dao"
	"github.com/letterScape/backend/dto"
	"github.com/letterScape/backend/vo"
)

type WNFTStatService struct{}

func (stat *WNFTStatService) TradePage(c *gin.Context, input *dto.TradePageInput) (*vo.TradePageVO, error) {
	db, err := lib.GetGormPool("default")
	if err != nil {
		return nil, err
	}

	var list []vo.TradeRecordVO

	records, total, err := (&dao.WnftTradeRecord{}).Page(c, db, input)
	if err != nil {
		return nil, err
	}

	for _, record := range *records {
		dealPrice, err := EnlargeSymbol(c, db, record.ChainId, record.DealPrice)
		if err != nil {
			return nil, err
		}
		recordVO := vo.TradeRecordVO{
			TradeId:   record.TradeId,
			WnftId:    record.WnftId,
			Buyer:     record.Buyer,
			Seller:    record.Seller,
			DealPrice: dealPrice,
			DealTime:  record.CreateTime.Format("2006-01-02 15:04:05"),
		}
		list = append(list, recordVO)
	}

	page := &vo.TradePageVO{List: &list, Total: total}
	return page, nil
}

func (stat *WNFTStatService) DealPriceStat(c *gin.Context, input *dto.DealPriceStatInput) (*[][]string, error) {
	db, err := lib.GetGormPool("default")
	if err != nil {
		return nil, err
	}

	query := &dao.QueryListInput{
		WnftId:      input.WnftId,
		StartTime:   input.StartTime,
		EndTime:     input.EndTime,
		OrderByTime: sql.ASC,
	}

	list, err := (&dao.WnftTradeRecord{}).List(c, db, query)
	if err != nil {
		return nil, err
	}

	if len(*list) == 0 {
		return &[][]string{}, nil
	}

	var priceStatList [][]string
	for _, item := range *list {
		dealPrice, err := EnlargeSymbol(c, db, item.ChainId, item.DealPrice)
		if err != nil {
			return nil, err
		}
		priceStat := []string{
			item.CreateTime.Format("2006-01-02 15:04:05"),
			dealPrice,
		}
		priceStatList = append(priceStatList, priceStat)
	}

	return &priceStatList, nil
}

func (stat *WNFTStatService) TransactionPage(c *gin.Context, input *dto.TransactionPageInput) (*vo.TransactionPageVO, error) {
	db, err := lib.GetGormPool("default")
	if err != nil {
		return nil, err
	}

	var list []vo.TransactionVO

	records, total, err := (&dao.TxRecord{}).Page(c, db, input)
	if err != nil {
		return nil, err
	}

	for _, record := range *records {
		recordVO := vo.TransactionVO{
			TxId:       record.TxId,
			DetailId:   record.DetailId,
			TxHash:     record.TxHash,
			TxStatus:   record.TxStatus,
			TxType:     transaction.Type2String(record.TxType),
			CreateTime: record.CreateTime.Format("2006-01-02 15:04:05"),
		}
		list = append(list, recordVO)
	}

	page := &vo.TransactionPageVO{List: &list, Total: total}
	return page, nil
}
