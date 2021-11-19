package router

import (
	"fmt"
	"sync"

	"github.com/R-Niagra/go-dex/core"
	"github.com/R-Niagra/go-dex/erc20"
	"github.com/R-Niagra/go-dex/router/lib"
)

type UniswapV2Router02 struct {
	factory core.IUniswapV2Factory
	mu      sync.Mutex
}

//NewRouter returns new instance of the router
func NewRouter(_factory core.IUniswapV2Factory) *UniswapV2Router02 {
	router := &UniswapV2Router02{
		factory: _factory,
	}
	return router
}

func (r *UniswapV2Router02) _addLiquidity(tokenA erc20.IErcToken, tokenB erc20.IErcToken, amountADesired uint64, amountBDesired uint64, amountAMin uint64, amountBMin uint64, to string) (uint64, uint64) {
	if !r.factory.PairExists(tokenA, tokenB) { //create pair if it doesn't exists
		fmt.Println("Pair doesn't exist. creating pair")
		pairAdd := tokenA.TokenSymbol() + tokenB.TokenSymbol() + "Pair"
		ercAdd := tokenA.TokenSymbol() + tokenB.TokenSymbol() + "ERC"
		r.factory.CreatePair(tokenA, tokenB, pairAdd, ercAdd, to)
	}
	reserveA, reserveB := lib.GetReserves(r.factory, tokenA, tokenB)
	if reserveA == 0 && reserveB == 0 { //if it is empty then the creater decides the rate
		return amountADesired, amountBDesired
	}

	//Using rate find tokenB given the desired tokenA amount
	amountBOptimal := lib.Quote(amountADesired, reserveA, reserveB)

	if amountBOptimal <= amountBDesired {
		if amountBOptimal < amountBMin {
			panic("router: insufficient B amount")
		}
		return amountADesired, amountBOptimal
	}

	amountAOptimal := lib.Quote(amountBDesired, reserveB, reserveA)
	if amountAOptimal > amountADesired || amountAOptimal < amountAMin {
		panic("out of acceptible window")
	}

	return amountAOptimal, amountBDesired
}

//AddLiquidity add liquidity in the pool if the rate falls under desired window
func (r *UniswapV2Router02) AddLiquidity(tokenA erc20.IErcToken, tokenB erc20.IErcToken, amountADesired uint64, amountBDesired uint64, amountAMin uint64, amountBMin uint64, senderAdd string) {
	//shall ensure deadline
	///
	amountA, amountB := r._addLiquidity(tokenA, tokenB, amountADesired, amountBDesired, amountAMin, amountBMin, senderAdd)
	pairAdd := lib.GetPairAddress(r.factory, tokenA, tokenB)

	//send tokenA to pair contract address
	tokenA.Transfer(senderAdd, pairAdd, amountA)
	fmt.Println("AddLiquidity: tokenA sent to pair Add")
	//send tokenB to pair contract address
	tokenB.Transfer(senderAdd, pairAdd, amountB)
	fmt.Println("AddLiquidity: tokenB sent to pair Add")

	//TODO: add liquidity
	pair := r.factory.GetPoolPair(tokenA, tokenB)
	//mint tokens and update reserves
	pair.Mint(senderAdd)

	fmt.Println("AddLiquidity: successful")
}

func (r *UniswapV2Router02) _swap(amounts []uint64, path []erc20.IErcToken, toAdd string) {
	pathLen := len(path)
	for i := 0; i < pathLen-1; i++ {
		inpTok, outTok := path[i], path[i+1]
		token0, _ := lib.SortTokens(inpTok, outTok)
		amountOut := amounts[i+1]

		var amount0Out, amount1Out uint64
		if inpTok.GetAddress() == token0.GetAddress() {
			amount0Out = 0
			amount1Out = amountOut
		} else {
			amount1Out = 0
			amount0Out = amountOut
		}

		pair := r.factory.GetPoolPair(inpTok, outTok)

		pair.Swap(amount0Out, amount1Out, toAdd)

	}
}

//SwapExactTokensForTokens swaps the exact number of input tokens with the maximum output token if rate is within the range
func (r *UniswapV2Router02) SwapExactTokensForTokens(amountIn uint64, amountOutMin uint64, path []erc20.IErcToken, toAdd string, senderAdd string) {
	amounts := lib.GetAmountsOut(r.factory, amountIn, path)
	// fmt.Println("Output amount is: ", amounts[len(amounts)-1])
	if amounts[len(amounts)-1] < amountOutMin {
		panic("Insufficient output amount")
	}

	//transfer amountIn to the pair contract
	pairAdd := lib.GetPairAddress(r.factory, path[0], path[1])
	// fmt.Println("in swap, ", path[0].BalanceOf(pairAdd), pairAdd)

	success := path[0].TransferFrom(senderAdd, pairAdd, amounts[0])
	if !success {
		panic("Transfer failed")
	}
	// fmt.Println("in swap2, ", path[0].BalanceOf(pairAdd), pairAdd)

	r._swap(amounts, path, toAdd)
}

//SwapExactTokensForTokens swaps the exact number of input tokens with the maximum output token if rate is within the range
func (r *UniswapV2Router02) SwapExactTokensForTokensP(amountIn uint64, amountOutMin uint64, path []erc20.IErcToken, toAdd string, senderAdd string, swapMu *sync.Mutex) {
	// r.mu.Lock()

	swapMu.Lock()
	amounts := lib.GetAmountsOut(r.factory, amountIn, path)
	fmt.Println("Output amount is: ", amounts[len(amounts)-1])
	if amounts[len(amounts)-1] < amountOutMin {
		panic("Insufficient output amount")
	}

	//transfer amountIn to the pair contract
	pairAdd := lib.GetPairAddress(r.factory, path[0], path[1])
	// fmt.Println("in swap, ", path[0].BalanceOf(pairAdd), pairAdd)

	success := path[0].TransferFromP(senderAdd, pairAdd, amounts[0])
	if !success {
		panic("Transfer failed")
	}
	// fmt.Println("in swap2, ", path[0].BalanceOf(pairAdd), pairAdd)

	r._swap(amounts, path, toAdd)
	swapMu.Unlock()
	// r.mu.Unlock()
}
