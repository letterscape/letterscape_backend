package dao

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"math/big"
	"time"
)

type ConfigParams struct {
	ConfigId   string    `json:"configId" gorm:"column:config_id;primaryKey"`
	Type       string    `json:"type" gorm:"column:type"`
	Param      string    `json:"param" gorm:"column:param"`
	ParamName  string    `json:"paramName" gorm:"column:param_name"`
	Value      string    `json:"value" gorm:"column:value"`
	ValueName  string    `json:"valueName" gorm:"column:value_name"`
	Remark     string    `json:"remark" gorm:"column:remark"`
	CreateTime time.Time `json:"createTime" gorm:"column:create_time"`
	ModifyTime time.Time `json:"modifyTime" gorm:"column:modify_time"`
}

func (config *ConfigParams) TableName() string {
	return "t_config_params"
}

func (config *ConfigParams) SelectOne(c *gin.Context, db *gorm.DB) (*ConfigParams, error) {

	dbSession := db.WithContext(c)

	configParams := &ConfigParams{}
	err := dbSession.Where("type = ? and param = ?", config.Type, config.Param).First(configParams).Error
	if err != nil {
		return nil, err
	}
	return configParams, nil
}

func (config *ConfigParams) Save(c *gin.Context, db *gorm.DB) error {
	dbSession := db.WithContext(c)
	err := dbSession.Save(config).Error
	if err != nil {
		return err
	}
	return nil
}

func (config *ConfigParams) GetSymbolMultiplier(c *gin.Context, db *gorm.DB, chainId string) (*big.Int, error) {
	config.Type = "symbol"
	config.Param = chainId
	param, err := config.SelectOne(c, db)
	if err != nil {
		return nil, err
	}
	multiplier := new(big.Int)
	multiplier.SetString(param.Value, 10)
	return multiplier, nil
}
