package poly

import (
	"fmt"
	"io"
	"math/big"
)

type Field struct {
	N *big.Int
}

type Element struct {
	val   *big.Int
	field *Field
}

func (e *Element) Mul(e2 *Element) *Element {
	if e2.val == nil {
		fmt.Println("BadBoy")
	}
	t := new(big.Int)
	t.Mul(e.val, e2.val)
	e.val.Mod(t, e.field.N)
	return e
}

func (e *Element) Clone() *Element {
	ne := &Element{new(big.Int), e.field}
	ne.val.Set(e.val)
	return ne
}

func (e *Element) Add(e2 *Element) *Element {
	e.val.Add(e.val, e2.val)
	e.val.Mod(e.val, e.field.N)
	return e
}

func (e *Element) Exp(i *big.Int) *Element {
	e.val.Exp(e.val, i, e.field.N)
	return e
}

func (f *Field) NewElement(i *big.Int) *Element {
	e := &Element{new(big.Int), f}
	e.val.Mod(i, f.N)
	return e
}

func (f *Field) NewElementInt(i int) *Element {
	return f.NewElement(big.NewInt(int64(i)))
}

func (f *Field) RandomElement(rd io.Reader) (*Element, error) {
	buf := make([]byte, len(f.N.Bytes()))
	_, err := rd.Read(buf)
	if err != nil {
		return nil, err
	}
	v := new(big.Int)
	v.SetBytes(buf)
	v = v.Mod(v, f.N)
	e := &Element{v, f}
	return e, nil
}

func (e *Element) Inverse() *Element {
	e.val.ModInverse(e.val, e.field.N)
	return e
}

func (e *Element) Neg() *Element {
	e.val.Neg(e.val)
	e.val.Mod(e.val, e.field.N)
	return e
}
