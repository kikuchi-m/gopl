package bank1

type req struct {
	amount  int
	success chan<- bool
}

var deposits = make(chan req)
var balances = make(chan int)

func Deposit(amount uint) {
	done := make(chan bool)
	deposits <- req{int(amount), done}
	// always succeed
	<-done
}

func Withdraw(amount uint) bool {
	done := make(chan bool)
	deposits <- req{-int(amount), done}
	return <-done
}

func Balance() int {
	return <-balances
}

func teller() {
	var balance int
	for {
		select {
		case req := <-deposits:
			balance += req.amount
			if balance < 0 {
				balance += -req.amount
				req.success <- false
			} else {
				req.success <- true
			}
		case balances <- balance:
		}
	}
}

func init() {
	go teller()
}
