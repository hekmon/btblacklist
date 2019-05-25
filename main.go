package main

import (
	"fmt"
	"os"

	"github.com/hekmon/btblacklist/ripe"

	"github.com/kylelemons/godebug/pretty"
)

func main() {
	searches := []string{
		"trident media guard",
		"trident mediaguard",
		"trident mediguard",
		"hadopi",
	}

	results := make([][]ripe.Range, len(searches))

	for index, search := range searches {
		fmt.Println("->", search)
		ranges, err := ripe.Search(search)
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
