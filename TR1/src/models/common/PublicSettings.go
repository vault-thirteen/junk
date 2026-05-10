package cm

type PublicSettings struct {
	Version            string `json:"version"`
	TTL                int    `json:"ttl"`
	SiteName           string `json:"siteName"`
	SiteDomain         string `json:"siteDomain"`
	SessionMaxDuration int    `json:"sessionMaxDuration"`
	MessageEditTime    int    `json:"messageEditTime"`
	PageSize           int    `json:"pageSize"`
}
