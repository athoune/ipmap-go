package main

import (
	"compress/bzip2"
	"fmt"
	"os"

	"github.com/athoune/ipmap-go/csv"
)

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	r := csv.New(bzip2.NewReader(f))
	for r.Next() {
		line, err := r.Value()
		if err != nil {
			panic(err)
		}
		fmt.Println(line)
	}
}
