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
	_, err := keep_events.NewTransaksiCreatedEventDataFromDispatcher(eventData)
	if err != nil {
		logrus.Error(err.Error())
		return
	}
	x.updateBarang("created")
}
func (x *BarangEventListenerHandler) TransaksiUpdated(eventData any) {
	_, err := keep_events.NewTransaksiUpdatedEventDataFromDispatcher(eventData)
	if err != nil {
		logrus.Error(err.Error())
		return
	}
	x.updateBarang("updated")
}
func (x *BarangEventListenerHandler) TransaksiSoftDelete(eventData any) {
	_, err := keep_events.NewTransaksiSoftDeleteEventDataFromDispatcher(eventData)
	if err != nil {
		logrus.Error(err.Error())
		return
	}
	x.updateBarang("softDelete")
}
func (x *BarangEventListenerHandler) TransaksiRestore(eventData any) {
	_, err := keep_events.NewTransaksiRestoreEventDataFromDispatcher(eventData)
	if err != nil {
		logrus.Error(err.Error())
		return
	}
	x.updateBarang("restore")
}

func (x *BarangEventListenerHandler) updateBarang(action string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	affected, err := x.service.UpdateBarangFromTransaksi(ctx)
	if err != nil {
		logrus.Error(err.Error())
		return
	}

	logrus.WithField("listener", "keep.barang."+action).Infof("affected:%d", affected)
}
