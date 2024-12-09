package services

import (
	"errors"
	"github.com/e421083458/golang_common/lib"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/letterScape/backend/chain"
	"github.com/letterScape/backend/constants/transaction"
	"github.com/letterScape/backend/dao"
	"github.com/letterScape/backend/dto"
	"github.com/letterScape/backend/middleware"
	"github.com/letterScape/backend/utils"
	"github.com/letterScape/backend/vo"
	"gorm.io/gorm"
	"log"
	"strings"
	"time"
)

type WNFTInfoService struct{}

func (service *WNFTInfoService) Info(c *gin.Context) {

}

func (service *WNFTInfoService) Detail(c *gin.Context, wnftId string) (detailVO *vo.WnftDetailVO, err error) {
	db, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	detail, err := (&dao.WnftBase{}).GetDetailById(c, db, wnftId)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return nil, err
	}

	detailVO = &vo.WnftDetailVO{}
	err = copier.Copy(&detailVO, &detail)
	if err != nil {
		return nil, err
	}

	detailVO.Price, err = EnlargeSymbol(c, db, detailVO.ChainId, detailVO.Price)
	if err != nil {
		return nil, err
	}

	return detailVO, nil
}

func (service *WNFTInfoService) Page(c *gin.Context, pageInput *dto.WnftPageInput) (*vo.WnftDetailVOList, error) {
	db, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return nil, err
	}
	wnftList, total, err := (&dao.WnftBase{}).PageList(c, db, pageInput)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return nil, err
	}

	var voList []vo.WnftDetailVO
	size := len(*wnftList)
	for i := 0; i < size; i++ {
		wnft := (*wnftList)[i]
		detailVO := &vo.WnftDetailVO{}
		err := copier.Copy(&detailVO, &wnft)
		if err != nil {
			return nil, err
		}
		detailVO.Price, err = EnlargeSymbol(c, db, detailVO.ChainId, detailVO.Price)
		if err != nil {
			return nil, err
		}
		voList = append(voList, *detailVO)
	}
	pageList := &vo.WnftDetailVOList{List: &voList, Total: total}
	return pageList, nil
}

func (service *WNFTInfoService) Mint(c *gin.Context, input *dto.MintWnftInput) error {
	chainContext := &chain.Context{}
	chainContext.SetChainOpt(chain.Mapping[input.ChainId])
	status, err := chainContext.QueryTxStatus(common.HexToHash(input.TxHash))
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return err
	}

	db, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return err
	}

	tx := db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	// there is only one tokenId in one chain
	wnftId := input.TokenId + "chid" + input.ChainId

	price, err := ShrinkSymbol(c, db, input.ChainId, input.Price)
	if err != nil {
		return err
	}

	wnftBase := &dao.WnftBase{
		WnftId:        wnftId,
		TokenId:       input.TokenId,
		ChainId:       input.ChainId,
		Owner:         input.Owner,
		Price:         price,
		LastDealPrice: "0",
		Interval:      input.Interval,
		Deadline:      time.Now().Unix() + input.Interval*3600,
		IsPaid:        true,
		IsListed:      false,
		IsExpired:     false,
		IsBurnt:       false,
		CreateTime:    time.Now(),
		ModifyTime:    time.Now(),
		Readonly:      status != transaction.Success && status != transaction.Failed,
	}

	// it could be inserted when nft first created or be updated a burnt nft
	if err := wnftBase.Save(c, db); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2002, err)
		return err
	}

	// query nft whether has existed
	wnftDetailQuery, err := (&dao.WnftDetail{}).FindById(c, db, wnftId)
	if err != nil {
		tx.Rollback()
		return err
	}

	wnftDetail := &dao.WnftDetail{
		WnftId:     wnftId,
		Title:      input.Title,
		Content:    input.Content,
		CreateTime: time.Now(),
		ModifyTime: time.Now(),
	}

	if len(wnftDetailQuery.WnftId) > 0 {
		wnftDetail.Hostname = wnftDetailQuery.Hostname
		wnftDetail.OriginUri = wnftDetailQuery.OriginUri
	} else {
		if len(input.Hostname) == 0 || len(input.OriginUri) == 0 {
			tx.Rollback()
			return errors.New("uri cannot be empty")
		}
		wnftDetail.Hostname = input.Hostname
		wnftDetail.OriginUri = input.OriginUri
	}

	if err := wnftDetail.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2002, err)
		return err
	}

	// update old record to not latest
	err = db.Model(&dao.TxRecord{}).Where("is_latest = 1 and detail_id = ?", wnftId).Updates(map[string]interface{}{"is_latest": 0, "modify_time": time.Now()}).Error
	if err != nil {
		tx.Rollback()
		log.Fatal("wnft tx update failed!", wnftId, err)
	}

	txRecord := &dao.TxRecord{
		TxId:       strings.ReplaceAll(uuid.New().String(), "-", ""),
		ChainId:    input.ChainId,
		TxObject:   transaction.Wnft,
		TxHash:     input.TxHash,
		TxStatus:   status,
		TxType:     transaction.MINT,
		DetailId:   wnftId,
		CreateTime: time.Now(),
		ModifyTime: time.Now(),
		IsLatest:   true,
	}
	if err := txRecord.Save(c, db); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2002, err)
		return err
	}

	return tx.Commit().Error
}

