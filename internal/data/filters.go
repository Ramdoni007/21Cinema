package data

import (
	"strings"

	"github.com/Ramdoni007/21Cinema/internal/validator"
)

type Filters struct {
	Page        int
	PageSize    int
	Sort        string
	SortSfeList []string
}

func ValidateFilters(v *validator.Validator, f Filters) {
	v.Check(f.Page > 0, "page", "page must be greater than zero")

	v.Check(f.Page <= 10_000_000, "page", "page must be maximum than 10 millions")
	v.Check(f.Page > 0, "page_size", "page_size must be greater than zero")
	v.Check(f.Page <= 100, "page_size", "page_size must be greater than a maximum 100")

	v.Check(validator.In(f.Sort, f.SortSfeList...), "sort", "invalid sort value")
}

func (f Filters) sortColumn() string {
	for _, safeValue := range f.SortSfeList {
		if f.Sort == safeValue {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}

	panic("unsafe sort parameter: " + f.Sort)
}

func (f Filters) sortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}

	return "ASC"
}
