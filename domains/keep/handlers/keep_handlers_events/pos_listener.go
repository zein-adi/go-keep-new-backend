package keep_handlers_events

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_events"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_service_interfaces"
	"time"
)

func NewPosEventListenerHandler(posService keep_service_interfaces.IPosServices) *PosEventListenerHandler {
	return &PosEventListenerHandler{
		service: posService,
	}
}

type PosEventListenerHandler struct {
	service keep_service_interfaces.IPosServices
}

/*
 * Transaksi
 */

func (x *PosEventListenerHandler) TransaksiCreated(eventData any) {
	data, err := keep_events.NewTransaksiCreatedEventDataFromDispatcher(eventData)
	if err != nil {
		logrus.Error(err.Error())
		return
	}
	x.updateSaldo(keep_events.TransaksiCreated, data.Data.PosAsalId, data.Data.PosTujuanId)
}
func (x *PosEventListenerHandler) TransaksiUpdated(eventData any) {
	data, err := keep_events.NewTransaksiUpdatedEventDataFromDispatcher(eventData)
	if err != nil {
		logrus.Error(err.Error())
		return
	}
	x.updateSaldo(keep_events.TransaksiUpdated, data.Old.PosAsalId, data.Old.PosTujuanId, data.New.PosAsalId, data.New.PosTujuanId)
}
func (x *PosEventListenerHandler) TransaksiSoftDeleted(eventData any) {
	data, err := keep_events.NewTransaksiSoftDeleteEventDataFromDispatcher(eventData)
	if err != nil {
		logrus.Error(err.Error())
		return
	}
	x.updateSaldo(keep_events.TransaksiSoftDeleted, data.Data.PosAsalId, data.Data.PosTujuanId)
}
func (x *PosEventListenerHandler) TransaksiRestored(eventData any) {
	data, err := keep_events.NewTransaksiRestoreEventDataFromDispatcher(eventData)
	if err != nil {
		logrus.Error(err.Error())
		return
	}
	x.updateSaldo(keep_events.TransaksiRestored, data.Data.PosAsalId, data.Data.PosTujuanId)
}
func (x *PosEventListenerHandler) updateSaldo(action string, ids ...string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	affected, err := x.service.UpdateSaldoFromTransaksi(ctx, ids)
	if err != nil {
		logrus.Error(err.Error())
		return
	}

	logrus.
		WithField("listener", "keep.pos").
		WithField("event", action).
		Infof("affected:%d", affected)
}
