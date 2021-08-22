package zeebe

import (
	"context"
	"fmt"
	"time"

	"git.coryptex.com/lib/sago/sagolog"
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
		sagolog.Log(sagolog.INFO,
			fmt.Sprintf("ZB Broker %s:%d", broker.Host, broker.Port))
		for _, partition := range broker.Partitions {
			sagolog.Log(sagolog.INFO,
				fmt.Sprintf("ZB Partition %d:%s", partition.PartitionId, roleToString(partition.Role)))
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
	_, err := client.NewDeployProcessCommand().AddResourceFile(path).Send(ctx)
	if err != nil {
		return err
	}
	return nil
}

func CreateProcessInstance(client zbc.Client, pid string, variables map[string]interface{}) error {
	request, err := client.NewCreateInstanceCommand().
		BPMNProcessId(pid).
		LatestVersion().
		VariablesFromMap(variables)

	if err != nil {
		return err
	}

	_, err = request.Send(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func FailJob(client worker.JobClient, job entities.Job, msg string) error {
	_, err := client.NewFailJobCommand().JobKey(job.GetKey()).Retries(job.Retries - 1).
		ErrorMessage(msg).
		Send(context.Background())
	if err != nil {
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
