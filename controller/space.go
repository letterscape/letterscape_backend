package controller

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/ipfs/boxo/files"
	"github.com/letterScape/backend/dto"
	"github.com/letterScape/backend/middleware"
	"github.com/letterScape/backend/services"
	"github.com/letterScape/backend/utils"
	"io"
	"log"
	"os"
	"time"
)

type SpaceController struct{}

func SpaceRegister(router *gin.RouterGroup) {
	space := SpaceController{}
	router.GET("/info", space.Info)
	router.GET("/page", space.Page)
	router.POST("/create", space.Create).OPTIONS("/create", space.Create)
	router.POST("/publish", space.Publish).OPTIONS("/publish", space.Publish)
	router.POST("/upload", space.Upload).OPTIONS("/upload", space.Upload)
	router.GET("/fetch", space.Fetch)
}

func (space *SpaceController) Info(context *gin.Context) {
	contentId, exists := context.GetQuery("id")
	if !exists {
		middleware.ResponseError(context, 2003, errors.New("resource not found"))
	}

	service := &services.SpaceService{}
	contentVO, err := service.Info(context, contentId)
	if err != nil {
		middleware.ResponseError(context, 2001, err)
		return
	}
	middleware.ResponseSuccess(context, contentVO)
}

func (space *SpaceController) Page(context *gin.Context) {
	pageInput := &dto.SpacePageInput{}
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

	service := &services.SpaceService{}
	wnftList, err := service.Page(context, pageInput)
	if err != nil {
		middleware.ResponseError(context, 2002, err)
		return
	}

	middleware.ResponseSuccess(context, wnftList)
	return
}

func (space *SpaceController) Create(context *gin.Context) {
	input := &dto.SpaceContentSaveInput{}
	if err := input.BindingValidParams(context); err != nil {
		middleware.ResponseError(context, 2001, err)
		return
	}

	service := &services.SpaceService{}
	err := service.Create(context, input)
	if err != nil {
		middleware.ResponseError(context, 2002, err)
		return
	}

	middleware.ResponseSuccess(context, "success")
	return
}

func (space *SpaceController) Publish(context *gin.Context) {
	inpput := &dto.ContentPublishInput{}
	if err := inpput.BindingValidParams(context); err != nil {
		middleware.ResponseError(context, 2001, err)
		return
	}
	service := &services.SpaceService{}
	if err := service.Publish(context, inpput); err != nil {
		middleware.ResponseError(context, 2002, err)
		return
	}
	middleware.ResponseSuccess(context, "success")
	return
}

func (space *SpaceController) Upload(context *gin.Context) {
	input := &dto.UploadContentInput{}
	if err := input.BindingValidParams(context); err != nil {
		middleware.ResponseError(context, 2001, err)
		return
	}
	tempFile, err := os.CreateTemp("", input.ContentId+"_"+time.Now().Format("20060102150405")+".txt")
	if err != nil {
		middleware.ResponseError(context, 2002, err)
		return
	}

	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			middleware.ResponseError(context, 2003, err)
			return
		}
	}(tempFile.Name())

	// write content into file
	_, err = io.Copy(tempFile, bytes.NewReader([]byte(input.Content)))
	if err != nil {
		middleware.ResponseError(context, 2002, err)
		return
	}

	// reset file pointer to the start
	_, err = tempFile.Seek(0, 0)
	if err != nil {
		middleware.ResponseError(context, 2002, err)
		log.Fatalf("Failed to reset file pointer: %v", err)
	}

	cid, err := utils.StoreFileIntoIpfs(context, tempFile)
	if err != nil {
		middleware.ResponseError(context, 2002, err)
		return
	}

	middleware.ResponseSuccess(context, cid)
	return
}

func (space *SpaceController) Fetch(context *gin.Context) {
	resource, exists := context.GetQuery("resource")
	if !exists {
		middleware.ResponseError(context, 2003, errors.New("resource not found"))
	}

	file, err := utils.FetchFileFromIpfs(context, resource)
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

	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	if err != nil {
		middleware.ResponseError(context, 2002, err)
		return
	}

	middleware.ResponseSuccess(context, buf.String())
	return
}
