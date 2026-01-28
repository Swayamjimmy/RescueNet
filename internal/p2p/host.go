package p2p

import (
	"fmt"
	"log"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
)

// CreateHost creates a libp2p host bound to the provided TCP port.
// Returns the host, its connection string, or an error.
func CreateHost(PORT_HOST string) (host.Host, string, error) {
	if PORT_HOST == "" {
		return nil, "", fmt.Errorf("port must be provided")
	}

	// Correct multiaddr format for ipv4 TCP listener.
	listenAddr := "/ip4/0.0.0.0/tcp/" + PORT_HOST

	h, err := libp2p.New(libp2p.ListenAddrStrings(listenAddr))
	if err != nil {
		return nil, "", fmt.Errorf("unable to create host on %s: %w", listenAddr, err)
	}

	connectionString := listenAddr + "/p2p/" + h.ID().String() // full connection string peers need to dial

	log.Println("Hello, my Peer ID is: ", h.ID())
	log.Println("Listening on: ", listenAddr)
	log.Println("My Host Address is: ", h.Addrs())

	return h, connectionString, nil
}
