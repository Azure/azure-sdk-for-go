package main

import "testing"

func Test_versionle(t *testing.T) {
	const dateWithAlpha, dateWithBeta = "2016-02-01-alpha", "2016-02-01-beta"
	const semVer1dot2, semVer1dot3 = "2018-03-03-1.2", "2018-03-03-1.3"

	testCases := []struct {
		left  string
		right string
		want  bool
	}{
		{"2017-12-01", "2018-03-04", true},
		{"2018-03-04", "2017-12-01", false},
		{semVer1dot2, semVer1dot3, true},
		{semVer1dot3, semVer1dot2, false},
		{semVer1dot2, semVer1dot2, true},
		{dateWithAlpha, dateWithBeta, true},
		{dateWithBeta, dateWithAlpha, false},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			t.Logf("\n Left: %s\nRight: %s", tc.left, tc.right)
			if got, err := versionle(tc.left, tc.right); err != nil {
				t.Error(err)
			} else if got != tc.want {
				t.Logf("\n got: %v\nwant: %v", got, tc.want)
				t.Fail()
			}
		})
	}
}
