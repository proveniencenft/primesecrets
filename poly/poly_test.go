package poly

import (
	"fmt"
	"math/big"
	"testing"
)

func TestNewPoly(t *testing.T) {
	n, _ := new(big.Int).SetString("73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001", 16)
	field := &Field{n}

	pol, err := field.NewPoly(15, big.NewInt(17000000042))

	if err != nil {
		t.Error(err)
	}

	for i, v := range pol.Coefficients {
		fmt.Println(i, v.val)
	}

	shares, err := pol.GenerateShares(20)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(shares)
	v, err := Lagrange(shares[:])
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Decoded secret:", v)
	v, err = Lagrange(shares[1:])
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Decoded secret:", v)
	v, err = Lagrange(shares[2:])
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Decoded secret:", v)

}
