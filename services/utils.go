package services

import (
	"github.com/gin-gonic/gin"
	"github.com/letterScape/backend/dao"
	"gorm.io/gorm"
	"math/big"
	"strings"
)

// todo Cache the config data

// EnlargeSymbol 10^18WEI -> 1ETH
func EnlargeSymbol(c *gin.Context, db *gorm.DB, chainId string, amt string) (string, error) {
	config := &dao.ConfigParams{}
	multiplier, err := config.GetSymbolMultiplier(c, db, chainId)
	if err != nil {
		return amt, err
	}
	amtBig := new(big.Rat)
	amtBig.SetString(amt)
	return trimTrailingZeros(new(big.Rat).Quo(amtBig, multiplier).FloatString(10)), nil
}

// ShrinkSymbol 1ETH -> 10^18WEI
func ShrinkSymbol(c *gin.Context, db *gorm.DB, chainId string, amt string) (string, error) {
	config := &dao.ConfigParams{}
	multiplier, err := config.GetSymbolMultiplier(c, db, chainId)
	if err != nil {
		return amt, err
	}
	amtBig := new(big.Rat)
	amtBig.SetString(amt)
	return trimTrailingZeros(amtBig.Mul(amtBig, multiplier).FloatString(10)), nil
}

func trimTrailingZeros(input string) string {
	if strings.Contains(input, ".") {
		input = strings.TrimRight(input, "0")
		input = strings.TrimRight(input, ".") // if no number after dot, trim the dot
	}
	return input
}
