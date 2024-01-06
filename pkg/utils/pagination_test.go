package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewPagination(t *testing.T) {
	asserts := assert.New(t)

	p := NewPagination(1, 10)
	asserts.Equal(Pagination{
		Page:     1,
		PageSize: 10,
	}, p)
}

func TestPagination_Offset(t *testing.T) {
	asserts := assert.New(t)

	{
		p := Pagination{
			Page:     1,
			PageSize: 10,
		}
		asserts.Equal(0, p.Offset())
	}

	{
		p := Pagination{
			Page:     2,
			PageSize: 10,
		}
		asserts.Equal(10, p.Offset())
	}

	{
		p := Pagination{
			Page:     0,
			PageSize: 10,
		}
		asserts.Equal(0, p.Offset())
	}
}

func TestPagination_Limit(t *testing.T) {
	asserts := assert.New(t)

	{
		p := Pagination{
			Page:     1,
			PageSize: 10,
		}
		asserts.Equal(10, p.Limit())
	}

	{
		p := Pagination{
			Page:     1,
			PageSize: 0,
		}
		asserts.Equal(10, p.Limit())
	}

	{
		p := Pagination{
			Page:     1,
			PageSize: -1,
		}
		asserts.Equal(10, p.Limit())
	}
}