func (service *WNFTInfoService) List(c *gin.Context, input *dto.WnftListInput) error {
	chainContext := &chain.Context{}
	chainContext.SetChainOpt(&chain.EthereumOpts{})
	status, err := chainContext.QueryTxStatus(common.HexToHash(input.TxHash))
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return err
	}

	db, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return err
	}

	tx := db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	wnftQuery, err := (&dao.WnftBase{}).FindById(c, db, input.WnftId)
	if err != nil {
		tx.Rollback()
		return err
	}

	wnftBase := &dao.WnftBase{
		WnftId:     input.WnftId,
		Owner:      wnftQuery.Owner,
		IsListed:   true,
		ModifyTime: time.Now(),
		Readonly:   status != transaction.Success && status != transaction.Failed,
	}

	if err := wnftBase.UpdateByIdWithWriteable(c, db); err != nil {
		tx.Rollback()
		return err
	}

	// update old record to not latest
	err = db.Model(&dao.TxRecord{}).Where("is_latest = 1 and detail_id = ?", input.WnftId).Updates(map[string]interface{}{"is_latest": 0, "modify_time": time.Now()}).Error
	if err != nil {
		tx.Rollback()
		log.Fatal("wnft tx update failed!", input.WnftId, err)
	}

	txRecord := &dao.TxRecord{
		TxId:       strings.ReplaceAll(uuid.New().String(), "-", ""),
		ChainId:    wnftQuery.ChainId,
		TxObject:   transaction.Wnft,
		TxHash:     input.TxHash,
		TxStatus:   status,
		TxType:     transaction.LIST,
		DetailId:   input.WnftId,
		CreateTime: time.Now(),
		ModifyTime: time.Now(),
		IsLatest:   true,
	}
	if err := txRecord.Save(c, db); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2002, err)
		return err
	}

	return tx.Commit().Error
}

