package main

import (
	"context"
	"fmt"
	"log"

	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
)

func main() {

	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Ошибка при выполнении запроса ListWorkflowExecutions:", err)
		return
	}
	defer c.Close()

	response, err := c.ListWorkflow(context.Background(), &workflowservice.ListWorkflowExecutionsRequest{
		PageSize: 10,
		Query:    "WorkflowType='PizzaOrderWorkflow' and ExecutionStatus='Running'",
	})
	if err != nil {
		log.Fatalln("Ошибка при выполнении запроса ListWorkflowExecutions:", err)
		return
	}

	for _, execution := range response.Executions {
		fmt.Println("WorkflowID:", execution.Execution.WorkflowId, "Status:", execution.Status)
	}

}
