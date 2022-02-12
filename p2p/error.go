package p2p

import "errors"

var (
	ErrSyncDataTooDeepBehind = errors.New("top block too deep behind")
)
