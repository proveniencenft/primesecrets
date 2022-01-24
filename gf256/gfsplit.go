package gf256

import (
	"crypto/rand"
	"fmt"
)

type gfpoly struct {
	Coeff []byte
}

type Share struct {
	Point  byte
	Value  []byte
	Degree byte
}

func (p *gfpoly) eval(pt byte) byte {
	v := p.Coeff[0]
	for i := 1; i < len(p.Coeff); i++ {
		v ^= Mul(p.Coeff[i], Exp(pt, i))

	}
	return v
}

func newgfpoly(zerovalue byte, degree int) *gfpoly {
	if degree < 0 || degree > 254 {
		panic("Nonsensical poly degree")
	}
	c := make([]byte, degree+1)
	for c[degree] == 0 {
		rand.Read(c)
	}

	c[0] = zerovalue
	return &gfpoly{c}
}

func SplitBytes(tosplit []byte, nshares, threshold int) ([]Share, error) {
	if nshares <= 0 || nshares > 254 || threshold > nshares {
		fmt.Errorf("wrong No of shares or threshold")
	}
	shares := make([]Share, nshares)
	for i := 0; i < nshares; i++ {
		shares[i].Point = byte(i + 1)
		shares[i].Degree = byte(threshold - 1)
	}
	for j := range tosplit {
		c := newgfpoly(tosplit[j], threshold-1)
		for i := 0; i < nshares; i++ {
			//vi := c.eval(byte(i + 1))
			shares[i].Value = append(shares[i].Value, c.eval(byte(i+1)))
		}

	}
	return shares, nil
}

func RecoverBytes(shares []Share) ([]byte, error) {
	deg := shares[0].Degree
	leng := len(shares[0].Value)

	//deduplicate shares
	dup := map[byte]bool{}
	unique := []Share{}
	for _, s := range shares {
		if dup[s.Point] {
			continue
		}
		if s.Degree != deg || len(s.Value) != leng {

			return nil, fmt.Errorf("Inconsistent shares")
		}
		unique = append(unique, s)
		dup[s.Point] = true

	}
	ideg := int(deg)
	if len(unique) <= ideg {

		return nil, fmt.Errorf("Not enough shares")
	}
	result := []byte{}
	for k := 0; k < len(unique[0].Value); k++ {
		v := byte(0)
		for i := 0; i <= int(deg); i++ {
			a := shares[i].Value[k]
			//fmt.Println("a:", a)
			for j := 0; j <= ideg; j++ {
				if i == j {
					continue
				}
				a = Mul(a, Mul(shares[j].Point, Inv(shares[i].Point^shares[j].Point)))

			}
			//fmt.Println("L:", a)
			v ^= a

		}
		result = append(result, v)
	}

	return result, nil

}
