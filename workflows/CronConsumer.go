package workflows

import (
	"fmt"
	"time"

	"github.com/quincy0/temproal-workflows/mq"
	"go.temporal.io/sdk/workflow"
)

func ConsumeCron(ctx workflow.Context, parallelism int) (results []string, err error) {
	workflow.GetLogger(ctx).Info("ConsumeCron workflow started.", "StartTime", workflow.Now(ctx))

	cwo := workflow.ChildWorkflowOptions{
		WorkflowExecutionTimeout: 10 * time.Second,
	}
	ctx = workflow.WithChildOptions(ctx, cwo)

	for i := 0; i < parallelism; i++ {
		workflow.Go(ctx, func(ctx workflow.Context) {
			var result string
			err = workflow.ExecuteChildWorkflow(ctx, ConsumeChildWorkflow).Get(ctx, &result)
			if err != nil {
				return
			}
			fmt.Println(result)
			results = append(results, result)
		})
	}

	_ = workflow.Await(ctx, func() bool {
		return err != nil || len(results) == parallelism
	})

	return
}

func ConsumeChildWorkflow(ctx workflow.Context) (string, error) {
	result := ""
	for {
		select {
		case <-time.After(5 * time.Second):
			if len(result) == 0 {
				return "empty msg", nil
			}
			return result, nil
		case msg := <-mq.MsgChan:
			result += "\n" + msg
			return msg, nil
		}
	}

}
