package network

import (
	"fmt"
	"net"
	"sync"
)

// Peer represents a single peer in the network.
type Peer struct {
	Address string // Peer address (IP:Port)
	Conn    net.Conn
}

// P2PNetwork manages the peer-to-peer network.
type P2PNetwork struct {
	Peers    map[string]*Peer // Connected peers
	mutex    sync.Mutex       // Mutex to handle concurrent access
	listener net.Listener     // Network listener
}

// NewP2PNetwork initializes a new P2P network.
func NewP2PNetwork() *P2PNetwork {
	return &P2PNetwork{
		Peers: make(map[string]*Peer),
	}
}

// Start starts the P2P network on a given port.
func (network *P2PNetwork) Start(port string) error {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("failed to start P2P network: %v", err)
	}

	network.listener = listener
	fmt.Printf("P2P network started on port %s\n", port)

	go network.acceptConnections()

	return nil
}

// acceptConnections handles incoming peer connections.
func (network *P2PNetwork) acceptConnections() {
	for {
		conn, err := network.listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}

		go network.handleConnection(conn)
	}
}

// handleConnection manages a single peer connection.
func (network *P2PNetwork) handleConnection(conn net.Conn) {
	peer := &Peer{
		Address: conn.RemoteAddr().String(),
		Conn:    conn,
	}

	network.mutex.Lock()
	network.Peers[peer.Address] = peer
	network.mutex.Unlock()

	fmt.Printf("Connected to peer: %s\n", peer.Address)

	// Example: Reading messages from peer (basic handshake or communication)
	go func() {
		defer conn.Close()
		buf := make([]byte, 1024)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				fmt.Printf("Connection with %s closed\n", peer.Address)
				network.mutex.Lock()
				delete(network.Peers, peer.Address)
				network.mutex.Unlock()
				break
			}

			message := string(buf[:n])
			fmt.Printf("Received from %s: %s\n", peer.Address, message)
		}
	}()
}

// ConnectToPeer connects to a new peer by address.
func (network *P2PNetwork) ConnectToPeer(address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to connect to peer %s: %v", address, err)
	}

	peer := &Peer{
		Address: address,
		Conn:    conn,
	}

	network.mutex.Lock()
	network.Peers[address] = peer
	network.mutex.Unlock()

	fmt.Printf("Connected to peer: %s\n", address)
	return nil
}
