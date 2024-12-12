package s

import (
	"errors"
	base2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	c "github.com/vault-thirteen/SimpleBB/pkg/common/models/server"
)

// SystemSettings are system settings.
type SystemSettings struct {
	NotificationTtl        base2.Count `json:"notificationTtl"`
	NotificationCountLimit base2.Count `json:"notificationCountLimit"`
	PageSize               base2.Count `json:"pageSize"`
	DKeySize               base2.Count `json:"dKeySize"`

	// This setting must be synchronised with settings of the Gateway module.
	IsTableOfIncidentsUsed base2.Flag `json:"isTableOfIncidentsUsed"`

	// This setting is used only when a table of incidents is enabled.
	BlockTimePerIncident BlockTimePerIncident `json:"blockTimePerIncident"`

	IsDebugMode base2.Flag `json:"isDebugMode"`
}

// BlockTimePerIncident is block time in seconds for each type of incident.
type BlockTimePerIncident struct {
	IllegalAccessAttempt            base2.Count `json:"illegalAccessAttempt"`            // 1.
	ReadingNotificationOfOtherUsers base2.Count `json:"readingNotificationOfOtherUsers"` // 2.
	WrongDKey                       base2.Count `json:"wrongDKey"`                       // 3.
}

func (s SystemSettings) Check() (err error) {
	if (s.NotificationTtl == 0) ||
		(s.NotificationCountLimit == 0) ||
		(s.PageSize == 0) ||
		(s.DKeySize == 0) {
		return errors.New(c.MsgSystemSettingError)
	}

	// Incidents.
	if s.IsTableOfIncidentsUsed {
		if (s.BlockTimePerIncident.IllegalAccessAttempt == 0) ||
			(s.BlockTimePerIncident.ReadingNotificationOfOtherUsers == 0) ||
			(s.BlockTimePerIncident.WrongDKey == 0) {
			return errors.New(c.MsgSystemSettingError)
		}
	}

	return nil
}
