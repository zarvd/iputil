package iputil

import (
	"net/netip"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGroupPrefixesByFamily(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Name       string
		Input      []netip.Prefix
		ExpectedV4 []netip.Prefix
		ExpectedV6 []netip.Prefix
	}{
		{
			Name:       "empty",
			Input:      []netip.Prefix{},
			ExpectedV4: []netip.Prefix{},
			ExpectedV6: []netip.Prefix{},
		},
		{
			Name: "single IPv4",
			Input: []netip.Prefix{
				netip.MustParsePrefix("192.168.0.0/24"),
			},
			ExpectedV4: []netip.Prefix{
				netip.MustParsePrefix("192.168.0.0/24"),
			},
		},
		{
			Name: "single IPv6",
			Input: []netip.Prefix{
				netip.MustParsePrefix("2001:db8::/48"),
			},
			ExpectedV4: []netip.Prefix{},
			ExpectedV6: []netip.Prefix{
				netip.MustParsePrefix("2001:db8::/48"),
			},
		},
		{
			Name: "multiple IPv4",
			Input: []netip.Prefix{
				netip.MustParsePrefix("192.168.0.0/24"),
				netip.MustParsePrefix("192.168.1.0/24"),
			},
			ExpectedV4: []netip.Prefix{
				netip.MustParsePrefix("192.168.0.0/24"),
				netip.MustParsePrefix("192.168.1.0/24"),
			},
		},
		{
			Name: "multiple IPv6",
			Input: []netip.Prefix{
				netip.MustParsePrefix("2001:db8::/48"),
				netip.MustParsePrefix("2001:db8:1::/48"),
			},
			ExpectedV4: []netip.Prefix{},
			ExpectedV6: []netip.Prefix{
				netip.MustParsePrefix("2001:db8::/48"),
				netip.MustParsePrefix("2001:db8:1::/48"),
			},
		},
		{
			Name: "mixed IPv4 and IPv6",
			Input: []netip.Prefix{
				netip.MustParsePrefix("192.168.0.0/24"),
				netip.MustParsePrefix("2001:db8::/48"),
				netip.MustParsePrefix("192.168.1.0/24"),
				netip.MustParsePrefix("2001:db8:1::/48"),
			},
			ExpectedV4: []netip.Prefix{
				netip.MustParsePrefix("192.168.0.0/24"),
				netip.MustParsePrefix("192.168.1.0/24"),
			},
			ExpectedV6: []netip.Prefix{
				netip.MustParsePrefix("2001:db8::/48"),
				netip.MustParsePrefix("2001:db8:1::/48"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			t.Parallel()

			gotV4, gotV6 := GroupPrefixesByFamily(tt.Input)
			assert.ElementsMatch(t, gotV4, tt.ExpectedV4)
			assert.ElementsMatch(t, gotV6, tt.ExpectedV6)
		})
	}
}

func TestContainsPrefix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Name     string
		Prefix   netip.Prefix
		Other    netip.Prefix
		Expected bool
	}{
		{
			Name:     "empty prefix",
			Prefix:   netip.Prefix{},
			Other:    netip.MustParsePrefix("192.168.0.0/24"),
			Expected: false,
		},
		{
			Name:     "empty other",
			Prefix:   netip.MustParsePrefix("192.168.0.0/24"),
			Other:    netip.Prefix{},
			Expected: false,
		},
		{
			Name:     "other prefix is subset",
			Prefix:   netip.MustParsePrefix("192.168.0.0/22"),
			Other:    netip.MustParsePrefix("192.168.0.0/23"),
			Expected: true,
		},
		{
			Name:     "other prefix is superset",
			Prefix:   netip.MustParsePrefix("192.168.0.0/22"),
			Other:    netip.MustParsePrefix("192.168.0.0/21"),
			Expected: false,
		},
		{
			Name:     "non-overlapping prefix",
			Prefix:   netip.MustParsePrefix("192.168.0.0/23"),
			Other:    netip.MustParsePrefix("192.168.3.0/24"),
			Expected: false,
		},
		{
			Name:     "equal prefix",
			Prefix:   netip.MustParsePrefix("192.168.0.0/24"),
			Other:    netip.MustParsePrefix("192.168.0.0/24"),
			Expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.Expected, ContainsPrefix(tt.Prefix, tt.Other))
		})
	}
}

