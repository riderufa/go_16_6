package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type BankClient struct {
	dataMutex sync.Mutex
	balance   int
}

func (c *BankClient) Deposit(amount int) {
	c.dataMutex.Lock()
	c.balance += amount
	c.dataMutex.Unlock()
}

func (c *BankClient) Withdrawal(amount int) error {
	c.dataMutex.Lock()
	defer c.dataMutex.Unlock()
	c.balance -= amount
	if c.balance < 0 {
		c.balance += amount
		return fmt.Errorf("balance is not enough")
	}
	return nil
}

func (c *BankClient) Balance() int {
	c.dataMutex.Lock()
	defer c.dataMutex.Unlock()
	return c.balance
}

func main() {
	var bankClient BankClient
	var command string
	var amount int

	for i := 0; i < 10; i++ {
		go func() {
			for {
				rDeposit := rand.Intn(10) + 1
				bankClient.Deposit(rDeposit)
				rDuration := rand.Intn(1000) + 500
				time.Sleep(time.Duration(rDuration) * time.Millisecond)
			}
		}()
	}

	for i := 0; i < 5; i++ {
		go func() {
			for {
				rWithdrawal := rand.Intn(5) + 1
				_ = bankClient.Withdrawal(rWithdrawal)
				rDuration := rand.Intn(1000) + 500
				time.Sleep(time.Duration(rDuration) * time.Millisecond)
			}
		}()
	}

	fmt.Println("You can use commands: balance, deposit, withdrawal, exit")
	for {
		_, err := fmt.Scanln(&command)
		if err != nil {
			fmt.Println("Выход")
			return
		}
		switch command {
		case "balance":
			fmt.Println(bankClient.Balance())
		case "deposit":
			fmt.Println("Enter deposit amount")
			_, err := fmt.Scanln(&amount)
			if err != nil {
				fmt.Println("Unsupported type")
			} else {
				bankClient.Deposit(amount)
			}
		case "withdrawal":
			fmt.Println("Enter withdrawal amount")
			_, err := fmt.Scanln(&amount)
			if err != nil {
				fmt.Println("Unsupported type")
			} else {
				err = bankClient.Withdrawal(amount)
				if err != nil {
					fmt.Println(err)
				}
			}
		case "exit":
			return
		default:
			fmt.Println("Unsupported command. You can use commands: balance, deposit, withdrawal, exit")
		}
	}
}
