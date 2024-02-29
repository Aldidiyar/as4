package queryParameter

import (
	"arch/pkg/type/pagination"
	"arch/pkg/type/sort"
)

type QueryParameter struct {
	Sorts      sort.Sorts
	Pagination pagination.Pagination
}
