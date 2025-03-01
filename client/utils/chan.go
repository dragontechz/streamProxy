package utils

import (
	"fmt"
	"net"
)

func MakeConn(dst string, channel chan net.Conn) {
	//channel := make(chan int)
	for {
		if len(channel) > 3 {
			fmt.Printf("greter than authorised num\n")
			continue
		}
		conn, err := net.Dial("tcp", dst)
		if err != nil {
			fmt.Printf("erro dialing to server in MakEConn func : %v", err)
			break
		}
		channel <- conn

	}

}
