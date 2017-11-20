package varName

import (
	// "fmt"
	"testing"
)

func TestMakeTableName(t *testing.T) {
	cases := []struct {
		input  *MakeTableNameParams
		output string
	}{
		{NewTableNameParams("220 BEA EconData Employment 2010-2015"),
			"bea_econ_data_employment_2010"},
		{&MakeTableNameParams{
			InputName:     "220 BEA EconData Employment 2010-2015",
			SkipWords:     &defaultSkipwords,
			Substitutions: &defaultSubstitutions,
			Delim:         " ",
			MaxLen:        20,
			RemoveOnly:    false,
			NoRepeats:     true,
			Alignment:     Right,
			NameCasing:    Kebab,
		},
			"employment-2010-2015"},
		{&MakeTableNameParams{
			InputName:     "220 BEA ===EconData Employment 2010-2015",
			SkipWords:     &defaultSkipwords,
			Substitutions: &defaultSubstitutions,
			Delim:         " ",
			MaxLen:        25,
			RemoveOnly:    true,
			NoRepeats:     true,
			Alignment:     Edge,
			NameCasing:    Camel,
		},
			"BeaEmployment20102015"},
		{&MakeTableNameParams{
			InputName:     "aaa bbb ccc ddd",
			SkipWords:     &defaultSkipwords,
			Substitutions: &defaultSubstitutions,
			Delim:         " ",
			MaxLen:        9,
			RemoveOnly:    true,
			NoRepeats:     true,
			Alignment:     Edge,
			NameCasing:    Camel,
		},
			"aaaBbbDdd"},
	}

	for i, c := range cases {
		p := c.input
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
		{"newName",
			existing1,
			"newName",
		},
	}

	for i, c := range cases {
		got := MakeNameUnique(c.input, c.existing)
		if c.output != got {
			t.Errorf("case %d error mismatch: expected %s, got %s", i, c.output, got)
			continue
		}
	}

}
