package zeebe

import (
	"context"
	"log"
	"time"

	"github.com/camunda-cloud/zeebe/clients/go/pkg/entities"
	"github.com/camunda-cloud/zeebe/clients/go/pkg/pb"
	"github.com/camunda-cloud/zeebe/clients/go/pkg/worker"
	"github.com/camunda-cloud/zeebe/clients/go/pkg/zbc"
	"github.com/pkg/errors"
)

func GetZeebeClient(address string) (zbc.Client, error) {
	client, err := zbc.NewClient(&zbc.ClientConfig{
		GatewayAddress:         address,
		UsePlaintextConnection: true,
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}

func ListBrokers(client zbc.Client) {
	ctx := context.Background()
	topology, err := client.NewTopologyCommand().Send(ctx)
	panicOnErr(err)

	for _, broker := range topology.Brokers {
		log.Println("Broker", broker.Host, ":", broker.Port)
		for _, partition := range broker.Partitions {
			log.Println("Partition", partition.PartitionId, ":", roleToString(partition.Role))
		}
	}
}

func StartJobWorker(client zbc.Client, jobType string, handler worker.JobHandler) (close func()) {
	jobWorker := client.NewJobWorker().JobType(jobType).Handler(handler).Open()

	return func() {
		jobWorker.Close()
		jobWorker.AwaitClose()
	}
}

func DeployProcess(client zbc.Client, path string) error {
	// After the client is created
	ctx := context.Background()
	response, err := client.NewDeployProcessCommand().AddResourceFile(path).Send(ctx)
	if err != nil {
		return err
	}

	log.Println(response.String())
	return nil
}

func CreateProcessInstance(client zbc.Client, pid string, variables map[string]interface{}) {
	request, err :=
		client.NewCreateInstanceCommand().
			BPMNProcessId(pid).
			LatestVersion().
			VariablesFromMap(variables)

	panicOnErr(err)

	ctx := context.Background()

	msg, err := request.Send(ctx)
	panicOnErr(err)

	log.Println(msg.String())
}

func FailJob(client worker.JobClient, job entities.Job, msg string) error {
	log.Println(msg)
	_, err := client.NewFailJobCommand().JobKey(job.GetKey()).Retries(job.Retries - 1).
		ErrorMessage(msg).
		Send(context.Background())
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func PublishMessage(ctx context.Context, client zbc.Client, msgName, correlationId string, vars map[string]interface{}) error {
	req, err := client.NewPublishMessageCommand().
		MessageName(msgName).
		CorrelationKey(correlationId).
		TimeToLive(time.Minute).
		VariablesFromMap(vars)

	if err != nil {
		return errors.Wrapf(err,
			"failed to create zb %s:%s PublishMessageCommand\nerr: %v",
			msgName, correlationId, err)
	}
	_, err = req.Send(ctx)
	if err != nil {
		return errors.Wrapf(err,
			"failed to send zb %s:%s PublishMessageCommand\nerr: %v",
			msgName, correlationId, err)
	}
	log.Printf("zb message %s:%s published\nvariables: %v\n", msgName, correlationId, vars)
	return nil
}

func roleToString(role pb.Partition_PartitionBrokerRole) string {
	switch role {
	case pb.Partition_LEADER:
		return "Leader"
	case pb.Partition_FOLLOWER:
		return "Follower"
	default:
		return "Unknown"
	}
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
