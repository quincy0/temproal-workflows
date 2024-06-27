package starter

import (
	"context"
	"log"
	"time"

	"github.com/quincy0/temproal-workflows/workflows"
	"github.com/redis/go-redis/v9"
	"go.temporal.io/sdk/client"
)

const (
	WorkflowCronProducerID = "cron_producer"
	WorkflowCronConsumerID = "cron_consumer"
)

func getLock() bool {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	lockRes, err := rdb.SetNX(context.Background(), "starterLock", 1, 10*time.Minute).Result()
	if !lockRes || err != nil {
		return false
	}
	return true
}

func InitCron() {
	if !getLock() {
		return
	}
	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		return
	}
	defer c.Close()

	workflowProducerOptions := client.StartWorkflowOptions{
		ID:           WorkflowCronProducerID,
		TaskQueue:    "cron",
		CronSchedule: "* * * * *",
	}
	runProducer, err := c.ExecuteWorkflow(context.Background(), workflowProducerOptions, workflows.ProduceCron, 200)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}

	workflowConsumerOptions := client.StartWorkflowOptions{
		ID:           WorkflowCronConsumerID,
		TaskQueue:    "cron",
		CronSchedule: "*/2 * * * *",
	}
	runConsumer, err := c.ExecuteWorkflow(context.Background(), workflowConsumerOptions, workflows.ConsumeCron, 200)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}

	log.Println("Started workflow", "CronProducer WorkflowID", runProducer.GetID(), "RunID", runProducer.GetRunID())
	log.Println("Started workflow", "CronConsumer WorkflowID", runConsumer.GetID(), "RunID", runConsumer.GetRunID())
}
