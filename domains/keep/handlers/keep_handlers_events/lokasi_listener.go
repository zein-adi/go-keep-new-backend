package keep_handlers_events

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_events"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_service_interfaces"
	"time"
)

func NewLokasiEventListenerHandler(posService keep_service_interfaces.ILokasiServices) *LokasiEventListenerHandler {
	return &LokasiEventListenerHandler{
		service: posService,
	}
}

type LokasiEventListenerHandler struct {
	service keep_service_interfaces.ILokasiServices
}

func (x *LokasiEventListenerHandler) TransaksiCreated(eventData any) {
	_, err := keep_events.NewTransaksiCreatedEventDataFromDispatcher(eventData)
	if err != nil {
		logrus.Error(err.Error())
		return
	}
	x.updateLokasi("created")
}
func (x *LokasiEventListenerHandler) TransaksiUpdated(eventData any) {
	_, err := keep_events.NewTransaksiUpdatedEventDataFromDispatcher(eventData)
	if err != nil {
		logrus.Error(err.Error())
		return
	}
	x.updateLokasi("updated")
}
func (x *LokasiEventListenerHandler) TransaksiSoftDeleted(eventData any) {
	_, err := keep_events.NewTransaksiSoftDeleteEventDataFromDispatcher(eventData)
	if err != nil {
		logrus.Error(err.Error())
		return
	}
	x.updateLokasi("softDelete")
}
func (x *LokasiEventListenerHandler) TransaksiRestored(eventData any) {
	_, err := keep_events.NewTransaksiRestoreEventDataFromDispatcher(eventData)
	if err != nil {
		logrus.Error(err.Error())
		return
	}
	x.updateLokasi("restore")
}

func (x *LokasiEventListenerHandler) updateLokasi(action string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	affected, err := x.service.UpdateLokasiFromTransaksi(ctx)
	if err != nil {
		logrus.Error(err.Error())
		return
	}

	logrus.WithField("listener", "keep.lokasi."+action).Infof("affected:%d", affected)
}
