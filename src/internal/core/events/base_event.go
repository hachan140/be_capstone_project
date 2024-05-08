package events

type Event struct {
	ID          string `json:"id"`
	RequestID   string `json:"request_id"`
	RefEventID  string `json:"ref_event_id"`
	Event       string `json:"event"`
	UserID      string `json:"user_id"`
	ServiceCode string `json:"service_code"`
	Timestamp   int64  `json:"timestamp"`
	EventTime   int64  `json:"event_time"`
	DeviceID    string `json:"device_id"`
	Payload     struct {
		Data      interface{} `json:"data"`
		Meta      interface{} `json:"metadata,omitempty"`
		Attribute string      `json:"attribute,omitempty"`
		Tags      interface{} `json:"tags,omitempty"`
	} `json:"payload"`
	PayloadID string `json:"payload_id"`
}
