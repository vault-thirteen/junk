package base

import "math"

func CalculateTotalPages(totalItems Count, pageSize Count) Count {
	return Count(math.Ceil(float64(totalItems) / float64(pageSize)))
}
