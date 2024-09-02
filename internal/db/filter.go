package db

import (
	"fmt"
	"math"
	"strings"
)

type Filter struct {
	PageInt        int
	PageSize       int
	SortFields     []string // 排序字段，对应的数据库字段必须属于SortSafeFields，如：[id， user_id, -id, -user_id]
	SortSafeFields []string // 允许排序的安全字段，需要预设，如[id, -id]: id -> ASC, -id -> DESC
}

// sortSQL 对字段进行排序
//
// 如果有返回则是：" order by xx ASC, yy Desc ... "
//
// 没有则返回空串 " "
func (f Filter) sortSQL() string {
	if len(f.SortFields) == 0 || len(f.SortSafeFields) == 0 {
		return " "
	}

	sortList := []string{}

	for _, sortField := range f.SortFields {
		for _, sortSafeField := range f.SortSafeFields {
			if sortField == sortSafeField {
				sortList = append(sortList, sortDirection(sortField))
			}
		}
	}

	if len(sortList) == 0 {
		return " "
	}

	sql := " ORDER BY " + strings.Join(sortList, ",") + " "
	return sql
}

// limit 查询记录数，0<n<=100,超出范围则设为默认值
func (f Filter) limit() int {
	if f.PageSize > 100 {
		f.PageSize = 100
	}
	if f.PageSize <= 0 {
		f.PageSize = 10
	}

	return f.PageSize
}

// offset 查询偏移量，最小为1
func (f Filter) offset() int {
	if f.PageInt <= 0 {
		f.PageInt = 1
	}
	return (f.PageInt - 1) * f.PageSize
}

func (f Filter) limitSQL() string {
	return fmt.Sprintf("LIMIT %d OFFSET %d", f.limit(), f.offset())
}

type Metadata struct {
	CurrentPage  int `json:"currentPage,omitempty"`
	PageSize     int `json:"pageSize,omitempty"`
	FirstPage    int `json:"firstPage,omitempty"`
	LastPage     int `json:"lastPage,omitempty"`
	TotalRecords int `json:"totalRecords,omitempty"`
}

func calculateMetadata(totalRecords, pageInt, pageSize int) Metadata {
	if totalRecords == 0 {

		return Metadata{}
	}

	return Metadata{
		CurrentPage:  pageInt,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
}

func sortDirection(field string) string {
	if strings.HasPrefix(field, "-") {
		field = strings.Replace(field, "-", "", 1)
		return field + " DESC"
	}

	return field + " ASC"
}
