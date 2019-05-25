package main

import (
	"fmt"
	"os"

	"github.com/kylelemons/godebug/pretty"

	"github.com/hekmon/btblacklist/ripe"
)

func main() {
	ripeController := ripe.New()

	searches := []string{
		"trident media guard",
		"trident mediaguard",
		"trident mediguard",
		"hadopi",
	}

	for _, search := range searches {
		fmt.Println(search)
		data, err := ripeController.Search(search)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		pretty.Print(data)
		fmt.Println()
	}
}
