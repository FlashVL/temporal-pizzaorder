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
		log.Fatalln("Не удалось создать клиента Temporal:", err)
	}
	defer c.Close()

	workflowID := "pizza-workflow"
	newAddress := "Новый адрес, новая улица"

	err = c.SignalWorkflow(context.Background(), workflowID, "", pizzaorder.UpdateAddressSignal, newAddress)
	if err != nil {
		log.Fatalln("Ошибка при отправке сигнала обновления адреса:", err)
	}

}
