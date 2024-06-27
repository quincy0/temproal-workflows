package worker

import (
	"log"

	"github.com/quincy0/temproal-workflows/workflows"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func Init() {
	// The client and worker are heavyweight objects that should be created once per process.
	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "cron", worker.Options{})

	w.RegisterWorkflow(workflows.ProduceCron)
	w.RegisterWorkflow(workflows.ConsumeCron)

	w.RegisterWorkflow(workflows.ProduceChildWorkflow)
	w.RegisterWorkflow(workflows.ConsumeChildWorkflow)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
