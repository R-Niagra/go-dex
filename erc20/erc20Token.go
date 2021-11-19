package erc20

import (
	"fmt"
	"sync"
)

type ERC20Token struct {
	Name    string
	Symbol  string
	Address string

	balances    map[string]uint64
	totalSupply uint64
	allowance   map[string]map[string]uint64
	mu          sync.Mutex
}

//NewIERCToken returns interface against ERC20Token
func NewIERCToken(_name string, _symbol string, _totalSupply uint64, sender_add string, _address string) IErcToken {
	return NewERC20Token(_name, _symbol, _totalSupply, sender_add, _address)
}

//NewERC20Token creates and return new instance of ERC20 token
func NewERC20Token(_name string, _symbol string, _totalSupply uint64, sender_add string, _address string) *ERC20Token {
	token := &ERC20Token{
		Name:        _name,
		Symbol:      _symbol,
		totalSupply: _totalSupply,
		Address:     _address,
	}

	token.balances = make(map[string]uint64)
	token.allowance = make(map[string]map[string]uint64)

	token.balances[sender_add] = _totalSupply

	return token
}

func (t *ERC20Token) GetIERCToken() IErcToken {
	return t
}

//TokenName returns the token name
func (t *ERC20Token) TokenName() string {
	return t.Name
}

//TokenSymbol returns the token symbol
func (t *ERC20Token) TokenSymbol() string {
	return t.Symbol
}

//TotalSupply return the total supply of token
func (t *ERC20Token) TotalSupply() uint64 {
	return t.totalSupply
}

//TotalSupply return the total supply of token
func (t *ERC20Token) GetAddress() string {
	return t.Address
}

//BalanceOf gives the token balance of user address
func (t *ERC20Token) BalanceOf(address string) uint64 {

	if balance, ok := t.balances[address]; ok {
		return balance
	}

	return 0
}

// //BalanceOfP gives the token balance of user address
// func (t *ERC20Token) BalanceOfP(address string) uint64 {
// 	t.mu.Lock()
// 	defer t.mu.Unlock()
// 	if balance, ok := t.balances[address]; ok {
// 		return balance
// 	}

// 	return 0
// }

//addBalance will add amount to the address
func (t *ERC20Token) addBalance(address string, amount uint64) {
	if _, ok := t.balances[address]; ok {
		t.balances[address] += amount
		return
	}

	// fmt.Println("Address doesn't exist")
	t.balances[address] = amount
}

//addBalance will add amount to the address
func (t *ERC20Token) subtractBalance(address string, amount uint64) {
	if _, ok := t.balances[address]; ok {
		t.balances[address] -= amount
		return
	}

}

//Transfer transfers token from sender address to the receiver address
func (t *ERC20Token) Transfer(sender string, receiver string, amount uint64) bool {
	senderBalance := t.BalanceOf(sender)
	if senderBalance < amount {
		fmt.Println("Insufficient balance")
		return false
	}

	//subtract amount from the sender
	t.balances[sender] -= amount

	//add amount to the receiver
	t.addBalance(receiver, amount)

	return true
}

//TransferP transfers token from sender address to the receiver address. It is safe to execute in parallel
func (t *ERC20Token) TransferP(sender string, receiver string, amount uint64) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	senderBalance := t.BalanceOf(sender)
	if senderBalance < amount {
		fmt.Println("Insufficient balance")
		return false
	}

	//subtract amount from the sender
	t.balances[sender] -= amount

	//add amount to the receiver
	t.addBalance(receiver, amount)

	return true
}

//TransferFrom transfers from owner to the buyer address given that owner has already set allowance
func (t *ERC20Token) TransferFrom(owner string, buyer string, amount uint64) bool {

	if amount > t.BalanceOf(owner) {
		fmt.Println("Owner does't have sufficient balance")
		panic("Insufficient funds")
		return false
	}
	if amount > t.allowed(owner, buyer) {
		res := fmt.Sprintf(owner, " buyer: ", buyer, "amount: ", amount, " actual", t.allowed(owner, buyer))
		panic(res)
		return false
	}

	t.subtractBalance(owner, amount)
	t.allowance[owner][buyer] -= amount
	// fmt.Println("Adding balance to user", buyer)
	t.addBalance(buyer, amount)
	// fmt.Println("Add done")
	return true
}

//TransferFromP transfers from owner to the buyer address given that owner has already set allowance. Safe to use in parallel
func (t *ERC20Token) TransferFromP(owner string, buyer string, amount uint64) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	if amount > t.BalanceOf(owner) {
		fmt.Println("Owner does't have sufficient balance")
		return false
	}
	if amount > t.allowed(owner, buyer) {
		res := fmt.Sprintf(owner, " buyer: ", buyer, "amount: ", amount, " actual", t.allowed(owner, buyer))
		panic(res)
		return false
	}

	t.subtractBalance(owner, amount)
	t.allowance[owner][buyer] -= amount
	// fmt.Println("Adding balance to user", buyer)
	t.addBalance(buyer, amount)
	// fmt.Println("Add done")
	return true
}

func (t *ERC20Token) allowed(owner string, buyer string) uint64 {
	if amount, ok := t.allowance[owner][buyer]; ok {
		return amount
	}
	return 0
}

// func (t *ERC20Token) allowedP(owner string, buyer string) uint64 {
// 	t.mu.Lock()
// 	defer t.mu.Unlock()
// 	if amount, ok := t.allowance[owner][buyer]; ok {
// 		return amount
// 	}
// 	return 0
// }

func (t *ERC20Token) Approve(owner string, spender string, value uint64) {
	if _, ok := t.allowance[owner]; !ok {
		t.allowance[owner] = make(map[string]uint64)
		// fmt.Println("Creating allowed map")
	}

	t.allowance[owner][spender] = value
}

func (t *ERC20Token) ApproveP(owner string, spender string, value uint64) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, ok := t.allowance[owner]; !ok {
		t.allowance[owner] = make(map[string]uint64)
		fmt.Println("Creating allowed map")
	}

	t.allowance[owner][spender] += value

}

func (t *ERC20Token) _mint(to string, value uint64) {
	t.totalSupply += value
	t.addBalance(to, value)
}

func (t *ERC20Token) _burn(from string, value uint64) {
	t.subtractBalance(from, value)
	t.totalSupply -= value
}
