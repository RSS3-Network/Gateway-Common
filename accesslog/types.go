package accesslog

import "time"

type Log struct {
	KeyID     *string   `json:"key_id,omitempty"` // ID string of key
	Path      string    `json:"path"`
	Status    int       `json:"status"`
	Timestamp time.Time `json:"ts"`
}
