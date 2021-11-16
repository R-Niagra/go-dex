package lib

import (
	"github.com/R-Niagra/go-dex/core"
	"github.com/R-Niagra/go-dex/erc20"
)

func SortTokens(tokenA erc20.IErcToken, tokenB erc20.IErcToken) (erc20.IErcToken, erc20.IErcToken) {
	if tokenA.GetAddress() < tokenB.GetAddress() {
		return tokenA, tokenB
	}
	return tokenB, tokenA
}

//GetReserves get reserves of the tokens
func GetReserves(factory core.IUniswapV2Factory, tokenA erc20.IErcToken, tokenB erc20.IErcToken) (uint64, uint64) {
	pair := factory.GetPoolPair(tokenA, tokenB)
	if pair == nil {
		panic("pair doesn't exist")
	}

	res0, res1, _ := pair.GetReserves()

	if tokenA.GetAddress() < tokenB.GetAddress() {
		return res0, res1
	}

	return res1, res0
}

func GetPairAddress(factory core.IUniswapV2Factory, tokenA erc20.IErcToken, tokenB erc20.IErcToken) string {
	pair := factory.GetPoolPair(tokenA, tokenB)
	if pair == nil {
		panic("pair doesn't exist")
	}

	return pair.GetPairAddress()
}

//Quote return equivalent amount of the other asset
func Quote(amountA uint64, reserveA uint64, reserveB uint64) uint64 {
	if amountA == 0 {
		panic("amount cannot be zero")
	}

	if reserveA <= 0 || reserveB <= 0 {
		panic("Insufficient liquidity")
	}

	amountB := (amountA * reserveB) / reserveA
	return amountB
}

// func GetReserves(factory core.IUniswapV2Factory, tokenA erc20.IErcToken, tokenB erc20.IErcToken) (uint64, uint64) {

// }

//GetAmountsOut calculates AmountOut for every pair in the path
func GetAmountsOut(factory core.IUniswapV2Factory, amountIn uint64, path []erc20.IErcToken) []uint64 {

	pathLen := len(path)
	if pathLen < 2 {
		panic("path length is less than 2")
	}
	amounts := make([]uint64, pathLen)
	amounts[0] = amountIn

	for i := 0; i < pathLen-1; i++ {
		reserveIn, reserveOut := GetReserves(factory, path[i], path[i+1])
		amounts[i+1] = GetAmountOut(amountIn, reserveIn, reserveOut)
	}
	return amounts
}

//GetAmountOut calculates the maximum output against the input amount of asset
func GetAmountOut(amountIn uint64, reserveIn uint64, reserveOut uint64) uint64 {

	if amountIn <= 0 {
		panic("input amount cannot be 0")
	}
	if reserveIn <= 0 || reserveOut <= 0 {
		panic("uniswapV2Library: INSUFFICIENT_LIQUIDITY")
	}

	amountInWithFee := amountIn * 997
	numerator := amountInWithFee * (reserveOut)
	denominator := (reserveIn * 1000) + (amountInWithFee)
	return numerator / denominator

}
