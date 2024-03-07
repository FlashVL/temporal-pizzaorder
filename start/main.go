package main

import (
	"context"
	"encoding/json"
	"log"

	"pizzaorder"

	"go.temporal.io/sdk/client"
)

func main() {

	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	options := client.StartWorkflowOptions{
		ID:        "pizza-workflow",
		TaskQueue: pizzaorder.TaskQueueName,
	}

	input := pizzaorder.Order{
		Id: "O-547",
	}

	we, err := c.ExecuteWorkflow(context.Background(), options, pizzaorder.PizzaOrderWorkflow, input)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}
	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

	var result string
	err = we.Get(context.Background(), &result)
	if err != nil {
		log.Fatalln("Unable get workflow result", err)
	}

	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalln("Unable to format result in JSON format", err)
	}
	log.Printf("Workflow result: %s\n", string(data))
}
