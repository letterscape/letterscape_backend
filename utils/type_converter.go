package utils

import (
	"fmt"
	"math/big"
)

func ToBigInt(value interface{}) (*big.Int, error) {
	switch v := value.(type) {
	case big.Int:
		return &v, nil
	case *big.Int:
		if v == nil {
			return nil, fmt.Errorf("nil *big.Int pointer")
		}
		return v, nil
	default:
		// 类型不匹配，返回错误
		return nil, fmt.Errorf("unsupported type: %T", value)
	}
}
