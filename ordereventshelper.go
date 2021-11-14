package orders

import (
	common "github.com/dev4fun007/autobot-common"
	"time"
)

func CreateOrderEvent(brokerName string, order common.Order, config common.BaseConfig, actionType common.ActionType, errMsg string) common.OrderEvent {
	totalAmountPaid := order.TotalQuantity*order.PricePerUnit + order.FeeAmount
	event := common.OrderEvent{
		Order:        order,
		EventError:   errMsg,
		BrokerName:   brokerName,
		StrategyName: config.Name,
		StrategyType: config.StrategyType,
		TotalAmount:  totalAmountPaid,
		Action:       actionType,
		Timestamp:    time.Now().Format(time.RFC3339),
	}
	return event
}
