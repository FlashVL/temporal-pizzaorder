package main

import (
	"context"
	"fmt"
	"log"

	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
)

func main() {
	// Создание клиента Temporal
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Ошибка создания клиента Temporal:", err)
		return
	}
	defer c.Close()

	var pageSize int32
	pageSize = 1

	var nextPageToken []byte

	i := 1

	for {
		response, err := c.ListWorkflow(context.Background(), &workflowservice.ListWorkflowExecutionsRequest{
			Namespace:     "default",
			PageSize:      pageSize,
			Query:         "WorkflowType='PizzaOrderWorkflow' and ExecutionStatus='Running'",
			NextPageToken: nextPageToken,
		})
		if err != nil {
			log.Fatalln("Ошибка при выполнении запроса ListWorkflowExecutions:", err)
			return
		}
		for _, execution := range response.Executions {
			fmt.Println("WorkflowID:", execution.Execution.WorkflowId, "RunID:", execution.Execution.RunId)
		}

		fmt.Println("Итерация:", i)
		i++
		if len(response.NextPageToken) == 0 {
			break
		}

		nextPageToken = response.NextPageToken
	}

	fmt.Println("Все страницы были успешно извлечены.")
}
