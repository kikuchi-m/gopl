package bank1

var deposits = make(chan int)
var done = make(chan bool)
var balances = make(chan int)

func Deposit(amount uint) {
	deposits <- int(amount)
	<-done
}

func Withdraw(amount uint) bool {
	deposits <- -int(amount)
	return <-done
}

func Balance() int {
	return <-balances
}

func teller() {
	var balance int
	for {
		select {
		case amount := <-deposits:
			balance += amount
			if balance < 0 {
				balance += -amount
				done <- false
			} else {
				done <- true
			}
		case balances <- balance:
		}
	}
}

func init() {
	go teller()
}
