package orders

import (
	common "github.com/dev4fun007/autobot-common"
	"time"
)

func CreateOrderEvent(brokerName string, order common.Order, config common.BaseConfig, errMsg string) common.OrderEvent {
	event := common.OrderEvent{
		Order:        order,
		EventError:   errMsg,
		BrokerName:   brokerName,
		StrategyName: config.Name,
		StrategyType: config.StrategyType,
		Timestamp:    time.Now().Format(time.RFC3339),
	}
	return event
}
