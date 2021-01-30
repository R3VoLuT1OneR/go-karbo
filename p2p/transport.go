package p2p

//import (
//	"context"
//	"github.com/libp2p/go-libp2p-core/transport"
//	ma "github.com/multiformats/go-multiaddr"
//	"github.com/libp2p/go-libp2p-core/peer"
//	manet "github.com/multiformats/go-multiaddr/net"
//	"time"
//)
//
//type P2PTransport struct {
//
//	ConnectTimeout time.Duration
//
//}
//
//func NewTransport() transport.Transport {
//	t := P2PTransport{}
//
//	return t
//}
//
//func (t *P2PTransport) Dial(ctx context.Context, raddr ma.Multiaddr, p peer.ID) (transport.CapableConn, error) {
//	if t.ConnectTimeout > 0 {
//		deadline := time.Now().Add(t.ConnectTimeout)
//		if d, ok := ctx.Deadline(); !ok || deadline.Before(d) {
//			var cancel func()
//			ctx, cancel = context.WithDeadline(ctx, deadline)
//			defer cancel()
//		}
//	}
//
//	var d manet.dialer
//	conn, err := d.DialContext(ctx, raddr)
//	if err != nil {
//		return nil, err
//	}
//
//	caconn := transport.CapableConn()
//}
//
//func (t *P2PTransport) CanDial(addr ma.Multiaddr) bool {
//
//}
//
//func (t *P2PTransport) Listen(laddr ma.Multiaddr) (transport.listener, error) {
//
//}
//
//func (t *P2PTransport) Protocols() []int {
//
//}
//
//func (t *P2PTransport) Proxy() bool {
//	return false
//}