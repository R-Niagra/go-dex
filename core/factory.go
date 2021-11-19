package core

import (
	"errors"
	"fmt"
	"sync"

	"github.com/R-Niagra/go-dex/erc20"
)

type UniswapV2Factory struct {
	GetPair  map[erc20.IErcToken]map[erc20.IErcToken]*UniswapV2Pair
	AllPairs []*UniswapV2Pair
	Address  string
	mu       sync.Mutex
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

func (fact *UniswapV2Factory) GetPoolPairP(tokenA erc20.IErcToken, tokenB erc20.IErcToken) *UniswapV2Pair {
	fact.mu.Lock()
	defer fact.mu.Unlock()
	if pair, ok := fact.GetPair[tokenA][tokenB]; ok {
		return pair
	}

	return nil
}

//CreatePair creates a liquidity pool using two tokens if the pool doesn't exists before
func (fact *UniswapV2Factory) CreatePair(tokenA erc20.IErcToken, tokenB erc20.IErcToken, pairAdd, ercAdd, senderAdd string) (*UniswapV2Pair, error) {
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
	pair, err := CreatePair(tokenA, tokenB, fact.Address, pairAdd, ercAdd, senderAdd)
	if err != nil {
		return nil, err
	}

	if _, ok := fact.GetPair[token0]; !ok {
		fact.GetPair[token0] = make(map[erc20.IErcToken]*UniswapV2Pair)
	}

	fact.GetPair[token0][token1] = pair

	if _, ok := fact.GetPair[token1]; !ok {
		fact.GetPair[token1] = make(map[erc20.IErcToken]*UniswapV2Pair)
	}

	fact.GetPair[token1][token0] = pair

	fact.AllPairs = append(fact.AllPairs, pair)

	fmt.Println("Pair created!!")

	return pair, nil
}
