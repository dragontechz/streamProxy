package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"stream_handler/utils"
	"sync"
)

var DEFAULT_RESPONSE string = "HTTP/1.1 200 OK\r\n"
var DefaultInstreamPort string = "1117"
var DefaultOutstreamPort string = "1118"

// var remoteConn_IP_map map[string]*net.Conn = make(map[string]*net.Conn)
var sessionId_remoteConn_map map[string]*net.Conn = make(map[string]*net.Conn)

var SESSIONID_REMOTE_CONN_map sync.Map

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
	go handler.run()
	proxy := utils.SOCKS5{Listn_addr: ":" + dstport}
	proxy.RUN_v5()
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

	go func() {
		for {
			conn, err := listner.Accept()
			if err != nil {
				fmt.Printf("erro in outstream fun : %v\n", err)
			}
			go h.outstreamPacket(conn)
		}
	}()
	for {
		conn, err := listner.Accept()
		if err != nil {
			fmt.Printf("erro in outstream fun : %v\n", err)
		}
		go h.outstreamPacket(conn)
	}
}

func (h *handler) outstreamPacket(conn net.Conn) {
	channel := make(chan net.Conn)
	go utils.MakeConn(h.dstPORT, channel)
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
		dst, exist := SESSIONID_REMOTE_CONN_map.Load(sessionId)
		if !exist {
			dstConn := <-channel
			SESSIONID_REMOTE_CONN_map.Store(sessionId, &dstConn)
			dst, _ := SESSIONID_REMOTE_CONN_map.Load(sessionId)
			remote_conn := dst.(*net.Conn)

			query, err = utils.Decompress_str(query)

			if err != nil {
				fmt.Printf("error in decompression: %v\n", err)
			}

			//fmt.Printf("decompressed query : %s\n", query)
			_, err = (*remote_conn).Write([]byte(query)) //
			if err != nil {
				fmt.Printf("error while writing to remote conn %v : ", err)
			}
			conn.Write([]byte("HTTP/1.1 200 Byteok \r\n\r\n"))
			break

		} else {
			remote_conn := dst.(*net.Conn)
			query, err = utils.Decompress_str(query)
			if err != nil {
				fmt.Printf("error in outstream packet in decompression: %v\n", err)
			}
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
		fmt.Printf("error in instreamPacket : %v\n", err)
	}

	data := string(buff[:n])
	sessionId, _ := utils.GetId(utils.ExtractQuery(data))
	if sessionId == "" {
		fmt.Printf("ERROR IN instreeampacket in getting packed id\n")
	}
	go func() {
		//fmt.Printf("sessionid that ask for data %s\n", sessionId)
		for {
			remoteConnInterface, _ := SESSIONID_REMOTE_CONN_map.Load(sessionId)

			//if !exist {
			//	continue
			//}

			remoteConn, ok := remoteConnInterface.(*net.Conn)
			if ok && remoteConn != nil {
				fmt.Printf("matching sessionid :%s\n", sessionId)
				go h.forwardPacket(conn, (*remoteConn), sessionId)
				break
			}
		}
	}()
}

func (h *handler) forwardPacket(dst, src net.Conn, sessionid string) {
	defer dst.Close()
	buff := make([]byte, h.BUFFSIZE*8)
	n, err := src.Read(buff)
	if err != nil {
		if err == io.EOF {
			SESSIONID_REMOTE_CONN_map.Delete(sessionid)
			fmt.Printf("session: %s of ip %s has been closed\n", sessionid, dst.RemoteAddr().String())
			src.Close()
		}
		fmt.Printf("error in forward packet %v\n", err)
		return
	}
	data := string(buff[:n])
	fmt.Printf("forwar packet of value: %f kb\n", float64(len(data))/1024.0)
	res := utils.InsertQuery(DEFAULT_RESPONSE, utils.Compress_str(data))

	dst.Write([]byte(res))
}
