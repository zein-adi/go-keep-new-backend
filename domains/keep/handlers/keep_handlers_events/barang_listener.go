package keep_handlers_events

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_events"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_service_interfaces"
	"time"
)

func NewBarangEventListenerHandler(posService keep_service_interfaces.IBarangServices) *BarangEventListenerHandler {
	return &BarangEventListenerHandler{
		service: posService,
	}
}

type BarangEventListenerHandler struct {
	service keep_service_interfaces.IBarangServices
}

func (x *BarangEventListenerHandler) TransaksiCreated(eventData any) {
	eventName, _, err := keep_events.NewTransaksiCreatedEventDataFromDispatcher(eventData)
	if err != nil {
		logrus.Error(err.Error())
		return
	}
	x.updateBarang(eventName)
}
func (x *BarangEventListenerHandler) TransaksiUpdated(eventData any) {
	eventName, _, err := keep_events.NewTransaksiUpdatedEventDataFromDispatcher(eventData)
	if err != nil {
		logrus.Error(err.Error())
		return
	}
	x.updateBarang(eventName)
}
func (x *BarangEventListenerHandler) TransaksiSoftDeleted(eventData any) {
	eventName, _, err := keep_events.NewTransaksiSoftDeleteEventDataFromDispatcher(eventData)
	if err != nil {
		logrus.Error(err.Error())
		return
	}
	x.updateBarang(eventName)
}
func (x *BarangEventListenerHandler) TransaksiRestored(eventData any) {
	eventName, _, err := keep_events.NewTransaksiRestoreEventDataFromDispatcher(eventData)
	if err != nil {
		logrus.Error(err.Error())
		return
	}
	x.updateBarang(eventName)
}

func (x *BarangEventListenerHandler) updateBarang(action string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	affected, err := x.service.UpdateBarangFromTransaksi(ctx)
	if err != nil {
		logrus.Error(err.Error())
		return
	}

	logrus.
		WithField("event", action).
		WithField("listener", "keep.barang").
		Infof("affected:%d", affected)
}
