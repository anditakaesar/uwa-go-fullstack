package common_test

import (
	"testing"

	"github.com/anditakaesar/uwa-go-fullstack/internal/common"
	"github.com/stretchr/testify/assert"
)

func TestPagination_GetOffset(test *testing.T) {
	test.Parallel()

	test.Run("success return 0", func(t *testing.T) {
		p := common.Pagination{
			Page: 1,
			Size: 10,
		}
		got := p.GetOffset()
		assert.Equal(t, 0, got)
	})

	test.Run("success return non zero", func(t *testing.T) {
		p := common.Pagination{
			Page: 2,
			Size: 10,
		}
		got := p.GetOffset()
		assert.Equal(t, 10, got)
	})
}

func TestSort_ToSQLSort(test *testing.T) {
	test.Parallel()

	test.Run("success ASC", func(t *testing.T) {
		s := common.Sort{
			Field:     "created_at",
			Direction: common.SORT_ASC,
		}
		got := s.ToSQLSort()
		assert.Equal(t, "created_at ASC", got)
	})

	test.Run("success ASC", func(t *testing.T) {
		s := common.Sort{
			Field:     "created_at",
			Direction: common.SORT_DESC,
		}
		got := s.ToSQLSort()
		assert.Equal(t, "created_at DESC", got)
	})
}
