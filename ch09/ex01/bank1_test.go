package bank1

import (
	"fmt"
	"sync"
	"testing"
)

func TestBank(t *testing.T) {
	done := make(chan struct{})

	go func() {
		Deposit(200)
		fmt.Println("=", Balance())
		done <- struct{}{}
	}()

	go func() {
		Deposit(100)
		done <- struct{}{}
	}()

	<-done
	<-done

	if got, want := Balance(), 300; got != want {
		t.Errorf("Balance = %d, want %d", got, want)
	}

	if Withdraw(400) {
		t.Errorf("Balance = %d", Balance())
	}

	if r, b := Withdraw(200), Balance(); !r || b != 100 {
		t.Errorf("Balance = %d", Balance())
	}

	// balance == 100
	var wg sync.WaitGroup
	wg.Add(20)
	for i := 0; i < 10; i++ {
		go func() {
			Deposit(200)
			wg.Done()
		}()
		go func() {
			Withdraw(2050)
			wg.Done()
		}()
	}
	wg.Wait()
	if b := Balance(); !(b == 2100 || b == 50) {
		t.Errorf("Balance = %d", Balance())
	}
	t.Logf("Balance = %d", Balance())
}
