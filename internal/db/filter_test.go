package db

import (
	"fmt"
	"strings"
	"testing"
)

func TestXXX(t *testing.T) {
	field := "-name-xx"
	field = strings.Replace(field, "-", "", 1)
	fmt.Println(field) // name-xx
}

func TestSortColumn(t *testing.T) {
	f := Filter{
		SortFields:     []string{"-id", "updated_at", "age"},
		SortSafeFields: []string{"id", "-id", "updated_at", "name"},
	}

	sql := f.sortColumn()
	fmt.Printf("%q\n", sql)
}
