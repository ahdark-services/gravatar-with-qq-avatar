package utils

import "go.opentelemetry.io/otel/attribute"

type Pagination struct {
	Page     int `json:"page" form:"page" query:"page"`
	PageSize int `json:"page_size" form:"page_size" query:"page_size"`
}

func NewPagination(page, pageSize int) Pagination {
	return Pagination{
		Page:     page,
		PageSize: pageSize,
	}
}

func (p Pagination) Offset() int {
	if p.Page <= 0 {
		return 0
	}

	return (p.Page - 1) * p.PageSize
}

func (p Pagination) Limit() int {
	if p.PageSize <= 0 {
		return 10
	}

	return p.PageSize
}

func (p Pagination) Attributes() []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.Int("page", p.Page),
		attribute.Int("page_size", p.PageSize),
	}
}
