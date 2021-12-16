package cryptonote

type chain struct {
	Blocks []*Block

	StartIndex int32
	EndIndex   int32

	Parent *Chain
}

type Chain interface {
}

func NewChain() (Chain, error) {
	chain := &chain{}

	return chain, nil
}
