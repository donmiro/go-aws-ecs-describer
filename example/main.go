package main

import (
	"fmt"

	ecsWrapper "github.com/donmiro/go-aws-ecs-describer"
)

func main() {
	region := "YOUR_REGION"
	name := "YOUR_CLUSTER_NAME"

	cluster, err := ecsWrapper.Cluster(region, name)
	if err != nil {
		fmt.Println(err)
	}

	description, err := cluster.GetClusterDescription()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(description)
}
