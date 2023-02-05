package main

import (
	"fmt"
	"github.com/anishathalye/seashells-server/datamanager"
	"io"
	"log"
	"net"
)

func runNetcatServer(manager *datamanager.DataManager) {
	ln, err := net.Listen("tcp", ":1337")
	if err != nil {
		log.Fatalln(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go func(conn net.Conn) {
			defer conn.Close()
			remoteIp := rSplitSingle(conn.RemoteAddr().String(), ":")
			var sess *datamanager.Session
			var id string
			for sess == nil {
				// rejection sampling of random ID
				id = randomString(idLength)
				sess = manager.Create(remoteIp, id)
			}
			defer sess.Finalize()
			conn.Write([]byte(fmt.Sprintf("serving at %s%s\n", baseUrl, id)))
			buf := make([]byte, 4096)
			for {
				n, err := conn.Read(buf)
				if err != nil {
					if err != io.EOF {
						log.Printf("error while reading: %v", err)
					}
					return
				}
				ok := sess.Append(buf[:n])
				if !ok {
					conn.Write([]byte(fmt.Sprintf("error: too many connections from your ip\n")))
					return
				}
			}
		}(conn)
	}
}
