package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/ipfs/boxo/files"
	"github.com/letterScape/backend/chain"
	"github.com/letterScape/backend/dto"
	"github.com/letterScape/backend/middleware"
	"github.com/letterScape/backend/services"
	"github.com/letterScape/backend/utils"
	"io"
	"log"
	"mime/multipart"
	"net/http"
)

type WnftInfoController struct{}

func WnftInfoRegister(router *gin.RouterGroup) {
	wnftInfo := WnftInfoController{}
	router.GET("/page", wnftInfo.Page)
	router.GET("/info", wnftInfo.Info)
	router.GET("/detail", wnftInfo.Detail)
	router.POST("/mint", wnftInfo.Mint).OPTIONS("/mint", wnftInfo.Mint)
	router.POST("/list", wnftInfo.List).OPTIONS("/list", wnftInfo.List)
	router.POST("/buy", wnftInfo.Buy).OPTIONS("/buy", wnftInfo.Buy)
	router.POST("/holdfee", wnftInfo.PayHoldfee).OPTIONS("/holdfee", wnftInfo.PayHoldfee)
	router.POST("/burn", wnftInfo.Burn).OPTIONS("/burn", wnftInfo.Burn)
	router.POST("/update", wnftInfo.Update)
	router.POST("/updateDetail", wnftInfo.UpdateDetail).OPTIONS("/updateDetail", wnftInfo.UpdateDetail)
	router.POST("/upload", wnftInfo.UploadResource).OPTIONS("/upload", wnftInfo.UploadResource)
	router.GET("/fetch", wnftInfo.FetchResource)
}

func (wnftInfo *WnftInfoController) Page(context *gin.Context) {
	pageInput := &dto.WnftPageInput{}
	if err := pageInput.BindingValidParams(context); err != nil {
		middleware.ResponseError(context, 2001, err)
		return
	}
	if pageInput.PageSize == 0 {
		pageInput.PageSize = 10
	}
	if pageInput.Page == 0 {
		pageInput.Page = 1
	}

	service := &services.WNFTInfoService{}
	wnftList, err := service.Page(context, pageInput)
	if err != nil {
		middleware.ResponseError(context, 2002, err)
		return
	}

	middleware.ResponseSuccess(context, wnftList)
	return
}

func (wnftInfo *WnftInfoController) Info(context *gin.Context) {

}

func (wnftInfo *WnftInfoController) Detail(context *gin.Context) {
	wnftId, exists := context.GetQuery("wnftId")
	if !exists {
		middleware.ResponseError(context, 2001, errors.New("params not exists"))
		return
	}
	if len(wnftId) == 0 {
		middleware.ResponseError(context, 2001, errors.New("id is required"))
	}
	wnftInfoService := services.WNFTInfoService{}
	detail, err := wnftInfoService.Detail(context, wnftId)
	if err != nil {
		middleware.ResponseError(context, 2002, err)
		return
	}
	middleware.ResponseSuccess(context, detail)
	return
}

func (wnftInfo *WnftInfoController) Mint(context *gin.Context) {
	input := &dto.MintWnftInput{}
	if err := input.BindingValidParams(context); err != nil {
		middleware.ResponseError(context, 2001, err)
		return
	}

	wnftInfoService := services.WNFTInfoService{}
	err := wnftInfoService.Mint(context, input)
	if err != nil {
		middleware.ResponseError(context, 2002, err)
		return
	}
	middleware.ResponseSuccess(context, "mint success")
	return
}

func (wnftInfo *WnftInfoController) List(context *gin.Context) {
	input := &dto.WnftListInput{}
	if err := input.BindingValidParams(context); err != nil {
		middleware.ResponseError(context, 2001, err)
		return
	}

	wnftInfoService := services.WNFTInfoService{}
	err := wnftInfoService.List(context, input)
	if err != nil {
		middleware.ResponseError(context, 2002, err)
		return
	}
	middleware.ResponseSuccess(context, "list success")
	return
}

