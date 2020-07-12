package main

import (
	"fmt"

	"github.com/staroids/starctl/pkg/api"
	v1 "github.com/staroids/starctl/pkg/api/v1"
)

func GetOrgFromName(client *api.StaroidClient, orgName string) (*v1.StaroidOrg, error) {
	orgs, err := client.V1().Org().GetAll()
	if err != nil {
		return nil, err
	}

	var org *v1.StaroidOrg = nil
	for _, o := range *orgs {
		if orgName == fmt.Sprintf("%s/%s", o.Provider, o.Name) {
			org = &o
			break
		}
	}

	if org == nil {
		err = fmt.Errorf("Org '%s' not found\n", orgName)
		return nil, err
	}

	return org, nil
}

func GetClusterFromName(client *api.StaroidClient, org *v1.StaroidOrg, clusterName string) (*v1.StaroidCluster, error) {
	clusters, err := client.V1().Cluster().WithOrg(org.Provider, org.Name).GetAll()
	if err != nil {
		return nil, err
	}
	var cluster *v1.StaroidCluster = nil
	for _, c := range *clusters {
		if c.Name == clusterName {
			cluster = &c
			break
		}
	}
	if cluster == nil {
		err = fmt.Errorf("Cluster '%s' not found\n", clusterName)
		return nil, err
	}

	return cluster, nil
}

func GetNamespaceFromAlias(client *api.StaroidClient, org *v1.StaroidOrg, cluster *v1.StaroidCluster, nsAlias string) (*v1.StaroidNamespace, error) {
	namespaces, err := client.V1().Namespace().WithOrg(org.Provider, org.Name).WithClusterID(cluster.ID).GetAll()
	if err != nil {
		return nil, err
	}

	var ns *v1.StaroidNamespace = nil
	for _, n := range *namespaces {
		if n.Alias == nsAlias {
			ns = &n
			break
		}
	}
	if ns == nil {
		err = fmt.Errorf("Namespace alias '%s' not found\n", nsAlias)
		return nil, err
	}

	return ns, nil
}
