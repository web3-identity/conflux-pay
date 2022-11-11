package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/web3-identity/conflux-pay/models"
	"github.com/web3-identity/conflux-pay/models/enums"
)

func RunNotifyTask() {

	lastHandledId := uint(0)
	// 每秒循环
	// 根据time,count,is_completed决定是否发送通知，发送完跟新time,count,is_completed
	// for {
	// time.Sleep(time.Second * 3)
	orders, err := models.FindNeedNotifyOrders(lastHandledId)
	if err != nil {
		logrus.WithError(err).Error("failed find orders need to notify")
		// continue
		return
	}

	for _, o := range orders {
		// if !o.IsPayNotifyCompleted && o.TradeState.IsStable() {
		go runPayNotifyTask(o)
		// }
		// if !o.IsRefundNotifyCompleted && o.RefundState.IsStable(o.TradeState) {
		go runRefundNotifyTask(o)
		// }
		// lastHandledId = o.ID
	}
	// }
}

func runPayNotifyTask(o *models.Order) {
	if o.IsPayNotifyCompleted || !o.TradeState.IsStable() {
		return
	}

	defer func() {
		o.IsPayNotifyCompleted = true
		o.Save()
	}()

	if o.AppPayNotifyUrl == nil {
		return
	}

	if _, err := url.ParseRequestURI(*o.AppPayNotifyUrl); err != nil {
		return
	}

	for {
		notifyTime := calcNextNotifyTime(o.PayNotifyCount)
		<-time.After(time.Until(notifyTime))
		if err := sendNotify(*o.AppPayNotifyUrl, &o.OrderCore); err != nil {
			o.PayNotifyCount++
			o.Save()
			continue
		} else {
			return
		}
	}
}

func runRefundNotifyTask(o *models.Order) {
	if o.IsRefundNotifyCompleted || !o.RefundState.IsStable(o.TradeState) {
		return
	}

	defer func() {
		o.IsRefundNotifyCompleted = true
		o.Save()
	}()

	if o.TradeState != enums.TRADE_STATE_REFUND {
		return
	}

	if o.AppRefundNotifyUrl == nil {
		return
	}

	if _, err := url.ParseRequestURI(*o.AppRefundNotifyUrl); err != nil {
		return
	}

	for {
		notifyTime := calcNextNotifyTime(o.RefundNotifyCount)
		<-time.After(time.Until(notifyTime))
		if err := sendNotify(*o.AppRefundNotifyUrl, &o.OrderCore); err != nil {
			o.RefundNotifyCount++
			o.Save()
			continue
		} else {
			return
		}
	}
}

func sendNotify(url string, orderCore *models.OrderCore) error {
	fmt.Println("send pay notify")
	payBody, _ := json.Marshal(orderCore)
	resp, err := http.DefaultClient.Post(url, "application/json", bytes.NewBuffer(payBody))
	if err != nil {
		logrus.WithError(err).Error("failed to send pay notification")
		return err
	}

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	logrus.WithField("status", resp.Status).Error("failed to send pay notification")
	return fmt.Errorf("failed status: %v", resp.Status)
}

func calcNextNotifyTime(count int) time.Time {
	t := time.Now().Add(time.Second * time.Duration(count))
	fmt.Printf("now: %v, notify time: %v", time.Now(), t)
	return t
}