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

func PrintTable(header *[]string, rows *[]*[]string) {
	colMax := make([]int, len(*header))

	// get max length of each colume value
	for _, row := range *rows {
		for i, col := range *row {
			if len(col) > colMax[i] {
				colMax[i] = len(col)
			}
			if len((*header)[i]) > colMax[i] {
				colMax[i] = len((*header)[i])
			}
		}
	}

	// build format string
	formatString := ""
	for i, colLen := range colMax {
		fs := fmt.Sprintf("%%-%ds", colLen)
		if i == 0 {
			formatString = fs
		} else {
			formatString = fmt.Sprintf("%s  %s", formatString, fs)
		}
	}
	formatString = fmt.Sprintf("%s\n", formatString)

	// print header
	headerArgs := make([]interface{}, len(*header))
	for i, v := range *header {
		headerArgs[i] = v
	}
	fmt.Printf(formatString, headerArgs...)

	// print rows
	for _, row := range *rows {
		rowArgs := make([]interface{}, len(*row))
		for i, v := range *row {
			rowArgs[i] = v
		}
		fmt.Printf(formatString, rowArgs...)
	}
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
