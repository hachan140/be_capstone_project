package events

type LogIn struct {
	Data *LogInData `json:"data"`
	Meta *Metadata  `json:"meta,omitempty"`
}

type LogInData struct {
	UserID           interface{} `json:"user_id"`
	Source           string      `json:"source"`
	Channel          string      `json:"channel"`
	IsNewInstall     bool        `json:"is_new_install"`
	IsFirstLogin     bool        `json:"is_first_login"`
	IsNewUser        bool        `json:"is_new_user"`
	PhoneNumber      string      `json:"phone_number"`
	HasPassword      bool        `json:"has_password"`
	RequestID        string      `json:"request_id,omitempty"`
	DeviceID         string      `json:"device_id,omitempty"`
	DeviceIDTnS      string      `json:"device_id_tns"`
	LoginTime        int64       `json:"login_time"`
	RequestUserAgent string      `json:"request_user_agent"`
	RequestMethod    string      `json:"request_method"`
	RequestUri       string      `json:"request_uri"`
}

type Metadata struct {
	DeviceID       string `json:"device_id,omitempty"`
	UserAgent      string `json:"user_agent"`
	RequestID      string `json:"request_id,omitempty"`
	RequestMethod  string `json:"request_method"`
	RequestPath    string `json:"request_path"`
	RequestPattern string `json:"request_pattern"`
	UserIDType     string `json:"user_id_type"`
}
