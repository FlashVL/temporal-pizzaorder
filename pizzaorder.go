package pizzaorder

import (
	"context"
	"log"
	"time"

	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/workflow"
)

const (
	TaskQueueName       = "pizzaOrderQueue"
	UpdateAddressSignal = "updateAddress"
	QueryType           = "orderStatus"
)

type Order struct {
	Id      string
	Address string
}

type OutDate struct {
	Status       string
	DeliveryTime int
	Address      string
}

type routeSignal struct {
	address string
}
type OrderStatus struct {
	Status       string
	DeliveryTime time.Time
}

type doughPreparationDetails struct {
	OrderID string
}

type doughPreparationResult struct {
	Success bool
	Message string
}

func PizzaOrderWorkflow(ctx workflow.Context, order Order) (OutDate, error) {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	currentStatus := OrderStatus{Status: "Принят", DeliveryTime: time.Time{}}

	workflow.SetQueryHandler(ctx, QueryType, func() (OrderStatus, error) {
		return currentStatus, nil
	})

	updateAddressChan := workflow.GetSignalChannel(ctx, UpdateAddressSignal)

	selector := workflow.NewSelector(ctx)

	var address string

	selector.AddReceive(updateAddressChan, func(c workflow.ReceiveChannel, _ bool) {
		c.Receive(ctx, &address)
	})

	currentStatus.Status = "Готовим тесто"

	childWorkflowOptions := workflow.ChildWorkflowOptions{
		ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
	}
	ctx = workflow.WithChildOptions(ctx, childWorkflowOptions)

	doughDetails := doughPreparationDetails{
		OrderID: order.Id,
	}

	var doughResult doughPreparationResult
	err := workflow.ExecuteChildWorkflow(ctx, PreparePizzaDoughWorkflow, doughDetails).Get(ctx, &doughResult)
	if err != nil {
		log.Println("Ошибка при подготовке теста:", err)
		return OutDate{}, err
	}

	currentStatus.Status = "Выпикаем пиццу"

	var status string
	err = workflow.ExecuteActivity(ctx, PreparePizzaActivity, order).Get(ctx, &status)

	if err != nil {
		log.Println("Ошибка при выполнении действия PreparePizza:", err)
		return OutDate{}, err
	}

	currentStatus.Status = "Готово"

	selector.Select(ctx)

	var deliveryTime int
	err = workflow.ExecuteActivity(ctx, DeliveryPizzaActivity, order).Get(ctx, &deliveryTime)

	if err != nil {
		log.Println("Ошибка при выполнении действия DeliveryPizzaActivity:", err)
		return OutDate{}, err
	}

	result := OutDate{
		Status:       status,
		DeliveryTime: deliveryTime,
		Address:      address,
	}

	return result, nil
}

func PreparePizzaActivity(ctx context.Context, order Order) (string, error) {

	time.Sleep(5 * time.Second)
	return "готово", nil
}

func DeliveryPizzaActivity(ctx context.Context, order Order) (int, error) {

	//time.Sleep(5 * time.Second)
	return 60, nil
}

func PreparePizzaDoughWorkflow(ctx workflow.Context, details doughPreparationDetails) (doughPreparationResult, error) {

	workflow.Sleep(ctx, 5*time.Second)

	return doughPreparationResult{
		Success: true,
		Message: "Тесто подготовлено",
	}, nil
}
