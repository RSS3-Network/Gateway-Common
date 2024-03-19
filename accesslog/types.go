package accesslog

import "time"

type Log struct {
	Key       *string   `json:"key,omitempty"` // UUID string of key
	Path      string    `json:"path"`
	Status    int       `json:"status"`
	Timestamp time.Time `json:"ts"`
}
