package model

import (
	"fmt"
	"time"
)

type Page struct {
	UserID   int64     `form:"userid" json:"userid"`
	ObjectID int64     `form:"objectid" json:"objectid"`
	PageFrom int       `form:"pagefrom" json:"pagefrom"` //从哪页开始
	PageSize int       `form:"pagesize" json:"pagesize"` //每页大小
	Key      string    `form:"key" json:"key"`           //关键词
	Asc      string    `form:"asc" json:"asc"`           //asc：“id”  id asc
	Desc     string    `form:"desc" json:"desc"`
	Name     string    `form:"name" json:"name"`
	DateFrom time.Time `form:"datafrom" json:"datafrom"` //时间点1
	DateTo   time.Time `form:"dateto" json:"dateto"`     //时间点2
	Total    int64     `form:"total" json:"total"`
}

func (p *Page) GetPageSize() int {
	if p.PageSize == 0 {
		return 100
	} else {
		return p.PageSize
	}
}

func (p *Page) GetPageFrom() int {
	if p.PageFrom < 0 {
		return 0
	} else {
		return p.PageFrom
	}
}

func (p *Page) GetOrderBy() string {
	if len(p.Asc) > 0 {
		return fmt.Sprintf(" %s asc", p.Asc)
	} else if len(p.Desc) > 0 {
		return fmt.Sprintf(" %s desc", p.Desc)
	} else {
		return ""
	}
}
