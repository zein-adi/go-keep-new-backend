package keep_events

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_entities"
	"time"
)

const (
	KantongHistoryCreated = "keep.kantong.history.created"
	KantongHistoryUpdated = "keep.kantong.history.updated"
	KantongHistoryDeleted = "keep.kantong.history.deleted"
)

type KantongHistoryCreatedEventData struct {
	Time time.Time
	Data keep_entities.KantongHistory
}
type KantongHistoryUpdatedEventData struct {
	Time time.Time
	Old  keep_entities.KantongHistory
	New  keep_entities.KantongHistory
}
type KantongHistoryDeletedEventData struct {
	Time time.Time
	Data keep_entities.KantongHistory
}

func NewKantongHistoryCreatedEventDataFromDispatcher(eventData any) (string, *KantongHistoryCreatedEventData, error) {
	data, ok := eventData.(KantongHistoryCreatedEventData)
	if !ok {
		err := errors.New(fmt.Sprintf("failed to cast %s eventdata from any to %s",
			"KantongHistoryCreated",
			"KantongHistoryCreatedEventData"))
		return KantongHistoryCreated, nil, err
	}
	return KantongHistoryCreated, &data, nil
}
func NewKantongHistoryUpdatedEventDataFromDispatcher(eventData any) (string, *KantongHistoryUpdatedEventData, error) {
	data, ok := eventData.(KantongHistoryUpdatedEventData)
	if !ok {
		err := errors.New(fmt.Sprintf("failed to cast %s eventdata from any to %s",
			"KantongHistoryUpdated",
			"KantongHistoryUpdatedEventData"))
		return KantongHistoryUpdated, nil, err
	}
	return KantongHistoryUpdated, &data, nil
}
func NewKantongHistoryDeletedEventDataFromDispatcher(eventData any) (string, *KantongHistoryDeletedEventData, error) {
	data, ok := eventData.(KantongHistoryDeletedEventData)
	if !ok {
		err := errors.New(fmt.Sprintf("failed to cast %s eventdata from any to %s",
			"KantongHistoryDeleted",
			"KantongHistoryDeletedEventData"))
		return KantongHistoryDeleted, nil, err
	}
	return KantongHistoryDeleted, &data, nil
}
