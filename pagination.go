package redis_pagination

import "math"

func newPagination(page int, limit int) *Pagination {
	return &Pagination{
		Page:  page,
		Limit: limit,
	}
}

type Pagination struct {
	Page       int      `json:"page" form:"page"`
	Limit      int      `json:"limit" form:"limit"`
	TotalRows  int64    `json:"total_rows" form:"total_rows"`
	TotalPages int      `json:"total_pages" form:"total_pages"`
	Rows       []string `json:"rows" form:"rows"`
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination) GetLimit() int {
	if p.Limit == 0 {
		p.Limit = 20
	}
	return p.Limit
}

func (p *Pagination) GetEnd() int {
	return p.Page * (p.Limit - 1)
}

func (p *Pagination) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *Pagination) countTotalPages(count int64) {
	p.TotalRows = count
	p.TotalPages = int(math.Ceil(float64(p.TotalRows) / float64(p.GetLimit())))
}
