package keep_events

import (
	"fmt"
	"github.com/pkg/errors"
	"time"
)

const (
	PosCreated     = "keep.pos.created"
	PosUpdated     = "keep.pos.updated"
	PosSoftDeleted = "keep.pos.softDeleted"
	PosRestored    = "keep.pos.restored"
	PosHardDeleted = "keep.pos.hardDeleted"
)

type PosEventData struct {
	Time     time.Time
	Id       string
	Nama     string
	Urutan   int
	Saldo    int
	ParentId string
	Level    int
	IsShow   bool
	IsLeaf   bool
	Status   string
}
type PosCreatedEventData PosEventData
type PosUpdatedEventData struct {
	Time time.Time
	Old  PosEventData
	New  PosEventData
}
type PosSoftDeletedEventData PosEventData
type PosRestoredEventData PosEventData
type PosHardDeletedEventData PosEventData

func NewPosCreatedEventDataFromDispatcher(eventData any) (string, *PosCreatedEventData, error) {
	data, ok := eventData.(PosCreatedEventData)
	if !ok {
		err := errors.New(fmt.Sprintf("failed to cast %s eventdata from any to %s",
			"PosCreated",
			"PosCreatedEventData"))
		return PosCreated, nil, err
	}
	return PosCreated, &data, nil
}
func NewPosUpdatedEventDataFromDispatcher(eventData any) (string, *PosUpdatedEventData, error) {
	data, ok := eventData.(PosUpdatedEventData)
	if !ok {
		err := errors.New(fmt.Sprintf("failed to cast %s eventdata from any to %s",
			"PosUpdated",
			"PosUpdatedEventData"))
		return PosUpdated, nil, err
	}
	return PosUpdated, &data, nil
}
func NewPosSoftDeleteEventDataFromDispatcher(eventData any) (string, *PosSoftDeletedEventData, error) {
	data, ok := eventData.(PosSoftDeletedEventData)
	if !ok {
		err := errors.New(fmt.Sprintf("failed to cast %s eventdata from any to %s",
			"PosSoftDeleted",
			"PosSoftDeletedEventData"))
		return PosSoftDeleted, nil, err
	}
	return PosSoftDeleted, &data, nil
}
func NewPosRestoreEventDataFromDispatcher(eventData any) (string, *PosRestoredEventData, error) {
	data, ok := eventData.(PosRestoredEventData)
	if !ok {
		err := errors.New(fmt.Sprintf("failed to cast %s eventdata from any to %s",
			"PosRestored",
			"PosRestoredEventData"))
		return PosRestored, nil, err
	}
	return PosRestored, &data, nil
}
