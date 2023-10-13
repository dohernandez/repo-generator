package main

import (
	"flag"

	_ "github.com/ethereum/go-ethereum/common"
	_ "github.com/google/uuid"
	_ "github.com/lib/pq"
)

var (
	outputPath = flag.String("output", "", "Path of the output file")
	model      = flag.String("model", "", "Name of struct to generate repo for")
	tag        = flag.String("tag", "", "Tag to use for the generated repo")
	//tag        = flag.StringArray("tag", "", "Tag to use for the generated repo")
)

func main() {

}
