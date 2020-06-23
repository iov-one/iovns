package types

import "testing"

func TestSnakeCaseAppend(t *testing.T) {
	cases := map[string]struct {
		Args     []string
		Expected string
	}{
		"success 0": {
			Args:     []string{},
			Expected: "",
		},
		"success 1": {
			Args:     []string{"register_domain", "1"},
			Expected: "register_domain_1",
		},
		"success 5": {
			Args:     []string{"register_domain", "1", "2", "3", "4", "5"},
			Expected: "register_domain_1_2_3_4_5",
		},
	}
	for _, c := range cases {
		r := buildSeedID(c.Args...)
		if c.Expected != r {
			t.Fatalf("expected %s got %s", c.Expected, r)
		}
	}
}
