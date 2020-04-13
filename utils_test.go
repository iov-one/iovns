package iovns

import "testing"

func TestSplitAccountKey(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expectedDomain := "domain"
		expectedAccount := "account"
		key := GetAccountKey(expectedDomain, expectedAccount)
		gotDomain, gotAccount := SplitAccountKey([]byte(key))
		if gotDomain != expectedDomain || expectedAccount != gotAccount {
			t.Fatalf("expected: (%s, %s) got: (%s, %s)", expectedDomain, expectedAccount, gotDomain, gotAccount)
		}
	})
}
