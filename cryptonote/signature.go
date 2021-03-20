package cryptonote

type EllipticCurvePointer [32]byte

type EllipticCurveScalar [32]byte

type Signature struct {
	C, R EllipticCurveScalar
}