package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ipfs/boxo/files"
	"github.com/letterScape/backend/chain"
	"github.com/letterScape/backend/dto"
	"github.com/letterScape/backend/middleware"
	"github.com/letterScape/backend/utils"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

type WnftResource struct{}

func WnftResourceRegister(router *gin.RouterGroup) {
	wnftResource := new(WnftResource)
	router.POST("/save", wnftResource.Save).OPTIONS("/save", wnftResource.Save)
	router.GET("/find", wnftResource.Find)
	router.POST("/upload", wnftResource.UploadResource).OPTIONS("/upload", wnftResource.UploadResource)
	router.GET("/fetch", wnftResource.FetchResource)
}

func (wnftResource *WnftResource) Save(context *gin.Context) {
	input := &dto.ResourceInput{}
	if err := input.BindingValidParams(context); err != nil {
		middleware.ResponseError(context, 2001, err)
		return
	}

	resourceJson, err := json.Marshal(input)
	if err != nil {
		fmt.Println("JSON serialized failed:", err)
		return
	}

	tempFile, err := os.CreateTemp("", input.ResourceId+"_"+time.Now().Format("20060102150405")+".txt")
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
	_, err = io.Copy(tempFile, bytes.NewReader(resourceJson))
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

	log.Println("resource cid: ", cid)
	middleware.ResponseSuccess(context, cid)
	return
}

func (wnftResource *WnftResource) Find(context *gin.Context) {
	input := &dto.FindResourceInput{}
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

	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	if err != nil {
		middleware.ResponseError(context, 2002, err)
		return
	}

	middleware.ResponseSuccess(context, buf.String())
	return
}

func (wnftResource *WnftResource) UploadResource(context *gin.Context) {
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
	log.Println("resourceId: ", cidStr)

	middleware.ResponseSuccess(context, cidStr)
	return
}

func (wnftResource *WnftResource) FetchResource(context *gin.Context) {
	resourceId, exists := context.GetQuery("resourceId")
	if !exists {
		middleware.ResponseError(context, 2003, errors.New("resource not found"))
	}

	if resourceId == "" {
		middleware.ResponseSuccess(context, "")
		return
	}

	file, err := utils.FetchFileFromIpfs(context, resourceId)
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
