package varName

import (
	// "fmt"
	"testing"
)

func TestMakeTableName(t *testing.T) {
	cases := []struct {
		input  string
		output string
	}{
		{"220 BEA EconData Employment 2010-2015",
			"BeaEconDataEmployment2010"},
	}

	for i, c := range cases {
		p := NewTableNameParams(c.input)
		got := MakeTableName(p)
		if c.output != got {
			t.Errorf("case %d error mismatch: expected %s, got %s", i, c.output, got)
			continue
		}
	}

}

func TestMakeNameUnique(t *testing.T) {
	existing1 := &map[string]bool{
		"BeaEconDataEmployment2010":     true,
		"BeaEconDataEmployment2010_2":   true,
		"BeaEconDataEmployment2010_100": true,
	}

	cases := []struct {
		input    string
		existing *map[string]bool
		output   string
	}{
		{"BeaEconDataEmployment2010",
			existing1,
			"BeaEconDataEmployment2010_101"},
	}

	for i, c := range cases {
		got := MakeNameUnique(c.input, c.existing)
		if c.output != got {
			t.Errorf("case %d error mismatch: expected %s, got %s", i, c.output, got)
			continue
		}
	}

}
