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

func Mul128(multiplier uint64, multiplicand uint64) (productLow uint64, productHi uint64) {
	// multiplier   = ab = a * 2^32 + b
	// multiplicand = cd = c * 2^32 + d
	// ab * cd = a * c * 2^64 + (a * d + b * c) * 2^32 + b * d
	a := multiplier >> 32
	b := multiplier & 0xFFFFFFFF
	c := multiplicand >> 32
	d := multiplicand & 0xFFFFFFFF

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
	productLow = bd + (adbc << 32)

	var productLowCarry uint64
	if productLow < bd {
		productLowCarry = 1
	} else {
		productLowCarry = 0
	}

	productHi = ac + (adbc >> 32) + (adbcCarry << 32) + productLowCarry

	return productLow, productHi
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
