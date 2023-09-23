package data

import (
	"math"
	"strings"

	"github.com/Ramdoni007/21Cinema/internal/validator"
)

type Filters struct {
	Page        int
	PageSize    int
	Sort        string
	SortSfeList []string
}

// Define a new Metadata struct for holding the pagination metadata.
type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

// The calculateMetadata() function calculates the appropriate pagination metadata
// values given the total number of records, current page, and page size values. Note
// that the last page value is calculated using the math.Ceil() function, which rounds
// up a float to the nearest integer. So, for example, if there were 12 records in total
// and a page size of 5, the last page value would be math.Ceil(12/5) = 3.
func calculateMetaData(totalrecords, page, pageSize int) Metadata {
	if totalrecords == 0 {
		// Note that we return an empty Metadata struct if there are no records.
		return Metadata{}
	}

	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalrecords) / float64(pageSize))),
		TotalRecords: totalrecords,
	}
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

func (f Filters) limit() int {
	return f.PageSize
}

func (f Filters) offset() int {
	return (f.Page - 1) * f.PageSize
}
