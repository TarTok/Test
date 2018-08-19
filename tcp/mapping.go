// Devices
package tcp

import (
	"bufio"
	//	"io"
	"log"
	"net"
	"strconv"
)

func ListenOne(port int) (net.Conn, error) {
	listen, err := net.Listen("tcp4", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}
	defer listen.Close()

	log.Printf("Begin listen port: %d", port)
	conn, err := listen.Accept()
	return conn, err

}

func Listen(port int, c chan net.Conn) {
	listen, err := net.Listen("tcp4", ":"+strconv.Itoa(port))
	if err != nil {
		log.Fatalf("Socket listen port %d failed,%s", port, err)
	}
	defer listen.Close()

	log.Printf("Begin listen port: %d", port)
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatalln(err)
			continue
		}
		select {
		case c <- conn:
		default:
			conn.Close() //второе соединение закрываем
		}
	}
}
func Mapping(rc, wc net.Conn, e chan error) {
	var (
		err error
		n   int
	)

	buf := make([]byte, 1024*100)
	r := bufio.NewReader(rc)
	w := bufio.NewWriter(wc)
	for {
		n, err = r.Read(buf)
		if err != nil {
			break
		}
		_, err = w.Write(buf[:n])
		if err != nil {
			break
		}
		//log.Println(n)
		err = w.Flush()
		if err != nil {
			break
		}

	}
	select {
	case e <- err:
	default:

	}

}
