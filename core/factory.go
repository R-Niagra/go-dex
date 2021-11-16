package core

import (
	"errors"
	"fmt"

	"github.com/R-Niagra/go-dex/erc20"
)

type UniswapV2Factory struct {
	GetPair  map[erc20.IErcToken]map[erc20.IErcToken]*UniswapV2Pair
	AllPairs []*UniswapV2Pair
	Address  string
}

func NewFactory(_address string) *UniswapV2Factory {
	newFactory := &UniswapV2Factory{
		Address: _address,
	}
	newFactory.GetPair = make(map[erc20.IErcToken]map[erc20.IErcToken]*UniswapV2Pair)
	return newFactory
}

//AllPairsLength returns the number of pairs created
func (fact *UniswapV2Factory) AllPairsLength() int {
	return len(fact.AllPairs)
}

//PairsExists checks if token pair exists in GetPair
func (fact *UniswapV2Factory) PairExists(tokenA erc20.IErcToken, tokenB erc20.IErcToken) bool {
	if _, ok := fact.GetPair[tokenA][tokenB]; ok {
		return true
	}

	return false
}

func (fact *UniswapV2Factory) GetPoolPair(tokenA erc20.IErcToken, tokenB erc20.IErcToken) *UniswapV2Pair {
	if pair, ok := fact.GetPair[tokenA][tokenB]; ok {
		return pair
	}

	return nil
}

//CreatePair creates a liquidity pool using two tokens if the pool doesn't exists before
func (fact *UniswapV2Factory) CreatePair(tokenA erc20.IErcToken, tokenB erc20.IErcToken) (*UniswapV2Pair, error) {
	tokenAadd := tokenA.GetAddress()
	tokenBadd := tokenB.GetAddress()

	if tokenAadd == tokenBadd {
		return nil, errors.New("pair can't have same token addresses")
	}

	var token0, token1 erc20.IErcToken

	if tokenAadd < tokenBadd {
		token0 = tokenA
		token1 = tokenB
	} else {
		token0 = tokenB
		token1 = tokenA
	}
	if token0.GetAddress() == "" {
		return nil, errors.New("address can't be zero")
	}
	//Check if pair already exists
	if _, ok := fact.GetPair[token0][token1]; ok {
		return nil, errors.New("pair exists")
	}

	//creating pair
	pair, err := CreatePair(tokenA, tokenB, fact.Address)
	if err != nil {
		return nil, err
	}

	fact.GetPair[token0][token1] = pair
	fact.GetPair[token1][token0] = pair

	fact.AllPairs = append(fact.AllPairs, pair)

	fmt.Println("Pair created!!")

	return pair, nil
}
