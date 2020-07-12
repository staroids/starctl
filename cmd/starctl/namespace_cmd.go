package main

import (
	"flag"
	"fmt"
	"os"

	v1 "github.com/staroids/starctl/pkg/api/v1"
)

func NamespaceCmdUsage() {
	fmt.Fprintf(os.Stdout, "namespace [create|list|get|star|stop|delete] <alias> [flags]\n")
}

func PrintNamespaces(cluster *[]v1.StaroidCluster, orgs *[]v1.StaroidOrg) {
	orgInfo := make(map[int64]*v1.StaroidOrg)
	if orgs != nil {
		for _, org := range *orgs {
			orgInfo[org.ID] = &org
		}
	}

	rows := make([]*[]string, 0)
	for _, cluster := range *cluster {
		org := orgInfo[cluster.OrgID]
		rows = append(rows, &[]string{cluster.Name, fmt.Sprintf("%s/%s", org.Provider, org.Name), fmt.Sprintf("%s/%s", cluster.Ske.Cloud, cluster.Ske.Region)})
	}
	header := []string{"NAME", "ORG", "SKE"}
	PrintTable(&header, &rows)
}

func NamespaceCmd(args []string) {
	getCmd := flag.NewFlagSet("namespace", flag.ExitOnError)

	if len(args) < 1 {
		NamespaceCmdUsage()
		os.Exit(1)
	}

	getCmd.Parse(args)

	switch args[0] {
	case "list":
	default:
		NamespaceCmdUsage()
		os.Exit(1)
	}

}
