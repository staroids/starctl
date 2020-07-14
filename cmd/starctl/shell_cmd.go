package main

import (
	"flag"
	"fmt"
	"os"
)

func ShellCmdUsage() {
	fmt.Fprintf(os.Stdout, "shell [flags] [start|stop] <namespace alias>\n")
}

func ShellCmd(args []string) {
	shellCmdFlag := flag.NewFlagSet("shell", flag.ExitOnError)
	orgName := shellCmdFlag.String("org", "", "organization (e.g. GITHUB/staroid)")
	clusterName := shellCmdFlag.String("cluster", "", "name of cluster")

	shellCmdFlag.Parse(args)

	if *orgName == "" {
		fmt.Println("'org' flag is missing")
		os.Exit(1)
	}

	if *clusterName == "" {
		fmt.Println("'cluster' flag is missing")
		os.Exit(1)
	}

	cmdArgs := shellCmdFlag.Args()

	if len(cmdArgs) < 1 {
		ShellCmdUsage()
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
	case "start":
		if argAlias == "" {
			ShellCmdUsage()
			os.Exit(1)
		}
		ns, err := staroidClient.V1().Namespace().
			WithOrg(org.Provider, org.Name).
			WithClusterID(cluster.ID).
			Get(argAlias)
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		if ns.Phase != "RUNNING" {
			fmt.Printf("Namespace %v is not running", argAlias)
			os.Exit(1)
		}

		err = staroidClient.V1().Namespace().
			WithOrg(org.Provider, org.Name).
			WithClusterID(cluster.ID).
			ShellStartById(ns.ID)
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
	case "stop":
		if argAlias == "" {
			ShellCmdUsage()
			os.Exit(1)
		}
		ns, err := staroidClient.V1().Namespace().
			WithOrg(org.Provider, org.Name).
			WithClusterID(cluster.ID).
			Get(argAlias)
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		if ns.Phase != "RUNNING" {
			fmt.Printf("Namespace %v is not running", argAlias)
			os.Exit(1)
		}

		err = staroidClient.V1().Namespace().
			WithOrg(org.Provider, org.Name).
			WithClusterID(cluster.ID).
			ShellStopById(ns.ID)
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
	default:
		ShellCmdUsage()
		os.Exit(1)
	}
}
