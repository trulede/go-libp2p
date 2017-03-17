package libp2pquic

import (
	"crypto/tls"
	"net"

	pstore "github.com/libp2p/go-libp2p-peerstore"
	tpt "github.com/libp2p/go-libp2p-transport"
	quicconn "github.com/marten-seemann/quic-conn"
	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr-net"
)

type listener struct {
	laddr        ma.Multiaddr
	quicListener net.Listener

	transport tpt.Transport
}

func newListener(laddr ma.Multiaddr, peers pstore.Peerstore, transport tpt.Transport) (*listener, error) {
	tlsConf := &tls.Config{}
	network, host, err := manet.DialArgs(laddr)
	if err != nil {
		return nil, err
	}
	qln, err := quicconn.Listen(network, host, tlsConf)
	if err != nil {
		return nil, err
	}
	return &listener{
		laddr:        laddr,
		quicListener: qln,
		transport:    transport,
	}, nil
}

func (l *listener) Accept() (tpt.Conn, error) {
	c, err := l.quicListener.Accept()
	if err != nil {
		return nil, err
	}

	mnc, err := manet.WrapNetConn(c)
	if err != nil {
		return nil, err
	}

	return &tpt.ConnWrap{
		Conn: mnc,
		Tpt:  l.transport,
	}, nil
}

func (l *listener) Close() error {
	return l.quicListener.Close()
}

func (l *listener) Addr() net.Addr {
	return l.quicListener.Addr()
}

func (l *listener) Multiaddr() ma.Multiaddr {
	return l.laddr
}

var _ tpt.Listener = &listener{}
