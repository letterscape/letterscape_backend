package services

import (
	"github.com/gin-gonic/gin"
	"github.com/letterScape/backend/dao"
	"gorm.io/gorm"
	"math/big"
)

// todo Cache the config data

// EnlargeSymbol 10^18WEI -> 1ETH
func EnlargeSymbol(c *gin.Context, db *gorm.DB, chainId string, amt string) (string, error) {
	config := &dao.ConfigParams{}
	multiplier, err := config.GetSymbolMultiplier(c, db, chainId)
	if err != nil {
		return amt, err
	}
	amtBig := big.NewInt(0)
	amtBig.SetString(amt, 10)
	return amtBig.Div(amtBig, multiplier).String(), nil
}

// ShrinkSymbol 1ETH -> 10^18WEI
func ShrinkSymbol(c *gin.Context, db *gorm.DB, chainId string, amt string) (string, error) {
	config := &dao.ConfigParams{}
	multiplier, err := config.GetSymbolMultiplier(c, db, chainId)
	if err != nil {
		return amt, err
	}
	amtBig := big.NewInt(0)
	amtBig.SetString(amt, 10)
	return amtBig.Mul(amtBig, multiplier).String(), nil
}
