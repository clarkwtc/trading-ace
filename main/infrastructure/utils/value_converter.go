package utils

import (
    "github.com/ethereum/go-ethereum/params"
    "math"
    "math/big"
)

func ForamtEther(value *big.Int) string {
    return new(big.Float).Quo(new(big.Float).SetInt(value), big.NewFloat(params.Ether)).Text('f', 18)
}

func ForamtGWei(value *big.Int) string {
    return new(big.Float).Quo(new(big.Float).SetInt(value), big.NewFloat(params.GWei)).Text('f', 9)
}

func ToUSDC(value *big.Int) *big.Int {
    scale := new(big.Float).SetFloat64(math.Pow10(6))
    valueFloat := new(big.Float).SetInt(value)

    result := new(big.Int)
    valueFloat.Mul(valueFloat, scale)
    valueFloat.Int(result)
    return result
}

func ForamtUSDC(value *big.Int) string {
    return new(big.Float).Quo(new(big.Float).SetInt(value), new(big.Float).SetFloat64(math.Pow10(6))).Text('f', 6)
}
