package common

import (
	"math"

	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
)

func BuildPaginationResponse(pagination *guestProto.PaginationParams, length int32, offset int32) *guestProto.Pagination {
	var lastPage int32
	if pagination.GetPerPage() != 0 {
		lastPage = int32(math.Ceil(float64(length) / float64(pagination.GetPerPage())))
	} else {
		lastPage = 0
	}

	return &guestProto.Pagination{
		Total:       length,
		PerPage:     pagination.GetPerPage(),
		CurrentPage: pagination.GetCurrentPage(),
		LastPage:    lastPage,
		From:        offset + 1,
		To:          offset + pagination.GetPerPage(),
	}
}
