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

	results := make([][]ripe.Range, len(searches))

	for index, search := range searches {
		fmt.Println("->", search)
		ranges, err := ripeController.Search(search)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		pretty.Print(ranges)
		results[index] = ranges
		fmt.Println()
	}

	fmt.Println("-> Uniq ranges")
	uniq := ripe.RemoveDuplicates(results)
	pretty.Print(uniq)
}
