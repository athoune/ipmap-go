package main

import (
	"bufio"
	"compress/bzip2"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/athoune/ipmap-go/csv"
	"github.com/athoune/ipmap-go/ipmap"
)

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	reader := csv.New(bzip2.NewReader(f))
	ranges, err := ipmap.New(reader)
	if err != nil {
		panic(err)
	}
	fmt.Println("file read, lets search")

	listen, err := net.Listen("tcp", "0.0.0.0:1234")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go func(conn net.Conn) {
			scan := bufio.NewScanner(conn)
			defer conn.Close()
			for scan.Scan() {
				line := scan.Text()
				line = strings.TrimSpace(line)
				log.Println(line)
				if line == "" {
					continue
				}
				containingNetworks, err := ranges.Ranger.ContainingNetworks(net.ParseIP(line))
				if err != nil {
					log.Printf("[%s] : %s\n", line, err)
					return
				}
				for _, network := range containingNetworks {
					fmt.Fprintf(conn, "%s\n", network)
				}
			}
		}(conn)

	}
}
