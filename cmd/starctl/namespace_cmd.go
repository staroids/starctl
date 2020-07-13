package main

import (
	"flag"
	"fmt"
	"os"

	v1 "github.com/staroids/starctl/pkg/api/v1"
)

func NamespaceCmdUsage() {
	fmt.Fprintf(os.Stdout, "namespace [create|list|get|start|stop|delete] <alias> [flags]\n")
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
	namespaceCmdFlag := flag.NewFlagSet("namespace", flag.ExitOnError)
	orgName := namespaceCmdFlag.String("org", "", "organization (e.g. GITHUB/staroid)")
	clusterName := namespaceCmdFlag.String("cluster", "", "name of cluster")
	commitLoc := namespaceCmdFlag.String("project", "GITHUB/staroids/namespace:master", "project:branch(#commit) (e.g. GITHUB/staroid/app:master, GITHUB/staroid/app:trunk#d10abcd)")

	namespaceCmdFlag.Parse(args)

	if *orgName == "" {
		fmt.Println("'org' flag is missing")
		os.Exit(1)
	}

	if *clusterName == "" {
		fmt.Println("'cluster' flag is missing")
		os.Exit(1)
	}

	cmdArgs := namespaceCmdFlag.Args()

	if len(cmdArgs) < 1 {
		NamespaceCmdUsage()
		os.Exit(1)
	}

	argAlias := ""
	if len(cmdArgs) > 1 {
		argAlias = cmdArgs[1]
	}

	staroidClient := CreateClient()

	org, err := GetOrgFromName(staroidClient, *orgName)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	cluster, err := GetClusterFromName(staroidClient, org, *clusterName)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	switch cmdArgs[0] {
	case "create":
		if argAlias == "" {
			NamespaceCmdUsage()
			os.Exit(1)
		}
		commit, err := v1.NewCommitFromCommitLocation(*commitLoc)
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
		ns, err := staroidClient.V1().Namespace().
			WithOrg(org.Provider, org.Name).
			WithClusterID(cluster.ID).
			WithCommit(commit).
			Create(argAlias)
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
		header := []string{"ALIAS", "NAME", "TYPE", "PHASE"}
		rows := make([]*[]string, 0)
		rows = append(rows, &[]string{ns.Alias, ns.Namespace, ns.Type, ns.Phase})
		PrintTable(&header, &rows)
	case "list":
		namespaces, err := staroidClient.V1().Namespace().
			WithOrg(org.Provider, org.Name).
			WithClusterID(cluster.ID).
			GetAll()
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		header := []string{"ALIAS", "NAME", "TYPE", "PHASE"}
		rows := make([]*[]string, 0)
		for _, ns := range *namespaces {
			rows = append(rows, &[]string{ns.Alias, ns.Namespace, ns.Type, ns.Phase})
		}
		PrintTable(&header, &rows)
	default:
		NamespaceCmdUsage()
		os.Exit(1)
	}

}
