package main

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/R-Niagra/go-dex/core"
	"github.com/R-Niagra/go-dex/erc20"
	"github.com/R-Niagra/go-dex/router"
)

//router02 contains all the service funcitons
var router02 *router.UniswapV2Router02

func main() {
	fmt.Println("Hello, Dex.")

	myAdd := "0x5ab9d116a53ef41063e3eae26a7ebe736720e9ba"

	//Creating some ERC20 Tokens
	token1 := erc20.NewERC20Token("Dai", "D", 200000000, myAdd, "token1")
	token2 := erc20.NewERC20Token("Niagra", "N", 50000000, myAdd, "token2")
	token3 := erc20.NewERC20Token("Omni", "O", 40000000, myAdd, "token3")
	token4 := erc20.NewERC20Token("Rio", "R", 80000000, myAdd, "token4")

	//Creating a new factory
	factory1 := core.NewFactory("factory1")

	//creating router to handle service calls
	router02 = router.NewRouter(factory1)

	//pool consisting of token1 and token2 pair
	pair1_2, err := factory1.CreatePair(token1, token2, "DNPair", "DNERC", myAdd)
	if err != nil {
		panic(err)
	}
	printPairReserves(pair1_2)

	pair1_3, err := factory1.CreatePair(token1, token3, "DOPair", "DOERC", myAdd)
	if err != nil {
		panic(err)
	}
	printPairReserves(pair1_3)

	pair3_4, err := factory1.CreatePair(token3, token4, "ORPair", "ORERC", myAdd)
	if err != nil {
		panic(err)
	}
	printPairReserves(pair3_4)

	//adding liquidity to the pair1_2 pool
	fmt.Println("Adding liquidity to the pair1_2 pool...")
	router02.AddLiquidity(token1, token2, 8000000, 4000000, 1999999, 99999, myAdd)
	printPairReserves(pair1_2)
	fmt.Println("Token1 and Token2 balance: ", token1.BalanceOf(myAdd), token2.BalanceOf(myAdd))
	fmt.Println("Token1 and Token2 pair balance: ", token1.BalanceOf(pair1_2.GetPairAddress()), token2.BalanceOf(pair1_2.GetPairAddress()))

	//adding liquidity to the pair1_3 pool
	fmt.Println("Adding liquidity to the pair1_3 pool...")
	router02.AddLiquidity(token1, token3, 6000000, 3000000, 599999, 99999, myAdd)
	printPairReserves(pair1_3)
	fmt.Println("Token1 and Token3 balance: ", token1.BalanceOf(myAdd), token3.BalanceOf(myAdd))
	fmt.Println("Token1 and Token3 pair balance: ", token1.BalanceOf(pair1_3.GetPairAddress()), token3.BalanceOf(pair1_3.GetPairAddress()))

	//adding liquidity to the pair3_4 pool
	fmt.Println("Adding liquidity to the pair2_4 pool...")
	router02.AddLiquidity(token3, token4, 6000000, 8000000, 599999, 99999, myAdd)
	printPairReserves(pair3_4)
	fmt.Println("Token1 and Token3 balance: ", token3.BalanceOf(myAdd), token4.BalanceOf(myAdd))
	fmt.Println("Token1 and Token3 pair balance: ", token3.BalanceOf(pair3_4.GetPairAddress()), token4.BalanceOf(pair3_4.GetPairAddress()))

	//Distributing erc tokens to do swap
	DistributeERC(token1, myAdd, 10000000)
	DistributeERC(token2, myAdd, 3000000)
	DistributeERC(token3, myAdd, 3000000)
	DistributeERC(token4, myAdd, 3000000)

	fmt.Println("Token1 Token2 Token3 balance: ", token1.BalanceOf(myAdd), token2.BalanceOf(myAdd), token3.BalanceOf(myAdd))

	// done := make(chan bool)
	tx := 1000
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	fmt.Println("Testing serial execution")
	completed := testSerialExecution(ctx, tx, token1, token2, pair1_2)
	fmt.Printf("completed: %d transactions in one Seconds(TPS)\n", completed)

	// fmt.Println("Testing parallel execution")
	// completed := testParallelExecution(ctx, tx, token1, token2, pair1_2, "User0")
	// fmt.Printf("completed: %d transactions in one Seconds(TPS)\n", completed)

	// fmt.Println("Testing parallel execution in multiple pools")
	// done := make(chan bool)

	// go func() {
	// 	fmt.Println("Testing parallel execution in first pool")
	// 	completed := testParallelExecution(ctx, tx, token1, token2, pair1_2, "User0")
	// 	fmt.Printf("completed: %d transactions in pool1 in one Seconds(TPS)\n", completed)
	// 	done <- true
	// }()

	// ctx1, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	// defer cancel()
	// go func() {
	// 	fmt.Println("Testing parallel execution in second pool")
	// 	res1 := testParallelExecution(ctx1, tx, token1, token3, pair1_3, "User1")
	// 	fmt.Printf("completed: %d transactions in pool2 in one Seconds(TPS)\n", res1)
	// 	done <- true
	// }()

	// for i := 0; i < 2; i++ {
	// 	<-done
	// }

}

