package server

import(
	"net"
	"log"
)

func StartServer() error {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:5678")
	if err != nil {
		return err
	}

	listener, err := net.ListenTCP("tcp4", tcpAddr)
	if err != nil {
		return err
	}

	go func() {
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				log.Print(err)
			}

			go RunConnection(conn)
		}
	}()
}

func RunConnection(conn *net.TCPConn) error {
	
}