package keep_events

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_entities"
	"time"
)

const (
	KantongCreated     = "keep.kantong.created"
	KantongUpdated     = "keep.kantong.updated"
	KantongSoftDeleted = "keep.kantong.softDeleted"
	KantongRestored    = "keep.kantong.restored"
	KantongHardDeleted = "keep.kantong.hardDeleted"
)

type KantongEventData struct {
	Time time.Time
	Data keep_entities.Kantong
}
type KantongCreatedEventData KantongEventData
type KantongUpdatedEventData struct {
	Time time.Time
	Old  keep_entities.Kantong
	New  keep_entities.Kantong
}
type KantongSoftDeletedEventData KantongEventData
type KantongRestoredEventData KantongEventData
type KantongHardDeletedEventData KantongEventData

func NewKantongCreatedEventDataFromDispatcher(eventData any) (string, *KantongCreatedEventData, error) {
	data, ok := eventData.(KantongCreatedEventData)
	if !ok {
		err := errors.New(fmt.Sprintf("failed to cast %s eventdata from any to %s",
			"KantongCreated",
			"KantongCreatedEventData"))
		return KantongCreated, nil, err
	}
	return KantongCreated, &data, nil
}
func NewKantongUpdatedEventDataFromDispatcher(eventData any) (string, *KantongUpdatedEventData, error) {
	data, ok := eventData.(KantongUpdatedEventData)
	if !ok {
		err := errors.New(fmt.Sprintf("failed to cast %s eventdata from any to %s",
			"KantongUpdated",
			"KantongUpdatedEventData"))
		return KantongUpdated, nil, err
	}
	return KantongUpdated, &data, nil
}
func NewKantongSoftDeleteEventDataFromDispatcher(eventData any) (string, *KantongSoftDeletedEventData, error) {
	data, ok := eventData.(KantongSoftDeletedEventData)
	if !ok {
		err := errors.New(fmt.Sprintf("failed to cast %s eventdata from any to %s",
			"KantongSoftDeleted",
			"KantongSoftDeletedEventData"))
		return KantongSoftDeleted, nil, err
	}
	return KantongSoftDeleted, &data, nil
}
func NewKantongRestoreEventDataFromDispatcher(eventData any) (string, *KantongRestoredEventData, error) {
	data, ok := eventData.(KantongRestoredEventData)
	if !ok {
		err := errors.New(fmt.Sprintf("failed to cast %s eventdata from any to %s",
			"KantongRestored",
			"KantongRestoredEventData"))
		return KantongRestored, nil, err
	}
	return KantongRestored, &data, nil
}
