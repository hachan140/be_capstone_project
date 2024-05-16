package events

import (
	"be-capstone-project/src/internal/core/common"
	"fmt"
	"github.com/google/uuid"
	"time"
)

func MakeEventSignInSignUp(eventName, serviceCode string, eventData *LogInData) *Event {
	eventTime := time.Now().UnixNano() / int64(time.Millisecond)
	event := &Event{
		ID:          uuid.New().String(),
		RequestID:   eventData.RequestID,
		Event:       eventName,
		ServiceCode: serviceCode,
		UserID:      fmt.Sprintf("%v", eventData.UserID),
		Timestamp:   eventTime,
		EventTime:   eventTime,
		DeviceID:    eventData.DeviceID,
	}
	event.Payload.Data = eventData
	event.Payload.Meta = Metadata{
		DeviceID:       eventData.DeviceID,
		UserAgent:      eventData.RequestUserAgent,
		RequestID:      eventData.RequestID,
		RequestMethod:  eventData.RequestMethod,
		RequestPath:    eventData.RequestUri,
		RequestPattern: eventData.RequestUri,
		UserIDType:     common.UserMobileType,
	}
	return event
}
