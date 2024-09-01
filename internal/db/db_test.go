package db

import "testing"

func TestUnqQuerySQL(t *testing.T) {
	var fields = make(map[string]string, 0)
	fields["id"] = "100"
	fields["phone"] = ""
	fields["name"] = "fzzp"

	sql, params := unqQuerySQL(fields, "$")
	t.Log(sql)
	t.Log(params...)
	t.Log("...")
}