func (service *WNFTInfoService) Buy(c *gin.Context, input *dto.BuyWnftInput) error {
	chainContext := &chain.Context{}
	chainContext.SetChainOpt(&chain.EthereumOpts{})
	txData, err := chainContext.QueryTx(common.HexToHash(input.TxHash))
	if err != nil {
		return err
	}

	// less than one minute
	//if txData.Time.Sub(time.Now()).Abs() > 60 {
	//	return errors.New("duration of tx is too large")
	//}

	db, err := lib.GetGormPool("default")
	if err != nil {
		return err
	}

	tx := db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	// update wbftBase
	wnftBase, err := (&dao.WnftBase{}).FindById(c, db, input.WnftId)
	if err != nil {
		tx.Rollback()
		return err
	}

	tokenIdBig, err := utils.ToBigInt(txData.Data[0])
	if err != nil {
		tx.Rollback()
		return err
	}
	tokenId := hexutil.EncodeBig(tokenIdBig)
	if tokenId != wnftBase.TokenId {
		tx.Rollback()
		return errors.New("wrong tx data")
	}

	if txData.TxStatus == transaction.Success {
		wnft, err := chainContext.GetWNFT(tokenId)
		if err != nil {
			return err
		}
		wnftBase.Owner = wnft.Owner.String()
		wnftBase.Price = wnft.Price.String()
		wnftBase.LastDealPrice = wnft.LastDealPrice.String()
		wnftBase.Deadline = wnft.Deadline.Int64()
		wnftBase.IsPaid = true
		wnftBase.Readonly = false
	} else {
		wnftBase.Owner = txData.From
		lastDealPriceBig, err := utils.ToBigInt(txData.Data[1])
		if err != nil {
			tx.Rollback()
			return err
		}
		wnftBase.LastDealPrice = lastDealPriceBig.String()
		price, err := utils.ToBigInt(txData.Data[2])
		if err != nil {
			tx.Rollback()
			return err
		}
		wnftBase.Price = price.String()
		wnftBase.Readonly = true
	}
	err = wnftBase.UpdateByIdWithWriteable(c, db)
	if err != nil {
		tx.Rollback()
		return err
	}

	// insert wnftTradeRecord
	tradeRecord := &dao.WnftTradeRecord{
		TradeId:    strings.ReplaceAll(uuid.New().String(), "-", ""),
		WnftId:     input.WnftId,
		Seller:     txData.From,
		Buyer:      txData.To,
		DealPrice:  wnftBase.LastDealPrice,
		ChainId:    wnftBase.ChainId,
		CreateTime: time.Now(),
	}
	if err := tradeRecord.Save(c, db); err != nil {
		tx.Rollback()
		return err
	}

	// update old record to not latest
	err = db.Model(&dao.TxRecord{}).Where("is_latest = 1 and detail_id = ?", input.WnftId).Updates(map[string]interface{}{"is_latest": 0, "modify_time": time.Now()}).Error
	if err != nil {
		tx.Rollback()
		log.Fatal("wnft tx update failed!", input.WnftId, err)
	}

	// insert txRecord
	txRecord := &dao.TxRecord{
		TxId:       strings.ReplaceAll(uuid.New().String(), "-", ""),
		ChainId:    wnftBase.ChainId,
		TxObject:   transaction.Wnft,
		TxHash:     input.TxHash,
		TxStatus:   txData.TxStatus,
		TxType:     transaction.BUY,
		DetailId:   input.WnftId,
		CreateTime: time.Now(),
		ModifyTime: time.Now(),
		IsLatest:   true,
	}
	if err := txRecord.Save(c, db); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (service *WNFTInfoService) PayHoldfee(c *gin.Context, input *dto.PayHoldfeeInput) error {
	db, err := lib.GetGormPool("default")
	if err != nil {
		return err
	}

	wnftQuery, err := (&dao.WnftBase{}).FindById(c, db, input.WnftId)
	if err != nil {
		return err
	}

	if wnftQuery.IsBurnt || wnftQuery.IsExpired {
		return errors.New("cannot pay holdfee any more")
	}

	chainContext := &chain.Context{}
	chainContext.SetChainOpt(chain.Mapping[wnftQuery.ChainId])
	status, err := chainContext.QueryTxStatus(common.HexToHash(input.TxHash))
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return err
	}

	tx := db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	newDeadline := wnftQuery.Deadline + wnftQuery.Interval*3600

	wnftBase := &dao.WnftBase{
		WnftId:     input.WnftId,
		Deadline:   newDeadline,
		IsPaid:     true,
		ModifyTime: time.Now(),
		Readonly:   status != transaction.Success && status != transaction.Failed,
	}

	if err := wnftBase.UpdateByIdWithWriteable(c, db); err != nil {
		tx.Rollback()
		return err
	}

	// update old record to not latest
	err = db.Model(&dao.TxRecord{}).Where("is_latest = 1 and detail_id = ?", input.WnftId).Updates(map[string]interface{}{"is_latest": 0, "modify_time": time.Now()}).Error
	if err != nil {
		tx.Rollback()
		log.Fatal("wnft tx update failed!", input.WnftId, err)
	}

	txRecord := &dao.TxRecord{
		TxId:       strings.ReplaceAll(uuid.New().String(), "-", ""),
		ChainId:    wnftQuery.ChainId,
		TxObject:   transaction.Wnft,
		TxHash:     input.TxHash,
		TxStatus:   status,
		TxType:     transaction.HOLDFEE,
		DetailId:   input.WnftId,
		CreateTime: time.Now(),
		ModifyTime: time.Now(),
		IsLatest:   true,
	}
	if err := txRecord.Save(c, db); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2002, err)
		return err
	}

	return tx.Commit().Error
}

