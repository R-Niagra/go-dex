package core

import (
	"github.com/R-Niagra/go-dex/erc20"
)

// //
// type IErcToken interface {
// 	TotalSupply() uint64
// 	BalanceOf(sender string) uint64
// 	// Allowance(address owner, address spender;
// 	TokenName() string
// 	TokenSymbol() string
// 	GetAddress() string
// 	Transfer(sender string, receiver string, amount uint64) bool
// 	// Approve(address spender, uint64) bool
// 	// TransferFrom(sender string, receiver string, amount uint64) bool
// 	// event Transfer(address indexed from, address indexed to, uint256 value);
// 	// event Approval(address indexed owner, address indexed spender, uint256 value);
// }

type IUniswapV2Factory interface {

	// function feeTo() external view returns (address);
	// function feeToSetter() external view returns (address);

	GetPoolPair(tokenA erc20.IErcToken, tokenB erc20.IErcToken) *UniswapV2Pair
	// function allPairs(uint) external view returns (address pair);
	PairExists(tokenA erc20.IErcToken, tokenB erc20.IErcToken) bool

	AllPairsLength() int

	CreatePair(tokenA erc20.IErcToken, tokenB erc20.IErcToken, pairAdd, ercAdd, senderAdd string) (*UniswapV2Pair, error)
	// function createPair(address tokenA, address tokenB) external returns (address pair);

	// function setFeeTo(address) external;
	// function setFeeToSetter(address) external;
}
