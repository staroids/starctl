package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/staroids/starctl/pkg/api"
	"github.com/staroids/starctl/pkg/auth"
	"github.com/staroids/starctl/pkg/constants"
)

var usage = func() {
	fmt.Fprintf(os.Stdout, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

func CreateClient() *api.StaroidClient {
	auth := auth.StaroidAuth{}
	err := auth.CheckAuth()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	client := api.StaroidClient{
		Auth: auth,
	}
	return &client
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "cluster":
		ClusterCmd(os.Args[2:])
	case "version":
		fmt.Printf("%s\n", constants.Version)
		os.Exit(0)
	default:
		usage()
		os.Exit(1)
	}
}
