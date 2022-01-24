package gf256

var gexp = make([]byte, 256)
var glog = make([]byte, 256)

const gen = 0x03

const poly = 0x011B

func mulInt(x, y int) int {
	y0 := 0
	for x > 0 {
		if x&1 == 0x01 {
			y0 ^= y

		}
		y <<= 1
		x >>= 1
		if y0&0x100 > 0 {
			y0 ^= poly
		}

	}
	return y0 & 0xff
}

func init() {
	el := 1
	for i := 0; i < 256; i++ {
		gexp[i] = byte(el)
		glog[el] = byte(i % 255)
		el = mulInt(el, gen)

	}
}

func Mul(x, y byte) byte {
	if x == 0 || y == 0 {
		return 0
	}
	s := int(glog[x]) + int(glog[y])
	if s > 255 {
		s = (s + 1) % 256
	}
	return gexp[s]
}

func Inv(x byte) byte {
	return gexp[255-glog[x]]
}

func Exp(x byte, p int) byte {
	if p == 0 {
		return 1
	}
	if x == 0 {
		return x
	}
	p = p % 255
	s := int(glog[x]) * p
	if s > 255 {
		s = s%256 + 1
	}
	return gexp[s]

}