func TestAggregatePrefixes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Name     string
		Input    []netip.Prefix
		Expected []netip.Prefix
	}{
		{
			Name:     "empty",
			Input:    []netip.Prefix{},
			Expected: []netip.Prefix{},
		},
		{
			Name: "single IPv4 prefix",
			Input: []netip.Prefix{
				netip.MustParsePrefix("192.168.0.0/24"),
			},
			Expected: []netip.Prefix{
				netip.MustParsePrefix("192.168.0.0/24"),
			},
		},
		{
			Name: "single IPv6 prefix",
			Input: []netip.Prefix{
				netip.MustParsePrefix("2001:db8::/48"),
			},
			Expected: []netip.Prefix{
				netip.MustParsePrefix("2001:db8::/48"),
			},
		},
		{
			Name: "two adjacent IPv4 prefixes",
			Input: []netip.Prefix{
				netip.MustParsePrefix("192.168.0.0/24"),
				netip.MustParsePrefix("192.168.1.0/24"),
			},
			Expected: []netip.Prefix{
				netip.MustParsePrefix("192.168.0.0/23"),
			},
		},
		{
			Name: "multiple adjacent IPv4 prefixes",
			Input: []netip.Prefix{
				netip.MustParsePrefix("192.168.0.0/24"),
				netip.MustParsePrefix("192.168.1.0/24"),
				netip.MustParsePrefix("192.168.2.0/23"),
			},
			Expected: []netip.Prefix{
				netip.MustParsePrefix("192.168.0.0/22"),
			},
		},
		{
			Name: "two overlapping IPv4 prefixes",
			Input: []netip.Prefix{
				netip.MustParsePrefix("192.168.0.0/24"),
				netip.MustParsePrefix("192.168.0.1/24"),
			},
			Expected: []netip.Prefix{
				netip.MustParsePrefix("192.168.0.0/24"),
			},
		},
		{
			Name: "multiple overlapping IPv4 prefixes",
			Input: []netip.Prefix{
				netip.MustParsePrefix("192.168.0.0/24"),
				netip.MustParsePrefix("192.168.0.1/24"),
				netip.MustParsePrefix("192.168.0.4/24"),
			},
			Expected: []netip.Prefix{
				netip.MustParsePrefix("192.168.0.0/24"),
			},
		},
		{
			Name: "two adjacent IPv6 prefixes",
			Input: []netip.Prefix{
				netip.MustParsePrefix("2001:db8::/48"),
				netip.MustParsePrefix("2001:db8:1::/48"),
			},
			Expected: []netip.Prefix{
				netip.MustParsePrefix("2001:db8::/47"),
			},
		},
		{
			Name: "multiple adjacent IPv6 prefixes",
			Input: []netip.Prefix{
				netip.MustParsePrefix("2001:db8::/48"),
				netip.MustParsePrefix("2001:db8:1::/48"),
				netip.MustParsePrefix("2001:db8:2::/47"),
			},
			Expected: []netip.Prefix{
				netip.MustParsePrefix("2001:db8::/46"),
			},
		},
		{
			Name: "two overlapping IPv6 prefixes",
			Input: []netip.Prefix{
				netip.MustParsePrefix("2001:db8::/48"),
				netip.MustParsePrefix("2001:db8:0:1::/48"),
			},
			Expected: []netip.Prefix{
				netip.MustParsePrefix("2001:db8::/48"),
			},
		},
		{
			Name: "multiple overlapping IPv6 prefixes",
			Input: []netip.Prefix{
				netip.MustParsePrefix("2001:db8::/48"),
				netip.MustParsePrefix("2001:db8:0:1::/48"),
				netip.MustParsePrefix("2001:db8:0:2::/48"),
			},
			Expected: []netip.Prefix{
				netip.MustParsePrefix("2001:db8::/48"),
			},
		},
		{
			Name: "mixed IPv4 and IPv6 prefixes",
			Input: []netip.Prefix{
				netip.MustParsePrefix("192.168.0.0/24"),
				netip.MustParsePrefix("2001:db8::/48"),
				netip.MustParsePrefix("192.168.1.0/24"),
				netip.MustParsePrefix("2001:db8:1::/48"),
			},
			Expected: []netip.Prefix{
				netip.MustParsePrefix("192.168.0.0/23"),
				netip.MustParsePrefix("2001:db8::/47"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			t.Parallel()

			actual := AggregatePrefixes(tt.Input)
			actualString := make([]string, len(actual))
			for i, p := range actual {
				actualString[i] = p.String()
			}

			expectedString := make([]string, len(tt.Expected))
			for i, p := range tt.Expected {
				expectedString[i] = p.String()
			}

			assert.ElementsMatch(t, expectedString, actualString)
		})
	}
}
