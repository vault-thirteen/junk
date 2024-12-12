package rpc

import (
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
)

type PageData struct {
	// Number of current page of items.
	PageNumber cmb.Count `json:"pageNumber"`

	// Number of all available pages.
	TotalPages cmb.Count `json:"totalPages"`

	// Number of items on a full page.
	PageSize cmb.Count `json:"pageSize"`

	// Number of items on the current page.
	ItemsOnPage cmb.Count `json:"itemsOnPage"`

	// Total number of all items.
	TotalItems cmb.Count `json:"totalItems"`
}
