package erc20

type IErcToken interface {
	TotalSupply() uint64
	BalanceOf(sender string) uint64
	// BalanceOfP(sender string) uint64
	// Allowance(address owner, address spender;
	TokenName() string
	TokenSymbol() string
	GetAddress() string
	Transfer(sender string, receiver string, amount uint64) bool
	TransferP(sender string, receiver string, amount uint64) bool
	TransferFrom(owner string, buyer string, amount uint64) bool
	TransferFromP(owner string, buyer string, amount uint64) bool
	Approve(owner string, spender string, value uint64)
	ApproveP(owner string, spender string, value uint64)
	// Approve(address spender, uint64) bool
	// TransferFrom(sender string, receiver string, amount uint64) bool
	// event Transfer(address indexed from, address indexed to, uint256 value);
	// event Approval(address indexed owner, address indexed spender, uint256 value);
}
