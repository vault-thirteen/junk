// ElectionRecord.go.

package eom

import "time"

type ElectionRecord struct {
	ServiceInstanceID string    `json:"serviceInstanceID"`
	LastUpdateTime    time.Time `json:"lastUpdateTime"`
}

type ElectionRecordRaw struct {
	ServiceInstanceID string
	LastUpdateTime    string
}
