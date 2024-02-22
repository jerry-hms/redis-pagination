package redis_pagination

import "math"

func newPagination(opt *Options) *pagination {
	return &pagination{
		Page:     opt.Page,
		PageSize: opt.PageSize,
	}
}

type pagination struct {
	Page       int      `json:"page" form:"page"`
	PageSize   int      `json:"page_size" form:"page_size"`
	TotalRows  int64    `json:"total_rows" form:"total_rows"`
	TotalPages int      `json:"total_pages" form:"total_pages"`
	Rows       []string `json:"rows" form:"rows"`
}

func (p *pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *pagination) GetLimit() int {
	if p.PageSize == 0 {
		p.PageSize = 20
	}
	return p.PageSize
}

func (p *pagination) GetEnd() int {
	return p.Page * (p.PageSize - 1)
}

func (p *pagination) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *pagination) countTotalPages(count int64) {
	p.TotalRows = count
	p.TotalPages = int(math.Ceil(float64(p.TotalRows) / float64(p.GetLimit())))
}
