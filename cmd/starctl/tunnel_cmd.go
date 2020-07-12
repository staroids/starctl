package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	chclient "github.com/jpillora/chisel/client"
	"github.com/staroids/starctl/pkg/constants"
	corev1 "k8s.io/api/core/v1"
)

func TunnelCmdUsage(flagSet *flag.FlagSet) {
	fmt.Fprintf(os.Stdout, "tunnel [flags] [remote] ([remote], [remote], ...)\n\n")
	flagSet.Usage()
}

func TunnelCmd(args []string) {
	tunnelCmdFlag := flag.NewFlagSet("tunnel", flag.ExitOnError)
	orgName := tunnelCmdFlag.String("org", "", "organization (e.g. GITHUB/staroid)")
	clusterName := tunnelCmdFlag.String("cluster", "", "name of cluster")
	nsAlias := tunnelCmdFlag.String("ns-alias", "", "namespace alias")
	kubeProxy := tunnelCmdFlag.Bool("kube-proxy", false, "Kubernetes API proxy")
	kubeProxyPort := tunnelCmdFlag.Int("kube-proxy-port", 8001, "Local port for Kubernetes API proxy")

	if len(args) < 1 {
		TunnelCmdUsage(tunnelCmdFlag)
		os.Exit(1)
	}

	// check required flags
	tunnelCmdFlag.Parse(args)

	if *orgName == "" {
		fmt.Println("'org' flag is missing")
		os.Exit(1)
	}

	if *clusterName == "" {
		fmt.Println("'cluster' flag is missing")
		os.Exit(1)
	}

	if *nsAlias == "" {
		fmt.Println("'ns-alias' flag is missing")
		os.Exit(1)
	}

	remotes := tunnelCmdFlag.Args()
	if *kubeProxy {
		remotes = append(remotes, fmt.Sprintf("%d:localhost:%d", *kubeProxyPort, 57683))
	}

	if len(remotes) == 0 {
		fmt.Println("Set at least one [remote] argument or set '-kube-proxy' flag")
		os.Exit(1)
	}

	// valid value
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

	namespace, err := GetNamespaceFromAlias(staroidClient, org, cluster, *nsAlias)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	resources, err := staroidClient.V1().Namespace().WithName(namespace.Namespace).GetAllResources()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	var shellService *corev1.Service = nil
	for _, service := range resources.Services.Items {
		if service.ObjectMeta.Labels[constants.K8S_LABEL_KEY_RESOURCE_SYSTEM] == constants.K8S_LABEL_VALUE_RESOURCE_SYSTEM_SHELL {
			shellService = &service
			break
		}
	}

	if shellService == nil {
		fmt.Printf("Shell service is not found\n", err)
		os.Exit(1)
	}

	tunnelServerURL := namespace.ServiceURL(shellService.GetName(), constants.TunnelServicePort)

	chConfig := chclient.Config{
		Server:           tunnelServerURL,
		KeepAlive:        0,
		MaxRetryCount:    -1,
		MaxRetryInterval: 0,
		Headers:          http.Header{},
		Remotes:          remotes,
	}
	chConfig.Headers.Set("Authorization", fmt.Sprintf("token %s", staroidClient.Auth.AccessToken()))
	chClient, err := chclient.NewClient(&chConfig)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	if *kubeProxy {
		fmt.Printf("--------------------\n")
		fmt.Printf("Kubernetes API proxy localhost:%d configured\n\n", *kubeProxyPort)
		fmt.Printf("Try 'kubectl --server localhost:%d -n %s <kubectl command>'\n", *kubeProxyPort, namespace.Namespace)
		fmt.Printf("--------------------\n")
	}

	err = chClient.Run()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}