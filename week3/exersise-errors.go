package main

import (
	"fmt"
)

type ErrNegativeSqrt float64

func (e ErrNegativeSqrt) Error() string {
	return fmt.Sprintf("cannot Sqrt negative number: %v", float64(e))
}

func Sqrt(x float64) (float64, error) {
	z := 1.0

	switch {
	case x > 0:
		for i := 0; i < 10; i++ {
			z -= (z*z - x) / (2 * z)
		}
	case x == 0:
		return 0, nil
	default:
		return 0, ErrNegativeSqrt(x)
	}

	return z, nil
}

func main() {
	fmt.Println(Sqrt(2))
	fmt.Println(Sqrt(0))
	fmt.Println(Sqrt(-2))
}
