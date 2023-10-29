package keep_handlers_events

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_events"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_request"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_service_interfaces"
	"time"
)

func NewKantongEventListenerHandler(kantongServices keep_service_interfaces.IKantongServices, kantongHistoryServices keep_service_interfaces.IKantongHistoryServices) *KantongEventListenerHandler {
	return &KantongEventListenerHandler{
		kantongServices:        kantongServices,
		kantongHistoryServices: kantongHistoryServices,
	}
}

type KantongEventListenerHandler struct {
	kantongServices        keep_service_interfaces.IKantongServices
	kantongHistoryServices keep_service_interfaces.IKantongHistoryServices
}

/*
 * Kantong History
 */

func (x *KantongEventListenerHandler) KantongHistoryCreated(eventData any) {
	eventName, data, err := keep_events.NewKantongHistoryCreatedEventDataFromDispatcher(eventData)
	if err != nil {
		logrus.Error(err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	asalId := ""
	tujuanId := ""
	jumlah := 0

	if data.Data.Jumlah >= 0 {
		tujuanId = data.Data.KantongId
		jumlah = data.Data.Jumlah
	} else {
		asalId = data.Data.KantongId
		jumlah = -1 * data.Data.Jumlah
	}

	affected, err := x.kantongServices.UpdateSaldo(ctx, asalId, tujuanId, jumlah, "", "", 0)
	if err != nil {
		logrus.Error(err.Error())
		return
	}

	logrus.
		WithField("event", eventName).
		WithField("listener", "keep.kantong").
		Infof("affected:%d", affected)
}

/*
 * Transaksi
 */

func (x *KantongEventListenerHandler) TransaksiCreated(eventData any) {
	eventName, data, err := keep_events.NewTransaksiCreatedEventDataFromDispatcher(eventData)
	if err != nil {
		logrus.Error(err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	affected := 0
	if data.Data.KantongAsalId != "" {
		_, err = x.kantongHistoryServices.Insert(ctx, &keep_request.KantongHistoryInsertUpdate{
			KantongId: data.Data.KantongAsalId,
			Jumlah:    data.Data.Jumlah * -1,
			Uraian:    data.Data.Uraian,
		})
		if err != nil {
			logrus.Error(err.Error())
			return
		}
		affected++
	}
	if data.Data.KantongTujuanId != "" {
		_, err = x.kantongHistoryServices.Insert(ctx, &keep_request.KantongHistoryInsertUpdate{
			KantongId: data.Data.KantongTujuanId,
			Jumlah:    data.Data.Jumlah,
			Uraian:    data.Data.Uraian,
		})
		if err != nil {
			logrus.Error(err.Error())
			return
		}
		affected++
	}

	logrus.
		WithField("event", eventName).
		WithField("listener", "keep.kantong").
		Infof("affected:%d", affected)
}
