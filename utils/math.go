package utils

import "sort"

func MedianSlice(sl []uint64) uint64 {
	l := len(sl)

	if l == 0 {
		return 0
	}

	if l == 1 {
		return sl[0]
	}

	sort.Slice(sl, func(i, j int) bool { return sl[i] < sl[j] })

	n := l / 2
	if l%2 == 0 {
		return (sl[n-1] + sl[n]) / 2
	} else {
		return sl[n]
	}
}

func ClampInt64(lo int64, v int64, hi int64) int64 {
	if v < lo {
		return lo
	}

	if v > hi {
		return hi
	}

	return v
}

func Mul128(multiplier uint64, multiplicand uint64) (productLo uint64, productHi uint64) {
	// multiplier   = ab = a * 2^32 + b
	// multiplicand = cd = c * 2^32 + d
	// ab * cd = a * c * 2^64 + (a * d + b * c) * 2^32 + b * d
	a := hiDword(multiplier)
	b := loDword(multiplier)
	c := hiDword(multiplicand)
	d := loDword(multiplicand)

	ac := a * c
	ad := a * d
	bc := b * c
	bd := b * d

	adbc := ad + bc

	var adbcCarry uint64
	if adbc < ad {
		adbcCarry = 1
	} else {
		adbcCarry = 0
	}

	// multiplier * multiplicand = product_hi * 2^64 + product_lo
	productLo = bd + (adbc << 32)

	var productLowCarry uint64
	if productLo < bd {
		productLowCarry = 1
	} else {
		productLowCarry = 0
	}

	productHi = ac + (adbc >> 32) + (adbcCarry << 32) + productLowCarry

	return productLo, productHi
}

// Div128p32 long division with 2^32 base
func Div128p32(dividendHi, dividendLo uint64, divisor uint32) (quotientHi, quotientLo uint64) {
	dividendDwords := make([]uint64, 4)
	dividendDwords[3] = hiDword(dividendHi)
	dividendDwords[2] = loDword(dividendHi)
	dividendDwords[1] = hiDword(dividendLo)
	dividendDwords[0] = loDword(dividendLo)

	remainder := uint32(0)
	quotientHi = divWithRemained(dividendDwords[3], divisor, &remainder) << 32
	quotientHi |= divWithRemained(dividendDwords[2], divisor, &remainder)
	quotientLo = divWithRemained(dividendDwords[1], divisor, &remainder) << 32
	quotientLo |= divWithRemained(dividendDwords[0], divisor, &remainder)

	return quotientLo, quotientHi
}

func MinInt64(a int64, b int64) int64 {
	if a > b {
		return b
	}

	return a
}

func MinUint64(a uint64, b uint64) uint64 {
	if a > b {
		return b
	}

	return a
}

func MaxInt64(a int64, b int64) int64 {
	if a > b {
		return a
	}

	return b
}

func MaxUint64(a uint64, b uint64) uint64 {
	if a > b {
		return a
	}

	return b
}

func divWithRemained(dividend uint64, divisor uint32, remainder *uint32) uint64 {
	dividend |= (uint64)(*remainder) << 32
	*remainder = uint32(dividend % uint64(divisor))
	return dividend / uint64(divisor)
}

func hiDword(val uint64) uint64 {
	return val >> 32
}

func loDword(val uint64) uint64 {
	return val & 0xFFFFFFFF
}