func (service *WNFTInfoService) Burn(c *gin.Context, input *dto.BurnWnftInput) error {
	db, err := lib.GetGormPool("default")
	if err != nil {
		return err
	}

	wnftQuery, err := (&dao.WnftBase{}).FindById(c, db, input.WnftId)
	if err != nil {
		return err
	}

	chainContext := &chain.Context{}
	chainContext.SetChainOpt(chain.Mapping[wnftQuery.ChainId])
	status, err := chainContext.QueryTxStatus(common.HexToHash(input.TxHash))
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return err
	}

	tx := db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	wnftBase := &dao.WnftBase{
		WnftId:     input.WnftId,
		IsPaid:     false,
		IsListed:   false,
		IsBurnt:    true,
		ModifyTime: time.Now(),
		Readonly:   status != transaction.Success && status != transaction.Failed,
	}

	if err := wnftBase.UpdateByIdWithWriteable(c, db); err != nil {
		tx.Rollback()
		return err
	}

	// update old record to not latest
	err = db.Model(&dao.TxRecord{}).Where("is_latest = 1 and detail_id = ?", input.WnftId).Updates(map[string]interface{}{"is_latest": 0, "modify_time": time.Now()}).Error
	if err != nil {
		tx.Rollback()
		log.Fatal("wnft tx update failed!", input.WnftId, err)
	}

	txRecord := &dao.TxRecord{
		TxId:       strings.ReplaceAll(uuid.New().String(), "-", ""),
		ChainId:    wnftQuery.ChainId,
		TxObject:   transaction.Wnft,
		TxHash:     input.TxHash,
		TxStatus:   status,
		TxType:     transaction.BURN,
		DetailId:   input.WnftId,
		CreateTime: time.Now(),
		ModifyTime: time.Now(),
		IsLatest:   true,
	}
	if err := txRecord.Save(c, db); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2002, err)
		return err
	}

	return tx.Commit().Error
}

func (service *WNFTInfoService) UpdateDetail(c *gin.Context, input *dto.WnftDetailInput) error {
	db, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return err
	}

	detail := &dao.WnftDetail{
		WnftId:         input.WnftId,
		Title:          input.Title,
		Content:        input.Content,
		IsTitlePicture: input.IsTitlePicture,
		ModifyTime:     time.Now(),
	}
	err = detail.UpdateById(c, db)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return err
	}
	return nil
}

func (service *WNFTInfoService) PollTx(c *gin.Context) {
	log.Printf("polling tx status")

	db, err := lib.GetGormPool("default")
	if err != nil {
		log.Fatal("get gorm pool failed", err)
	}

	input := &dto.TxRecordListInput{
		TxStatus: transaction.Pending,
		TxObject: transaction.Wnft,
		Size:     100,
	}
	records, err := (&dao.TxRecord{}).List(c, db, input)
	if err != nil {
		log.Fatal("query tx records failed", err)
	}
	chainContext := &chain.Context{}
	chainContext.SetChainOpt(&chain.EthereumOpts{})
	for _, record := range *records {
		status, err := chainContext.QueryTxStatus(common.HexToHash(record.TxHash))
		if err != nil {
			log.Fatal("polling tx records failed", err)
		}
		if status == transaction.Success {
			updateTxStatus(c, db, record, transaction.Success)
		} else if status == transaction.Failed {
			updateTxStatus(c, db, record, transaction.Failed)
		}
	}
}

func updateTxStatus(c *gin.Context, db *gorm.DB, record dao.TxRecord, status transaction.TxStatus) {
	record.TxStatus = status
	record.ModifyTime = time.Now()
	err := record.UpdateById(c, db)
	if err != nil {
		log.Fatal("update tx status failed", err)
	}
}

