package config

import (
	"encoding/json"
	"testing"
	"time"
)

func TestPzzDuration(t *testing.T) {
	var pzz struct {
		PzzTime PzzDuration `json:"pzzTime"`
	}

	var jsonString = `{"pzzTime": "2m"}`

	if err := json.Unmarshal([]byte(jsonString), &pzz); err != nil {
		t.Error(err)
	}

	t.Log(time.Duration(pzz.PzzTime.Duration))

	dur, _ := time.ParseDuration("2m")
	t.Log(dur)

	if time.Duration(pzz.PzzTime.Duration) != dur {
		t.Error("want equality, but not equality")
	}
}
