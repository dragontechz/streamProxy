package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"stream/utils"
	//"strings"
)

const (
	SESSIONIDLENGTH = 8
)

var DEFAULT_PAYLOAD string = `POST /chat HTTP/1.1 
Host: c.whatsapp.net 
User-Agent: Mozilla/5.0 (compatible; WAChat/1.2; +http://www.whatsHTTP/1.0app.com/contact)\start\end`
var DefaultInstreamPort string = "1118"
var DefaultOutstreamPort string = "1117"
var ADDR string = "" //"170.205.31.126"
var RECV_buff int = 1024 * 65

type packet struct {
	buff []byte
	n    int
}
type server struct {
	servingPort,
	instreamAddr,
	payload,
	outstreamAddr string
}

var channeloutstreamaddr chan net.Conn = make(chan net.Conn)
var channelintstreamaddr chan net.Conn = make(chan net.Conn)

func main() {
	inport := flag.String("inport", DefaultInstreamPort, "set the instream port")
	outport := flag.String("outport", DefaultOutstreamPort, "set the outstream port")
	payload := flag.String("payload", DEFAULT_PAYLOAD, "set the payload")
	serving_port := ":9090"
	flag.Parse()
	server := server{serving_port, ADDR + ":" + *inport, *payload, ADDR + ":" + *outport}
	server.Run()
}

func (s *server) Run() {
	listener, err := net.Listen("tcp", s.servingPort)
	fmt.Printf("client started on port %s waiting for connection to stream data\n", s.servingPort)
	if err != nil {
		fmt.Println("ERROR listening: ", err)
	}

	go utils.MakeConn(s.outstreamAddr, channeloutstreamaddr)
	go utils.MakeConn(s.instreamAddr, channelintstreamaddr)
	for {

		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("error in accepting: ", err)
		}
		fmt.Println("new connection")
		go s.handle_traffic(conn)
	}
}

// handles incommimg and outgoing traffic
func (s *server) handle_traffic(conn net.Conn) {
	sessionid := utils.GenerateRandomString(SESSIONIDLENGTH)
	go s.inSTREAM(conn, sessionid)
	go s.outSTREAM(conn, sessionid)
	//go s.outSTREAM(conn, sessionid)

}
func (s *server) inSTREAM(conn net.Conn, sessionid string) {
	defer conn.Close()
	//	channel := make(chan net.Conn)
	res_chan := make(chan string)
	//	go utils.MakeConn(s.outstreamAddr, channel)
	go func() {

		for {
			dst := <-channeloutstreamaddr
			req := utils.InsertQuery(s.payload, sessionid+":")

			_, err := dst.Write([]byte(req))
			if err != nil {
				fmt.Println("error in istream func : \n", err)
			}

			buff := make([]byte, RECV_buff)
			n, err := dst.Read(buff)
			if err != nil {
				if err == io.EOF {
				}
				fmt.Printf("cann't read from server: %v\n", err)
				dst.Close()
				break
			}
			data := string(buff[:n])

			//fmt.Printf("before extraction : %s\n", data)
			res := utils.ExtractQuery(data)
			fmt.Printf("before decompression %f kb\n", float64(len(res))/1024.0)
			res, err = utils.Decompress_str(res)
			if err != nil {
				fmt.Printf("error decompressing: %v\n", err)
			}
			fmt.Printf("after decompression %f kb\n", float64(len(res))/1024.0)
			res_chan <- res

			dst.Close()
			fmt.Printf("continuing\n")

			//conn.Write([]byte(res))
			continue

		}
	}()

	for {
		res := <-res_chan
		_, err := conn.Write([]byte(res))
		if err != nil {
			fmt.Printf("error in recv channel: %v\n", err)
			conn.Close()
		}
	}
}

func (s *server) outSTREAM(conn net.Conn, sessionid string) {
	for {
		buff := make([]byte, 1024*5)

		n, err := conn.Read(buff)
		if err != nil {
			if err == io.EOF {
				break
			}
			//fmt.Printf("error in outstream func from client: %v\n", err)
			//continue
		}
		if n < 1 {
			continue
		}
		dst := <-channelintstreamaddr
		req := string(buff[:n])
		maskedQuery := utils.InsertQuery(s.payload, sessionid+":"+utils.Compress_str(req))

		_, err = dst.Write([]byte(maskedQuery))
		if err != nil {
			fmt.Println("error in istream func : ", err)
			continue

		}
		for {
			_, err = dst.Read(buff)
			if err != nil {
				if err == io.EOF {
					fmt.Printf("coonnection was closed because of %v", err)
				}
				fmt.Printf("error in outstream func : %v", err)
				break
			}
			//if n < 1 {
			//	continue
			//}
			//res := string(buff[:n])
			//fmt.Printf("%s\n", res)
			break
		}

	}
}

// return a conn
/*func (s *server) makeConn(addr string) (net.Conn, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, nil
	}
	return conn, err
}*/

/*func forward_to_client(dst, client net.Conn) {
	for {
		buff := make([]byte, 1024*3)

		n, err := dst.Read(buff)
		if err != nil {
			fmt.Println("error in forward to client func: ", err)
			dst.Close()
			break
		}
		//if n <= 4 {
		//	continue
		//}
		data := string(buff[:n])
		res := extractRes(data)
		fmt.Printf("recved:%s", res)
		_, err = client.Write([]byte(res))
		if err != nil {
			fmt.Println("error in forwrd to client func:", err)
			break
		}
		//break
	}
}*/
