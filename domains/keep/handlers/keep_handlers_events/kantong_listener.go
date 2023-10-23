package keep_handlers_events

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_events"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_service_interfaces"
	"time"
)

func NewKantongEventListenerHandler(kantongServices keep_service_interfaces.IKantongServices) *KantongEventListenerHandler {
	return &KantongEventListenerHandler{
		service: kantongServices,
	}
}

type KantongEventListenerHandler struct {
	service keep_service_interfaces.IKantongServices
}

func (x *KantongEventListenerHandler) TransaksiCreated(eventData any) {
	data, err := keep_events.NewTransaksiCreatedEventDataFromDispatcher(eventData)
	if err != nil {
		logrus.Error(err.Error())
		return
	}
	x.updateSaldo("created", data.Data.KantongAsalId, data.Data.KantongTujuanId, data.Data.Jumlah,
		"", "", 0)
}
func (x *KantongEventListenerHandler) TransaksiUpdated(eventData any) {
	data, err := keep_events.NewTransaksiUpdatedEventDataFromDispatcher(eventData)
	if err != nil {
		logrus.Error(err.Error())
		return
	}
	x.updateSaldo("updated", data.New.KantongAsalId, data.New.KantongTujuanId, data.New.Jumlah,
		data.Old.KantongAsalId, data.Old.KantongTujuanId, data.Old.Jumlah)
}
func (x *KantongEventListenerHandler) TransaksiSoftDelete(eventData any) {
	data, err := keep_events.NewTransaksiSoftDeleteEventDataFromDispatcher(eventData)
	if err != nil {
		logrus.Error(err.Error())
		return
	}
	x.updateSaldo("softDelete", "", "", 0,
		data.Data.KantongAsalId, data.Data.KantongTujuanId, data.Data.Jumlah)
}
func (x *KantongEventListenerHandler) TransaksiRestore(eventData any) {
	data, err := keep_events.NewTransaksiRestoreEventDataFromDispatcher(eventData)
	if err != nil {
		logrus.Error(err.Error())
		return
	}
	x.updateSaldo("restore", data.Data.KantongAsalId, data.Data.KantongTujuanId, data.Data.Jumlah,
		"", "", 0)
}

func (x *KantongEventListenerHandler) updateSaldo(action, asalId, tujuanId string, jumlah int, oldAsalId, oldTujuanId string, oldJumlah int) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	affected, err := x.service.UpdateSaldo(ctx, asalId, tujuanId, jumlah, oldAsalId, oldTujuanId, oldJumlah)
	if err != nil {
		logrus.Error(err.Error())
		return
	}

	logrus.WithField("listener", "keep.kantong."+action).Infof("affected:%d", affected)
}
