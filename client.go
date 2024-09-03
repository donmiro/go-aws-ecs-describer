package aws_ecs_describer

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

type ECSCluster struct {
	region string
	name   string
}

func Cluster(region, name string) (*ECSCluster, error) {
	if region == "" || name == "" {
		return nil, errors.New("you must specify region & cluster name")
	}

	return &ECSCluster{
		region: region,
		name:   name,
	}, nil
}

func (cluster *ECSCluster) GetClusterDescription() (map[string]interface{}, error) {
	svc, err := svcCreator(cluster.region)
	if err != nil {
		return nil, err
	}

	// Get all services from cluster:
	listServicesInput := &ecs.ListServicesInput{
		Cluster: aws.String(cluster.name),
	}

	services, err := svc.ListServices(context.TODO(), listServicesInput)
	if err != nil {
		return nil, fmt.Errorf("unable to list services: %w", err)
	}

	// Services was found
	outputMap := make(map[string]interface{})
	if len(services.ServiceArns) > 0 {
		describeServicesInput := &ecs.DescribeServicesInput{
			Cluster:  aws.String(cluster.name),
			Services: services.ServiceArns,
		}

		describedServices, err := svc.DescribeServices(context.TODO(), describeServicesInput)
		if err != nil {
			return nil, fmt.Errorf("unable to describe services: %w", err)
		}

		for _, service := range describedServices.Services {
			// Getting task description
			tasks, err := describeTasks(svc, cluster.name, service.ServiceName)
			if err != nil {
				log.Fatalf("Failed to describe tasks: %v", err)
				continue
			}
			outputMap[service.ServiceName] = tasks
		}

		return outputMap, nil
	}

	return outputMap, fmt.Errorf("Unexpected error")
}

func svcCreator(awsRegion string) (*ecs.Client, error) {
	// Create AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))
	if err != nil {
		return nil, err
	}

	svc := ecs.NewFromConfig(cfg)

	return svc, nil
}

func describeTasks(svc *ecs.Client, clusterName string, serviceName *string) ([]ecs.Task, error) {
	// Get all tasks for each service
	listTasksInput := &ecs.ListTasksInput{
		Cluster:     aws.String(clusterName),
		ServiceName: serviceName,
	}

	tasks, err := svc.ListTasks(context.TODO(), listTasksInput)
	if err != nil {
		return nil, fmt.Errorf("unable to list tasks: %w", err)
	}

	if len(tasks.TaskArns) > 0 {
		describeTasksInput := &ecs.DescribeTasksInput{
			Cluster: aws.String(clusterName),
			Tasks:   tasks.TaskArns,
		}

		describedTasks, err := svc.DescribeTasks(context.TODO(), describeTasksInput)
		if err != nil {
			return nil, fmt.Errorf("unable to describe tasks: %w", err)
		}

		return describedTasks.Tasks, nil
	}

	return []ecs.Task{}, nil
}
