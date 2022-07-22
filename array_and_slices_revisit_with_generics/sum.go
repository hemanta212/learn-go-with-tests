package main

func Sum(numbers []int) int {
	return Reduce(numbers, func(sum, item int) int {
		return sum + item
	}, 0)
}

func SumAll(numbersToSum ...[]int) []int {
	return Reduce(numbersToSum, func(numbers []int, item []int) []int {
		return append(numbers, Sum(item))
	}, []int{})
}

func SumAllTails(numbersToSum ...[]int) []int {
	return Reduce(numbersToSum, func(numbers []int, item []int) []int {
		if len(item) == 0 {
			return append(numbers, 0)
		} else {
			return append(numbers, Sum(item[1:]))
		}
	}, []int{})
}

func Reduce[A, B any](collection []A, accumulator func(B, A) B, initialValue B) B {
	var result = initialValue
	for _, x := range collection {
		result = accumulator(result, x)
	}
	return result
}

type Transaction struct {
	From, To string
	Sum      int
}

func BalanceFor(transactions []Transaction, name string) int {
	adjustBalance := func(currBalance int, transaction Transaction) int {
		if transaction.From == name {
			return currBalance - transaction.Sum
		} else if transaction.To == name {
			return currBalance + transaction.Sum
		}
		return currBalance
	}
	return Reduce(transactions, adjustBalance, 0)
}