// sync wnft readonly status
func (service *WNFTInfoService) SyncWnftStatus(c *gin.Context) error {
	log.Printf("Sync wnft readonly status")

	db, err := lib.GetGormPool("default")
	if err != nil {
		log.Fatal("get gorm pool failed", err)
	}
	dbSession := db.WithContext(c)
	var list []dao.WnftTxDTO
	size := 100
	// todo wnft匹配tx_record最新记录是否为完成
	err = dbSession.Model(&dao.WnftBase{}).Select("t_wnft_base.*, t_tx_record.*").Joins("left join t_tx_record on t_wnft_base.wnft_id = t_tx_record.detail_id and t_tx_record.is_latest = 1").Where("t_wnft_base.readonly = 1 and (t_tx_record.tx_status = 'success' or t_tx_record.tx_status = 'failed')").Limit(size).Scan(&list).Error
	if err != nil {
		log.Fatal(err.Error())
		return err
	}

	log.Printf("there are %v records need to be synced", len(list))

	if len(list) == 0 {
		return nil
	}

	var wnftIds []string
	for _, txDTO := range list {
		wnftIds = append(wnftIds, txDTO.WnftId)
	}

	// todo 事务处理，同步链上wnft信息
	for _, txDTO := range list {
		// if nft was burnt, there is no data when query wnft on the chain
		if txDTO.IsBurnt {
			syncBurntWnft(c, db, &txDTO)
		} else {
			syncNotBurntWnft(c, db, &txDTO)
		}
	}

	//dbSession.Model(&dao.WnftBase{}).Where("wnft_id in (?)", wnftIds).Updates(map[string]interface{}{"readonly": 0, "modify_time": time.Now()})

	log.Printf("there are %v records have been synced", len(list))

	return nil
}

func syncBurntWnft(c *gin.Context, db *gorm.DB, txDTO *dao.WnftTxDTO) {
	wnftBase := &dao.WnftBase{
		WnftId:     txDTO.WnftId,
		ModifyTime: time.Now(),
		Readonly:   false,
	}

	tx := db.Begin()
	dbSession := db.WithContext(c)

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := wnftBase.UpdateById(c, db); err != nil {
		tx.Rollback()
		log.Fatal("wnft update failed!", txDTO.WnftId, err)
	}
	err := dbSession.Model(&dao.TxRecord{}).Where("is_latest = 1 and detail_id = ?", txDTO.WnftId).Updates(map[string]interface{}{"is_latest": 0, "modify_time": time.Now()}).Error
	if err != nil {
		tx.Rollback()
		log.Fatal("wnft tx update failed!", txDTO.WnftId, err)
	}

	tx.Commit()
}

func syncNotBurntWnft(c *gin.Context, db *gorm.DB, txDTO *dao.WnftTxDTO) {
	chainContext := &chain.Context{}
	chainContext.SetChainOpt(chain.Mapping[txDTO.ChainId])
	wnft, err := chainContext.GetWNFT(txDTO.TokenId)
	if err != nil {
		log.Printf("query wnft:[%v] failed: %v", txDTO.WnftId, err)
		return
	}
	wnftBase := &dao.WnftBase{
		WnftId:        txDTO.WnftId,
		Owner:         wnft.Owner.Hex(),
		Price:         wnft.Price.String(),
		LastDealPrice: wnft.LastDealPrice.String(),
		Interval:      wnft.Interval.Int64(),
		Deadline:      wnft.Deadline.Int64(),
		IsPaid:        wnft.IsPaid,
		IsListed:      wnft.IsListed,
		IsExpired:     wnft.IsExpired,
		ModifyTime:    time.Now(),
		Readonly:      false,
	}

	tx := db.Begin()
	dbSession := db.WithContext(c)

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := wnftBase.UpdateById(c, db); err != nil {
		tx.Rollback()
		log.Fatal("wnft update failed!", txDTO.WnftId, err)
	}
	err = dbSession.Model(&dao.TxRecord{}).Where("is_latest = 1 and detail_id = ?", txDTO.WnftId).Updates(map[string]interface{}{"is_latest": 0, "modify_time": time.Now()}).Error
	if err != nil {
		tx.Rollback()
		log.Fatal("wnft tx update failed!", txDTO.WnftId, err)
	}

	tx.Commit()
}
