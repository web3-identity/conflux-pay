package services

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/web3-identity/conflux-pay/models"
	"github.com/web3-identity/conflux-pay/utils"
)

// TODO: Monitor all orders already exists in db and auto close them when start programming.
func InitCloseOrderTask() {
	svr := NewOrderService()
	for {

		orders, err := models.FindNeedCloseOrders()

		if err == nil {
			var wg sync.WaitGroup
			for _, o := range orders {
				wg.Add(1)
				go func(order *models.Order) {
					defer func() {
						reason := recover()
						if reason != nil {
							logrus.WithField("order", order).WithField("reason", reason).Error("failed to close order")
						}
						// 无论是否error，重新Get刷新order
						svr.GetOrder(order.TradeNo)
						wg.Done()
					}()

					if err := utils.Retry(10, time.Second*5, func() error {
						_, err := svr.Close(order.TradeNo)
						return err
					}); err != nil {
						logrus.WithField("order", order).WithError(err).Error("failed to close order")
					}
				}(o)
			}
			wg.Wait()
			logrus.WithField("len", len(orders)).Info("close all old order completed")
			return
		}
		time.Sleep(10 * time.Second)
	}
}
