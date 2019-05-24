package main

import (
	"fmt"
	"os"

	"github.com/kylelemons/godebug/pretty"

	"github.com/hekmon/btblacklist/ripe"
)

func main() {
	ripeController := ripe.New()
	data, err := ripeController.Search("trident mediaguard")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	pretty.Print(data)
}
