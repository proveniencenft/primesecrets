package poly

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

var field *Field

func SetField(prime *big.Int) error {
	var err error
	field = &Field{prime}
	return err
}

type Poly struct {
	Coefficients []Element
	Field        Field
}

func (p *Poly) SplitSecret(nshares int) ([]Share, error) {
	if nshares < len(p.Coefficients) {
		return nil, fmt.Errorf("Not enough shares requested")
	}
	shares := make([]Share, nshares)
	for i := range shares {
		shares[i].Point = big.NewInt(int64(i + 1))
		shares[i].Value = p.EvalAtInt(i + 1).val
		shares[i].N = p.Field.N
		shares[i].D = len(p.Coefficients)
	}
	return shares, nil

}

//The field elements are assumed NOT to be in Montgomery form
func (f *Field) NewPoly(deg int, zVal *big.Int) (*Poly, error) {
	if f.N.Cmp(big.NewInt(int64(deg))) != 1 {
		return nil, fmt.Errorf("Poly degree out of range")
	}
	p := new(Poly)
	p.Coefficients = make([]Element, deg)
	p.Coefficients[0] = *field.NewElement(zVal)
	for deg > 1 {

		e, err := field.RandomElement(rand.Reader)
		if err != nil {
			return nil, err
		}
		if e.val.Cmp(big.NewInt(0)) == 0 {
			continue
		}
		deg -= 1
		p.Coefficients[deg] = *e

	}
	p.Field = *field
	return p, nil

}

//Evaluate the polynomial at x
func (pol Poly) EvalAt(x *Element) *Element {
	res := pol.Field.NewElementInt(0)
	for power, coeff := range pol.Coefficients {
		monom := x.Clone().Exp(big.NewInt(int64(power))).Mul(&coeff)
		//fmt.Println("power, val", power, monom.val)

		res.Add(monom)
	}
	return res

}

//Evaluate the polynomial at x
func (pol Poly) EvalAtInt(x int) *Element {
	return pol.EvalAt(pol.Field.NewElementInt(x))

}

type Share struct {
	Point *big.Int
	Value *big.Int
	N     *big.Int
	D     int //one day this may be not enough of a degree
}

func Lagrange(shares []Share) (*big.Int, error) {
	if len(shares) == 0 {
		return nil, fmt.Errorf("Empty share list")
	}
	N := shares[0].N
	D := shares[0].D
	field := &Field{N}

	//Remove duplicates (And do not divide by zero)
	sh := []Share{}
	dupcheck := map[string]bool{}
	for _, s := range shares {
		if dupcheck[s.Point.String()] {
			continue
		}
		if s.N != N || s.D != D {
			return nil, fmt.Errorf("Mismatching shares: %v/%v, %v/%v", N, s.N, D, s.D)
		}
		sh = append(sh, s)
		if len(sh) > D {
			break //Enough shares already
		}

	}

	if len(sh) < D {
		return nil, fmt.Errorf("Not enough shares")
	}

	Val := field.NewElementInt(0)
	for i, s := range sh {

		dupcheck[s.Point.String()] = true
		Lacc := field.NewElement(s.Value)

		for j, s2 := range sh {
			if i == j {
				continue
			}
			num := field.NewElement(s2.Point)
			num.Neg()
			den := num.Clone()
			den.Add(field.NewElement(s.Point))
			den.Inverse()
			num.Mul(den)
			Lacc.Mul(num)

		}
		Val.Add(Lacc)

	}

	return Val.val, nil
}
