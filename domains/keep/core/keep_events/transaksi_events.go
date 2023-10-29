package keep_events

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_entities"
	"time"
)

const (
	TransaksiCreated     = "keep.transaksi.created"
	TransaksiUpdated     = "keep.transaksi.updated"
	TransaksiSoftDeleted = "keep.transaksi.softDeleted"
	TransaksiRestored    = "keep.transaksi.restored"
	TransaksiHardDeleted = "keep.transaksi.hardDeleted"
)

type TransaksiEventData struct {
	Time            time.Time
	Id              string
	PosAsalId       string
	PosTujuanId     string
	KantongAsalId   string
	KantongTujuanId string
	Jumlah          int
	Lokasi          string
	Uraian          string
	Details         []*keep_entities.TransaksiDetail
}
type TransaksiCreatedEventData struct {
	Time time.Time
	Data TransaksiEventData
}
type TransaksiUpdatedEventData struct {
	Time time.Time
	Old  TransaksiEventData
	New  TransaksiEventData
}
type TransaksiSoftDeletedEventData struct {
	Time time.Time
	Data TransaksiEventData
}
type TransaksiRestoredEventData TransaksiSoftDeletedEventData
type TransaksiHardDeletedEventData TransaksiSoftDeletedEventData

func NewTransaksiCreatedEventDataFromDispatcher(eventData any) (string, *TransaksiCreatedEventData, error) {
	data, ok := eventData.(TransaksiCreatedEventData)
	if !ok {
		err := errors.New(fmt.Sprintf("failed to cast %s eventdata from any to %s",
			"TransaksiCreated",
			"TransaksiCreatedEventData"))
		return TransaksiCreated, nil, err
	}
	return TransaksiCreated, &data, nil
}
func NewTransaksiUpdatedEventDataFromDispatcher(eventData any) (string, *TransaksiUpdatedEventData, error) {
	data, ok := eventData.(TransaksiUpdatedEventData)
	if !ok {
		err := errors.New(fmt.Sprintf("failed to cast %s eventdata from any to %s",
			"TransaksiUpdated",
			"TransaksiUpdatedEventData"))
		return TransaksiUpdated, nil, err
	}
	return TransaksiUpdated, &data, nil
}
func NewTransaksiSoftDeleteEventDataFromDispatcher(eventData any) (string, *TransaksiSoftDeletedEventData, error) {
	data, ok := eventData.(TransaksiSoftDeletedEventData)
	if !ok {
		err := errors.New(fmt.Sprintf("failed to cast %s eventdata from any to %s",
			"TransaksiSoftDeleted",
			"TransaksiSoftDeletedEventData"))
		return TransaksiSoftDeleted, nil, err
	}
	return TransaksiSoftDeleted, &data, nil
}
func NewTransaksiRestoreEventDataFromDispatcher(eventData any) (string, *TransaksiRestoredEventData, error) {
	data, ok := eventData.(TransaksiRestoredEventData)
	if !ok {
		err := errors.New(fmt.Sprintf("failed to cast %s eventdata from any to %s",
			"TransaksiRestored",
			"TransaksiRestoredEventData"))
		return TransaksiRestored, nil, err
	}
	return TransaksiRestored, &data, nil
}
