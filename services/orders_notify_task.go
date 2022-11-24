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

var (
	//15s/15s/30s/3m/10m/20m/30m/30m/30m/60m/3h/3h/3h/6h/6h
	notifyIntervals = []time.Duration{
		time.Second * 0,
		time.Second * 15,
		time.Second * 15,
		time.Second * 30,
		time.Minute * 3,
		time.Minute * 10,
		time.Minute * 20,
		time.Minute * 30,
		time.Minute * 30,
		time.Minute * 60,
		time.Hour * 3,
		time.Hour * 3,
		time.Hour * 3,
		time.Hour * 6,
		time.Hour * 6,
	}
)

func RunNotifyTask() {

	for {
		// lastHandledId := uint(0)
		// 每秒循环
		// 根据time,count,is_completed决定是否发送通知，发送完跟新time,count,is_completed
		// for {
		// time.Sleep(time.Second * 3)
		orders, err := models.FindNeedNotifyOrders(0)
		if err != nil {
			logrus.WithError(err).Error("failed find orders need to notify")
			time.Sleep(10 * time.Second)
			continue
		}

		for _, o := range orders {
			go runPayNotifyTask(o)
			go runRefundNotifyTask(o)
		}
		return
	}

}

func runPayNotifyTask(o *models.Order) {
	fmt.Println("run pay notify task")
	if o.IsPayNotifyCompleted || (!o.TradeState.IsStable() && !o.TradeState.IsSuccess()) {
		return
	}

	// fmt.Println("aaa")
	defer func() {
		o.IsPayNotifyCompleted = true
		o.Save()
	}()

	// fmt.Println("bbb")
	if o.AppPayNotifyUrl == nil {
		return
	}

	// fmt.Println("ccc")
	if _, err := url.ParseRequestURI(*o.AppPayNotifyUrl); err != nil {
		return
	}

	// fmt.Println("ddd")
	for {
		notifyTime := calcNextNotifyTime(o.PayNotifyCount)
		fmt.Printf("pay notify time:%v \n", notifyTime)
		if notifyTime == nil {
			o.IsPayNotifyCompleted = true
			o.Save()
			return
		}

		<-time.After(time.Until(*notifyTime))
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
		if notifyTime == nil {
			o.IsRefundNotifyCompleted = true
			o.Save()
			return
		}
		<-time.After(time.Until(*notifyTime))
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
	fmt.Println("send notify")
	payBody, _ := json.Marshal(orderCore)
	resp, err := http.DefaultClient.Post(url, "application/json", bytes.NewBuffer(payBody))
	if err != nil {
		logrus.WithError(err).Error("failed to send pay notification")
		return err
	}

	if resp.StatusCode == http.StatusOK {
		logrus.WithField("url", url).Info("success send notidy")
		return nil
	}

	logrus.WithField("status", resp.Status).Error("failed to send pay notification")
	return fmt.Errorf("failed status: %v", resp.Status)
}

func calcNextNotifyTime(count int) *time.Time {
	if len(notifyIntervals) <= count {
		return nil
	}
	// t := time.Now().Add(time.Second * time.Duration(count))
	t := time.Now().Add(notifyIntervals[count])
	fmt.Printf("now: %v, notify time: %v", time.Now(), t)
	return &t
}
