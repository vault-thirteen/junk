package rm

import "math"

type ItemsPaginated struct {
	Items    any       `json:"items"`
	MetaData *MetaData `json:"metaData,omitempty"`
}

func NewItemsPaginated[T any](pageNumber int, pageSize int, items any, itemsTotalCount int) ItemsPaginated {
	itemsX := items.([]T)

	return ItemsPaginated{
		Items: itemsX,
		MetaData: &MetaData{
			Page: MetaDataPage{
				Number:     pageNumber,
				Size:       pageSize,
				TotalPages: int(math.Ceil(float64(itemsTotalCount) / float64(pageSize))),
			},
			Items: MetaDataItems{
				ItemsOnPage: len(itemsX),
				TotalItems:  itemsTotalCount,
			},
		},
	}
}
