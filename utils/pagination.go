package utils

import (
	"api-ticket/internal/controllers/http"
	"math"
)

func CalculatePagination(totalCount int, limit, offset *int) http.Paginations {
	pagination := http.Paginations{
		TotalCount: int64(totalCount),
	}

	if limit == nil {
		return pagination
	}

	limitValue := *limit
	pagination.PerPage = limitValue
	pagination.TotalPages = int(math.Ceil(float64(totalCount) / float64(limitValue)))
	pagination.LastPage = pagination.TotalPages

	if offset == nil {
		pagination.CurrentPage = 1
		if pagination.TotalPages > 1 {
			pagination.NextPage = 2
		}
		return pagination
	}

	offsetValue := *offset
	pagination.CurrentPage = int(offsetValue/limitValue) + 1

	if pagination.CurrentPage < pagination.TotalPages {
		pagination.NextPage = pagination.CurrentPage + 1
	}

	if pagination.CurrentPage > 1 {
		pagination.PrevPage = pagination.CurrentPage - 1
	}

	return pagination
}
