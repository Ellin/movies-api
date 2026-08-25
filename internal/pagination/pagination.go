package pagination

import (
	"fmt"
	"net/url"
	"strconv"
)

type Pagination struct {
	Page     int
	PageSize int
}

// Limit returns the number of results per page (i.e. page size).
func (p *Pagination) Limit() int {
	return p.PageSize
}

// Offset calculates how many entries to skip to reach the specified page.
// The page numbering is 0-based.
func (p *Pagination) Offset() int {
	return p.Page * p.PageSize
}

// Validate enforces pagination parameter limits.
func (p *Pagination) Validate() error {
	if p.Page < 0 {
		return fmt.Errorf("page must not be negative")
	}

	if p.PageSize < 1 || p.PageSize > 100 {
		return fmt.Errorf("size must be in range 1 - 100")
	}

	return nil
}

// Parse parses the pagination parameters from URL query values into a Pagination struct.
func Parse(query url.Values) (Pagination, error) {
	// Set default vaules
	pagination := Pagination{
		Page:     0,
		PageSize: 10,
	}

	// Parse "page" from query
	if queryPage := query.Get("page"); queryPage != "" {
		page, err := strconv.Atoi(queryPage)
		if err != nil {
			return Pagination{}, fmt.Errorf("page must be an integer")
		}
		pagination.Page = page
	}

	// Parse "size" from query
	if querySize := query.Get("size"); querySize != "" {
		size, err := strconv.Atoi(querySize)
		if err != nil {
			return Pagination{}, fmt.Errorf("size must be an integer")
		}
		pagination.PageSize = size
	}

	return pagination, nil
}
