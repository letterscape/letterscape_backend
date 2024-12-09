package services

import (
	"errors"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/letterScape/backend/dao"
	"github.com/letterScape/backend/dto"
	"github.com/letterScape/backend/vo"
	"time"
)

type SpaceService struct {
}

func (space *SpaceService) Info(c *gin.Context, contentId string) (*vo.SpaceContentVO, error) {
	db, err := lib.GetGormPool("default")
	if err != nil {
		return nil, err
	}

	content, err := (&dao.SpaceContent{}).FindById(c, db, contentId)
	if err != nil {
		return nil, err
	}
	contentVO := &vo.SpaceContentVO{}
	err = copier.Copy(&contentVO, &content)
	if err != nil {
		return nil, err
	}
	return contentVO, nil
}

func (space *SpaceService) Create(c *gin.Context, input *dto.SpaceContentSaveInput) error {
	db, err := lib.GetGormPool("default")
	if err != nil {
		return err
	}

	spaceContent := &dao.SpaceContent{
		ContentId:  input.ContentId,
		ChainId:    input.ChainId,
		Author:     input.Author,
		Title:      input.Title,
		Resource:   input.Resource,
		IsShown:    false,
		IsDeleted:  false,
		CreateTime: time.Now(),
		ModifyTime: time.Now(),
		Readonly:   false,
	}

	if err := spaceContent.Save(c, db); err != nil {
		return err
	}
	return nil
}

func (space *SpaceService) Publish(c *gin.Context, input *dto.ContentPublishInput) error {
	db, err := lib.GetGormPool("default")
	if err != nil {
		return err
	}
	content, err := (&dao.SpaceContent{}).FindById(c, db, input.ContentId)
	if err != nil {
		return err
	}
	if content == nil {
		return errors.New("content not found")
	}

	//todo check the author

	content.IsShown = true
	content.ModifyTime = time.Now()

	if err := content.UpdateById(c, db); err != nil {
		return err
	}
	return nil
}

func (space *SpaceService) Page(c *gin.Context, input *dto.SpacePageInput) (*vo.SpaceContentVOList, error) {
	db, err := lib.GetGormPool("default")
	if err != nil {
		return nil, err
	}

	contentList, total, err := (&dao.SpaceContent{}).PageList(c, db, input)
	if err != nil {
		return nil, err
	}

	var voList []vo.SpaceContentVO
	size := len(*contentList)
	for i := 0; i < size; i++ {
		content := (*contentList)[i]
		contentVO := &vo.SpaceContentVO{}
		err := copier.Copy(&contentVO, &content)
		if err != nil {
			return nil, err
		}
		voList = append(voList, *contentVO)
	}
	pageList := &vo.SpaceContentVOList{List: &voList, Total: total}
	return pageList, nil
}
