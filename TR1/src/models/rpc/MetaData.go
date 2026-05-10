package rm

type MetaData struct {
	Page  MetaDataPage  `json:"page"`
	Items MetaDataItems `json:"items"`
}

type MetaDataPage struct {
	// Number of the current page of items.
	Number int `json:"number"`

	// Size of the current page of items.
	Size int `json:"size"`

	// Total number of all pages.
	TotalPages int `json:"totalCount"`
}

type MetaDataItems struct {
	// Number of items on the current page.
	// It can be less than page size when there are not enough items.
	ItemsOnPage int `json:"onPage"`

	// Total number of all items on all pages.
	TotalItems int `json:"total"`
}
