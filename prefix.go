package iputil

import (
	"bytes"
	"fmt"
	"net/netip"
	"sort"
)

// GroupPrefixesByFamily groups prefixes by IP family.
// The first return value is the list of IPv4 prefixes, the second is the list of IPv6 prefixes.
func GroupPrefixesByFamily(vs []netip.Prefix) (v4, v6 []netip.Prefix) {
	for _, v := range vs {
		if v.Addr().Is4() {
			v4 = append(v4, v)
		} else {
			v6 = append(v6, v)
		}
	}
	return v4, v6
}

// ContainsPrefix checks if prefix p fully contains prefix o.
// It returns true if o is a subset of p, meaning all addresses in o are also in p.
// This is true when p overlaps with o and p has fewer or equal number of bits than o.
func ContainsPrefix(p, o netip.Prefix) bool {
	return p.Bits() <= o.Bits() && p.Overlaps(o)
}

// mergeAdjacentPrefixes attempts to merge two adjacent prefixes into a single prefix.
// It returns the merged prefix and a boolean indicating success.
// Note: This function only merges adjacent prefixes, not overlapping ones.
func mergeAdjacentPrefixes(p1, p2 netip.Prefix) (netip.Prefix, bool) {
	// Merge neighboring prefixes if possible
	if p1.Bits() != p2.Bits() || p1.Bits() == 0 {
		return netip.Prefix{}, false
	}

	var (
		bits    = p1.Bits()
		p1Bytes = p1.Addr().AsSlice()
		p2Bytes = p2.Addr().AsSlice()
	)
	if bitAt(p1Bytes, bits-1) == 0 {
		setBitAt(p1Bytes, bits-1, 1)
	} else {
		setBitAt(p2Bytes, bits-1, 1)
	}
	if !bytes.Equal(p1Bytes, p2Bytes) {
		return netip.Prefix{}, false
	}

	rv, err := p1.Addr().Prefix(bits - 1)
	if err != nil {
		panic(fmt.Sprintf("unreachable: %s", err))
	}
	return rv, true
}

// aggregatePrefixesForSingleIPFamily merges overlapping or adjacent prefixes into a single prefix.
// The input prefixes must be the same IP family (IPv4 or IPv6).
// For example,
// - [192.168.0.0/32, 192.168.0.1/32] -> [192.168.0.0/31] (adjacent)
// - [192.168.0.0/24, 192.168.0.1/32] -> [192.168.1.0/24] (overlapping)
func aggregatePrefixesForSingleIPFamily(prefixes []netip.Prefix) []netip.Prefix {
	if len(prefixes) <= 1 {
		return prefixes
	}

	sort.Slice(prefixes, func(i, j int) bool {
		addrCmp := prefixes[i].Addr().Compare(prefixes[j].Addr())
		if addrCmp == 0 {
			return prefixes[i].Bits() < prefixes[j].Bits()
		}
		return addrCmp < 0
	})

	var rv = []netip.Prefix{prefixes[0]}

	for i := 1; i < len(prefixes); i++ {
		last, p := rv[len(rv)-1], prefixes[i]
		if ContainsPrefix(last, p) {
			// Skip overlapping prefixes
			continue
		}
		rv = append(rv, p)

		// Merge adjacent prefixes if possible
		for len(rv) >= 2 {
			// Merge the last two prefixes if they are adjacent
			p, ok := mergeAdjacentPrefixes(rv[len(rv)-2], rv[len(rv)-1])
			if !ok {
				break
			}

			// Replace the last two prefixes with the merged prefix
			rv = rv[:len(rv)-2]
			rv = append(rv, p)
		}
	}
	return rv
}

// AggregatePrefixes merges overlapping or adjacent prefixes into a single prefix.
// It combines prefixes that can be represented by a larger, more inclusive prefix.
//
// Examples:
//   - Adjacent:    [192.168.0.0/32, 192.168.0.1/32] -> [192.168.0.0/31]
//   - Overlapping: [192.168.0.0/24, 192.168.0.1/32] -> [192.168.0.0/24]
func AggregatePrefixes(prefixes []netip.Prefix) []netip.Prefix {
	v4, v6 := GroupPrefixesByFamily(prefixes)
	return append(
		aggregatePrefixesForSingleIPFamily(v4),
		aggregatePrefixesForSingleIPFamily(v6)...,
	)
}
