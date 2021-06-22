package helper

import (
	"fmt"
	"strconv"
)

// PaginationParameters log param struct
type PaginationParameters struct {
	Page     int64
	StrPage  string
	Limit    int64
	StrLimit string
	Offset   int64
}

const (
	ErrorParameterInvalid = "parameter %s is invalid"
)

// ValidatePagination validates pagination parameters
func ValidatePagination(paging PaginationParameters) (PaginationParameters, error) {
	var err error

	if len(paging.StrPage) > 0 {
		paging.Page, err = strconv.ParseInt(paging.StrPage, 10, 32)
		if err != nil || paging.Page <= 0 {
			return paging, fmt.Errorf(ErrorParameterInvalid, "page")
		}
	}

	if len(paging.StrLimit) > 0 {
		paging.Limit, err = strconv.ParseInt(paging.StrLimit, 10, 32)
		if err != nil || paging.Limit <= 0 {
			return paging, fmt.Errorf(ErrorParameterInvalid, "limit")
		}
	}

	paging.Offset = (paging.Page - 1) * paging.Limit
	return paging, nil
}