func (wnftInfo *WnftInfoController) Buy(context *gin.Context) {
	input := &dto.BuyWnftInput{}
	if err := input.BindingValidParams(context); err != nil {
		middleware.ResponseError(context, 2001, err)
		return
	}

	wnftInfoService := services.WNFTInfoService{}
	err := wnftInfoService.Buy(context, input)
	if err != nil {
		middleware.ResponseError(context, 2002, err)
		return
	}
	middleware.ResponseSuccess(context, "buy success")
	return
}

func (wnftInfo *WnftInfoController) PayHoldfee(context *gin.Context) {
	input := &dto.PayHoldfeeInput{}
	if err := input.BindingValidParams(context); err != nil {
		middleware.ResponseError(context, 2001, err)
		return
	}
	wnftInfoService := services.WNFTInfoService{}
	if err := wnftInfoService.PayHoldfee(context, input); err != nil {
		middleware.ResponseError(context, 2002, err)
		return
	}
	middleware.ResponseSuccess(context, "pay holdfee success")
	return
}

func (wnftInfo *WnftInfoController) Burn(context *gin.Context) {
	input := &dto.BurnWnftInput{}
	if err := input.BindingValidParams(context); err != nil {
		middleware.ResponseError(context, 2001, err)
		return
	}

	wnftInfoService := services.WNFTInfoService{}
	err := wnftInfoService.Burn(context, input)
	if err != nil {
		middleware.ResponseError(context, 2002, err)
		return
	}
	middleware.ResponseSuccess(context, "burn success")
	return
}

func (wnftInfo *WnftInfoController) Update(context *gin.Context) {

}

func (wnftInfo *WnftInfoController) UpdateDetail(context *gin.Context) {
	input := &dto.WnftDetailInput{}
	if err := input.BindingValidParams(context); err != nil {
		middleware.ResponseError(context, 2001, err)
		return
	}
	wnftInfoService := services.WNFTInfoService{}
	err := wnftInfoService.UpdateDetail(context, input)
	if err != nil {
		middleware.ResponseError(context, 2002, err)
		return
	}
	middleware.ResponseSuccess(context, "mint success")
	return
}

func (wnftInfo *WnftInfoController) UploadResource(context *gin.Context) {
	file, err := context.FormFile("data")
	if err != nil {
		middleware.ResponseError(context, 2002, err)
		return
	}
	src, err := file.Open()
	if err != nil {
		middleware.ResponseError(context, 2002, err)
		return
	}
	defer func(src multipart.File) {
		err := src.Close()
		if err != nil {
			middleware.ResponseError(context, 2002, err)
			log.Println(err)
			return
		}
	}(src)

	cidStr, err := utils.StoreFileIntoIpfs(context, src)
	if err != nil {
		middleware.ResponseError(context, 2002, err)
		return
	}
	log.Println("cid: ", cidStr)

	middleware.ResponseSuccess(context, cidStr)
	return
}

func (wnftInfo *WnftInfoController) FetchResource(context *gin.Context) {
	input := &dto.FetchResourceInput{}
	if err := input.BindingValidParams(context); err != nil {
		middleware.ResponseError(context, 2001, err)
	}

	chainContext := &chain.Context{}
	chainContext.SetChainOpt(chain.Mapping[input.ChainId])
	// tokenURI is the cid
	tokenURI, err := chainContext.GetTokenURI(input.Fp)

	if err != nil {
		middleware.ResponseError(context, 2002, err)
		return
	}

	if tokenURI == "" {
		middleware.ResponseSuccess(context, "")
		return
	}

	file, err := utils.FetchFileFromIpfs(context, tokenURI)
	if err != nil {
		middleware.ResponseError(context, 2002, err)
		return
	}
	defer func(file files.File) {
		err := file.Close()
		if err != nil {
			middleware.ResponseError(context, 2002, err)
			return
		}
	}(file)

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil && err != io.EOF {
		middleware.ResponseError(context, 2002, err)
		return
	}

	contentType := http.DetectContentType(buffer)

	// reset the position of the point
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		middleware.ResponseError(context, 2002, err)
		return
	}

	log.Println("contentType: ", contentType)
	context.Header("Content-Type", contentType)
	_, err = io.Copy(context.Writer, file)
	if err != nil {
		middleware.ResponseError(context, 2002, err)
		return
	}
}