//testParallelExecution tests the parallel execution of transactions in the same pool
func testParallelExecution(ctx context.Context, numTx int, token1 *erc20.ERC20Token, token2 *erc20.ERC20Token, pair *core.UniswapV2Pair, user string) int {

	done := make(chan bool)
	fmt.Println("Executing transactions in parallel")

	// var mu sync.Mutex
	var swapMu sync.Mutex
	go func() {
		for i := 0; i < 15000; i++ {
			go func(i int) {

				// user := "User0"
				//approving amount for pair address
				fmt.Println("\nInitiating swap...")
				// printPairReserves(pair)

				//swpping one way
				if (i % 2) == 0 {
					fmt.Println("Swapping 1->2")

					token1.ApproveP(user, pair.GetPairAddress(), 1000)
					var path []erc20.IErcToken = []erc20.IErcToken{token1.GetIERCToken(), token2.GetIERCToken()}
					router02.SwapExactTokensForTokensP(1000, 10, path, user, user, &swapMu)

				} else {
					//swapping the other way
					fmt.Println("Swapping 2->1")

					token2.ApproveP(user, pair.GetPairAddress(), 500)
					var path []erc20.IErcToken = []erc20.IErcToken{token2.GetIERCToken(), token1.GetIERCToken()}
					router02.SwapExactTokensForTokensP(500, 10, path, user, user, &swapMu)

				}

				// printPairReserves(pair)
				done <- true
			}(i)
		}
	}()

	completed := 0
	for {
		select {
		case <-ctx.Done():
			return completed

		case <-done:
			completed++

		}
		// if completed >= numTx {
		// 	break
		// }
	}

}

//testSerialExecution tests the serial execution of transactions in the same pool
func testSerialExecution(ctx context.Context, numTx int, token1 *erc20.ERC20Token, token2 *erc20.ERC20Token, pair *core.UniswapV2Pair) int {
	fmt.Println("Executing transactions serialy")
	completed := 0

	go func() {
		for i := 0; ; i++ {

			user := "User0"
			//approving amount for pair address
			fmt.Println("\nInitiating swap...")
			// printPairReserves(pair)

			//swpping one way
			if (i % 2) == 0 {
				fmt.Println("Swapping 1->2")
				token1.Approve(user, pair.GetPairAddress(), 1000)
				var path []erc20.IErcToken = []erc20.IErcToken{token1.GetIERCToken(), token2.GetIERCToken()}
				router02.SwapExactTokensForTokens(1000, 10, path, user, user)
			} else {
				//swapping the other way
				fmt.Println("Swapping 2->1")
				token2.Approve(user, pair.GetPairAddress(), 500)
				var path []erc20.IErcToken = []erc20.IErcToken{token2.GetIERCToken(), token1.GetIERCToken()}
				router02.SwapExactTokensForTokens(500, 10, path, user, user)
			}

			// fmt.Println("Swap done!")
			printPairReserves(pair)
			completed++
		}

	}()

	//waiting for the context to finish
	for {
		<-ctx.Done()
		return completed

	}

}

//DistributeERC distributes ERC amoung 5 users
func DistributeERC(token *erc20.ERC20Token, senderAdd string, amount uint64) {
	for i := 0; i < 3; i++ { //Transferring token to user
		// user := token.TokenSymbol() + "User" + strconv.Itoa(i)
		user := "User" + strconv.Itoa(i)
		fmt.Println("Sending to: ", user)
		token.Transfer(senderAdd, user, amount)
	}
}

func printPairReserves(pair *core.UniswapV2Pair) {
	res0, res1, _ := pair.GetReserves()
	fmt.Println("Pair reserves are: ", res0, res1)
}
