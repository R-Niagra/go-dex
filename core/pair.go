package core

import (
	"errors"
	"fmt"

	"github.com/R-Niagra/go-dex/erc20"
)

const (
	MinimumLiquidity = 10000
)

type UniswapV2Pair struct {
	pairToken *erc20.ERC20Token

	reserve0           uint64
	reserve1           uint64
	blockTimestampLast uint32
	unlocked           uint64
	address            string

	FactoryAddress       string
	Token0Add            erc20.IErcToken
	Token1Add            erc20.IErcToken
	Price0CumulativeLast uint64
	Price1CumulativeLast uint64
	KLast                uint64
}

func CreatePair(_token0 erc20.IErcToken, _token1 erc20.IErcToken, _factoryAdd string, _pairAddress string, ercAdd string, senderAdd string) (*UniswapV2Pair, error) {
	// if _factoryAdd[0:6] != "factory" { //TODO only factory should be able to create address
	// 	return nil, errors.New("UniswapV2: FORBIDDEN")
	// }

	poolToken := erc20.NewERC20Token("Uniswap-v2", "Uni", 0, senderAdd, ercAdd)

	pair := &UniswapV2Pair{
		pairToken:      poolToken,
		Token0Add:      _token0,
		address:        _pairAddress,
		Token1Add:      _token1,
		FactoryAddress: _factoryAdd,
	}

	return pair, nil
}

func (pair *UniswapV2Pair) GetReserves() (_reserve0 uint64, _reserve1 uint64, _blockTimestampLast uint32) {
	_reserve0 = pair.reserve0
	_reserve1 = pair.reserve1
	_blockTimestampLast = pair.blockTimestampLast

	return _reserve0, _reserve1, _blockTimestampLast
}

func (pair *UniswapV2Pair) _safeTransfer(token erc20.IErcToken, to string, value uint64) bool {
	//transfering from pair address
	res := token.Transfer(pair.GetPairAddress(), to, value)
	if !res {
		fmt.Println("transfer failed")
		panic("Transfer failed")
	}
	return res

}

func (pair *UniswapV2Pair) GetPairAddress() string {
	return pair.address
}

//Swap swaps the given token with the other if the rate lies under
func (pair *UniswapV2Pair) Swap(amount0Out uint64, amount1Out uint64, to string) error {
	if amount0Out <= 0 && amount1Out <= 0 {
		return errors.New("insufficient amount to swap")
	}

	if amount0Out >= pair.reserve0 || amount1Out >= pair.reserve1 {
		return errors.New("uniswapV2: INSUFFICIENT_LIQUIDITY")
	}

	if amount0Out > 0 {
		pair._safeTransfer(pair.Token0Add, to, amount0Out)
	}
	if amount1Out > 0 {
		pair._safeTransfer(pair.Token1Add, to, amount1Out)
	}

	balance0 := pair.Token0Add.BalanceOf(pair.GetPairAddress())
	balance1 := pair.Token1Add.BalanceOf(pair.GetPairAddress())
	// fmt.Println("bal0 - bal1: ", balance0, balance1)
	//determining the amountIn to the pair
	var amount0In, amount1In uint64
	if balance0 > (pair.reserve0 - amount0Out) {
		amount0In = balance0 - (pair.reserve0 - amount0Out)
	}

	if balance1 > (pair.reserve1 - amount1Out) {
		amount1In = balance1 - (pair.reserve1 - amount1Out)
	}
	fmt.Println("Amount0In Amount1In: ", amount0In, amount1In)
	if amount0In == 0 && amount1In == 0 {
		panic("Insufficient input amount")
	}

	balance0Adjusted := (balance0 * 1000) - (amount0In * 3)
	balance1Adjusted := (balance1 * 1000) - (amount1In * 3)

	if (balance0Adjusted * balance1Adjusted) < (pair.reserve0 * pair.reserve1) {
		fmt.Println(balance0Adjusted*balance1Adjusted, pair.reserve0*pair.reserve1)
		panic("AMM formula doesn't hold")
	}

	//updating the reserves
	pair.reserve0 = balance0
	pair.reserve1 = balance1

	return nil
}

//Mint() mint uninswap tokens and update the reserve amount
func (pair *UniswapV2Pair) Mint(toAdd string) {

	// res0, res1, _ := pair.GetReserves()

	bal0 := pair.Token0Add.BalanceOf(pair.address)
	bal1 := pair.Token1Add.BalanceOf(pair.address)

	// amount0 := bal0 - res0
	// amount1 := bal1 - res1

	//TODO: fee shall be calculated and minted here

	//updating the reserves
	pair.reserve0 = bal0
	pair.reserve1 = bal1

}
