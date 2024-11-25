package transaction

import "fmt"

func ApplyTax(amount int64, taxRate int64) int64 {
	tax := amount * taxRate / 100
	netAmount := amount - tax
	return netAmount
}

func TransactionWithTax(from string, to string, amount int64, taxRate int64) {
	netAmount := ApplyTax(amount, taxRate)
	fmt.Printf("Transferring %d from %s to %s with tax %d applied\n", netAmount, from, to, amount-netAmount)
}
