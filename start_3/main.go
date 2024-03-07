package main

import (
	"context"
	"log"
	"pizzaorder"

	"go.temporal.io/sdk/client"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Ошибка создания клиента Temporal:", err)
	}
	defer c.Close()

	workflowID := "pizza-workflow"

	ctx := context.Background()

	resp, err := c.QueryWorkflow(ctx, workflowID, "", pizzaorder.QueryType)
	if err != nil {
		log.Fatalln("Ошибка запроса статуса заказа:", err)
	}

	if err != nil {
		log.Fatalln("Unable to query workflow", err)
	}
	var result interface{}
	if err := resp.Get(&result); err != nil {
		log.Fatalln("Unable to decode query result", err)
	}

	log.Printf("Текущий статус заказа: ", result)
}
