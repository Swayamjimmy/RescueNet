package p2p

import (
	"log"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
)

func CreateHost(PORT_HOST string) (host.Host, string, error) {

	PORT := "/ip4.127.0.0.1/tcp/" + PORT_HOST //where host will listen and returns host.Host (created host instance that can connect to others)

	h, err := libp2p.New(libp2p.ListenAddrStrings(PORT)) //libp2p.New() creates the host and ListenAddrStrings() tells host to listen to given port

	if err != nil {
		log.Fatal("Host Node Unable to Create")
		return nil, "", err
	}

	CONNECTION_STRING := PORT + "/p2p/" + h.ID().String() //full connection string that other peers need to connect to this host

	log.Println("Hello, my Peer ID is: ", h.ID())

	log.Println("Listening on: ", PORT)

	log.Println("My Host Address is: ", h.Addrs())

	return h, CONNECTION_STRING, nil

}
