package main

import (
	"bufio"
	"compress/bzip2"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/athoune/ipmap-go/csv"
	"github.com/athoune/ipmap-go/ipmap"
)

func main() {
	chrono := time.Now()
	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	reader := csv.New(bzip2.NewReader(f))
	ranges, err := ipmap.New(reader)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d lines read in %v, lets search\n", ranges.Length(), time.Now().Sub(chrono))

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
				chrono := time.Now()
				containingNetworks, ok := ranges.Tree.Get(net.ParseIP(line))
				log.Printf("%v", time.Now().Sub(chrono))
				if !ok {
					fmt.Fprintf(conn, "%s => nope\n", line)
				} else {
					fmt.Fprintf(conn, "%s %v\n", line, containingNetworks)
				}
			}
		}(conn)

	}
}
