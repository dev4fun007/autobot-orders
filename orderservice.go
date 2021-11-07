package orders

import (
	"context"
	"github.com/dev4fun007/autobot-common"
	"github.com/rs/zerolog/log"
)

const (
	OrderBatchSize  = 8
	OrderServiceTag = "OrderProcessorService"
)

type OrderProcessorService struct {
	limitChan    chan common.RequestLimitOrder
	marketChan   chan common.RequestMarketOrder
	repository   common.Repository
	brokerAction common.BrokerAction
}

func NewOrderProcessorService(action common.BrokerAction, repository common.Repository) *OrderProcessorService {
	return &OrderProcessorService{
		limitChan:    make(chan common.RequestLimitOrder, OrderBatchSize),
		marketChan:   make(chan common.RequestMarketOrder, OrderBatchSize),
		brokerAction: action,
		repository:   repository,
	}
}

func (processor OrderProcessorService) ExecuteMarketOrder(order common.RequestMarketOrder) {
	processor.marketChan <- order
}

func (processor OrderProcessorService) ExecuteLimitOrder(order common.RequestLimitOrder) {
	processor.limitChan <- order
}

func (processor OrderProcessorService) StartOrderService(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				{
					return
				}
			case mo := <-processor.marketChan:
				{
					orderResult, err := processor.brokerAction.ExecuteMarketOrder(mo)
					var event common.OrderEvent
					if err != nil {
						log.Error().Str(common.LogComponent, OrderServiceTag).
							Str("order-type", "Market").
							Err(err).
							Msg("error processing order")
						event = CreateOrderEvent(processor.brokerAction.GetBrokerName(), orderResult, mo.Config, err.Error())
					} else {
						log.Info().Str(common.LogComponent, OrderServiceTag).
							Str("order-type", "Market").
							Msg("order processed, creating event")
						event = CreateOrderEvent(processor.brokerAction.GetBrokerName(), orderResult, mo.Config, "")
					}
					_ = processor.repository.Save(ctx, event)
				}
			case lo := <-processor.limitChan:
				{
					orderResult, err := processor.brokerAction.ExecuteLimitOrder(lo)
					var event common.OrderEvent
					if err != nil {
						log.Error().Str(common.LogComponent, OrderServiceTag).
							Str("order-type", "Limit").
							Err(err).
							Msg("error processing order")
						event = CreateOrderEvent(processor.brokerAction.GetBrokerName(), orderResult, lo.Config, err.Error())
					} else {
						log.Info().Str(common.LogComponent, OrderServiceTag).
							Str("order-type", "Limit").
							Msg("order processed, creating event")
						event = CreateOrderEvent(processor.brokerAction.GetBrokerName(), orderResult, lo.Config, "")
					}
					_ = processor.repository.Save(ctx, event)
				}
			}
		}
	}()
}
