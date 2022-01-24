package gf256

import (
	"bytes"
	"crypto/rand"
	"testing"
)

func TestMul(t *testing.T) {
	//fmt.Println(gexp)
	//fmt.Println(Exp(0, 3))
	b := make([]byte, 3)
	for i := 0; i < 100; i++ {

		rand.Read(b)
		if assoc(b[0], b[1], b[2]) != 0 {
			t.Error("Not associative!")
		}

		if b[0] != 0 && Mul(b[0], Inv(b[0])) != 1 {
			t.Error("Wrong inverse!", b[0])
		}

	}

	secret := make([]byte, 4100)
	rand.Read(secret)
	sh, _ := SplitBytes(secret, 4, 3)
	rec, err := RecoverBytes(sh[1:])
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(secret, rec) {
		t.Errorf("Recovered wrong plaintext")
	}

}

func assoc(x, y, z byte) byte {
	f1 := Mul(x, z) ^ Mul(y, z)
	f2 := Mul((x ^ y), z)
	v := f1 ^ f2
	return v
}
