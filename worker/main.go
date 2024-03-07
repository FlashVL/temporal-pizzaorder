package main

import (
	"log"

	"pizzaorder"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()
	w := worker.New(c, pizzaorder.TaskQueueName, worker.Options{})

	w.RegisterWorkflow(pizzaorder.PizzaOrderWorkflow)
	w.RegisterWorkflow(pizzaorder.PreparePizzaDoughWorkflow)
	w.RegisterActivity(pizzaorder.PreparePizzaActivity)
	w.RegisterActivity(pizzaorder.DeliveryPizzaActivity)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
