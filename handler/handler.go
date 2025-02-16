package main

import (
	"flag"
	"fmt"
	//"io"
	"net"
	"stream_handler/utils"
)

var DEFAULT_RESPONSE string = "HTTP/1.1 200 OK\r\n\r\n/start/end"
var DefaultInstreamPort string = "1117"
var DefaultOutstreamPort string = "1118"

// var remoteConn_IP_map map[string]*net.Conn = make(map[string]*net.Conn)
var sessionId_remoteConn_map map[string]*net.Conn = make(map[string]*net.Conn)

//var packetId_packet_map map[string][]string = make(map[string][]string)

type handler struct {
	instreamPort, // response port
	outStreamPort, //recving port
	dstPORT string
	BUFFSIZE int
}

func main() {
	inport := flag.String("inport", DefaultInstreamPort, "set the instream port")
	outport := flag.String("outport", DefaultOutstreamPort, "set the outstream port")
	dstport := "8080"
	Buffsize := 1024 * 8
	flag.Parse()
	handler := handler{":" + *inport, ":" + *outport, ":" + dstport, Buffsize}
	handler.run()
}

func (h *handler) run() {
	go h.instreamHandler()
	h.outStreamHandler()

	//	go h.outStreamHandler()
	//	go h.instreamHandler()
}

func (h *handler) outStreamHandler() {
	listner, err := net.Listen("tcp", h.outStreamPort)
	if err != nil {
		fmt.Printf("error in oustreamhandler func: %v\n", err)
	}
	fmt.Printf("waiting for instream traffic on port: %s\n", h.outStreamPort)

	for {
		conn, err := listner.Accept()
		if err != nil {
			fmt.Printf("erro in outstream fun : %v\n", err)
		}
		go h.outstreamPacket(conn)
	}
}

func (h *handler) outstreamPacket(conn net.Conn) {
	for {
		buff := make([]byte, h.BUFFSIZE)
		n, err := conn.Read(buff)
		if err != nil {
			fmt.Printf("error in oustreampacket func : %v \n", err)
			break
		}
		if n < 5 {
			continue
		}
		data := string(buff[:n])
		req := utils.ExtractQuery(data)
		sessionId, query := utils.GetId(req)
		if sessionId == "" /* || query == ""*/ {
			fmt.Printf("error with arguement data : %s , req: %s ,query: %s ,sessionId: %s \n", data, req, query, sessionId)
			break
		}
		remote_conn, exist := sessionId_remoteConn_map[sessionId]
		if !exist {
			dstConn, err := h.makeremoteConn()
			if err != nil {
				fmt.Printf("error while mapping %v :\n", err)
				break
			}
			sessionId_remoteConn_map[sessionId] = &dstConn
			remote_conn, _ := sessionId_remoteConn_map[sessionId]

			_, err = (*remote_conn).Write([]byte(query))
			if err != nil {
				fmt.Printf("error while writing to remote conn %v : ", err)
			}
			conn.Write([]byte("HTTP/1.1 200 Byteok \r\n\r\n"))
			break

		} else {
			_, err = (*remote_conn).Write([]byte(query))
			if err != nil {
				fmt.Printf("error while writing to remote conn %v : ", err)
			}
			conn.Write([]byte("HTTP/1.1 200 Byteok \r\n\r\n"))
			break
		}
	}
}

func (h *handler) instreamHandler() {

	listner, err := net.Listen("tcp", h.instreamPort)
	if err != nil {
		fmt.Printf("error in instreamhandler func: %v\n", err)
	}
	fmt.Printf("waiting for outstream on port: %s\n", h.instreamPort)
	for {
		conn, err := listner.Accept()
		if err != nil {
			fmt.Printf("erro in intstream fun : %v\n", err)
		}
		go h.instreamPacket(conn)
	}
}

func (h *handler) instreamPacket(conn net.Conn) {
	// send response to client
	buff := make([]byte, h.BUFFSIZE)
	n, err := conn.Read(buff)
	if err != nil {
		fmt.Printf("error in instreamPacket : %v", err)
	}

	data := string(buff[:n])
	fmt.Printf("payload recved to ask for packet : %s\n", data)
	sessionId, _ := utils.GetId(utils.ExtractQuery(data))
	if sessionId == "" {
		fmt.Printf("ERROR IN instreeampacket in getting packed id\n")
	}
	go func() {
		for {
			remote_conn, exist := sessionId_remoteConn_map[sessionId]
			fmt.Printf("sessionid that ask for data %s\n", sessionId)
			fmt.Println(remote_conn)
			if !exist {
				fmt.Println("session not existing")
				continue
			}
			if remote_conn != nil {
				fmt.Println("matching sucessfull")
				forwardPacket(conn, *remote_conn)
				break
			}
		}
		//		infectedRes := DEFAULT_RESPONSE + "*/query:" + res
	}()
}

func forwardPacket(dst, src net.Conn) {
	buff := make([]byte, 1024)
	n, err := src.Read(buff)
	if err != nil {
		fmt.Printf("error in forward packet %v", err)
	}
	dst.Write(buff[:n])
	dst.Close()

}

func (h *handler) makeremoteConn() (net.Conn, error) {
	conn, err := net.Dial("tcp", h.dstPORT)
	if err != nil {
		return nil, err
	}
	return conn, err
}
