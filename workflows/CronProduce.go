package workflows

import (
	"fmt"
	"time"

	"github.com/quincy0/temproal-workflows/mq"
	"go.temporal.io/sdk/workflow"
)

func ProduceCron(ctx workflow.Context, parallelism int) (results []string, err error) {
	workflow.GetLogger(ctx).Info("Cron workflow started.", "StartTime", workflow.Now(ctx))

	cwo := workflow.ChildWorkflowOptions{
		WorkflowExecutionTimeout: 10 * time.Second,
	}
	ctx = workflow.WithChildOptions(ctx, cwo)

	for i := 0; i < parallelism; i++ {
		workflow.Go(ctx, func(ctx workflow.Context) {
			var result string
			err = workflow.ExecuteChildWorkflow(ctx, ProduceChildWorkflow).Get(ctx, &result)
			if err != nil {
				return
			}
			results = append(results, result)
		})
	}

	_ = workflow.Await(ctx, func() bool {
		return err != nil || len(results) == parallelism
	})

	return

}

func ProduceChildWorkflow(ctx workflow.Context) error {
	fmt.Println("produce data!")
	data := fmt.Sprintf("handle %s", time.Now().Format(time.StampNano))
	mq.Produce([]byte(data))
	return nil
}
