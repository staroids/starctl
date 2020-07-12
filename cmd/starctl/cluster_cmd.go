package main

import (
	"flag"
	"fmt"
	"os"

	v1 "github.com/staroids/starctl/pkg/api/v1"
)

func ClusterCmdUsage() {
	fmt.Fprintf(os.Stdout, "cluster [create|list|get|delete] <name> [flags]\n")
}

func PrintClusters(cluster *[]v1.StaroidCluster, orgs *[]v1.StaroidOrg) {
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

func ClusterCmd(args []string) {
	getCmd := flag.NewFlagSet("get", flag.ExitOnError)

	if len(args) < 1 {
		ClusterCmdUsage()
		os.Exit(1)
	}

	getCmd.Parse(args)

	switch args[0] {
	case "list":
		client := CreateClient()
		orgs, err := client.V1().Org().GetAll()
		if err != nil {
			fmt.Printf("%v", err)
			os.Exit(1)
		}

		allClusters := make([]v1.StaroidCluster, 0)
		for _, org := range *orgs {
			clusters, err := client.V1().Cluster().WithOrg(org.Provider, org.Name).GetAll()

			if err != nil {
				fmt.Printf("%v", err)
			}

			allClusters = append(allClusters, *clusters...)
		}
		PrintClusters(&allClusters, orgs)
	default:
		ClusterCmdUsage()
		os.Exit(1)
	}

}
