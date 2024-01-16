package main

import (
	"errors"
)

// Time complexity: O(n)
// Space complexity: O(1)
func sum_to_n_a(n int) (int, error) {
	if n < 0 {
		return 0, errors.New("n expected to be non-negative")
	}

	sum := 0
	for i := 1; i <= n; i++ {
		sum += i
	}
	return sum, nil
}

// Time complexity: O(n)
// Space complexity: O(1) with tail call elimination
func sum_to_n_b(n int) (int, error) {
	if n < 0 {
		return 0, errors.New("n expected to be non-negative")
	}
	return sum_to_n_rec(0, n), nil
}

func sum_to_n_rec(accumulator int, n int) int {
	if n == 0 {
		return accumulator
	}
	return sum_to_n_rec(accumulator+n, n-1)
}

// Time complexity: O(1)
// Space complexity: O(1)
func sum_to_n_c(n int) (int, error) {
	if n < 0 {
		return 0, errors.New("n expected to be non-negative")
	}
	return (n * (n + 1)) / 2, nil
}
