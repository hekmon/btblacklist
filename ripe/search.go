package ripe

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Range contains all information about an IP range
type Range struct {
	Name  string
	Range string
	Route string
}

// Search returns a list of ranges based on a search string
func Search(search string) (ranges []Range, err error) {
	// Prepare URL
	searchSplitted := strings.Split(search, " ")
	url := *baseURL
	queryValues := url.Query()
	queryValues.Set("q", "("+strings.Join(searchSplitted, " AND ")+")")
	url.RawQuery = queryValues.Encode()
	// Start request
	resp, err := client.Get(url.String())
	if err != nil {
		return
	}
	defer resp.Body.Close()
	// Decode
	var payload result
	decoder := json.NewDecoder(resp.Body)
	decoder.DisallowUnknownFields()
	if err = decoder.Decode(&payload); err != nil {
		err = fmt.Errorf("can't decode api answer payload: %v", err)
		return
	}
	// Extract results
	return extractRange(payload)
}

// RemoveDuplicates return a uniq list of range
func RemoveDuplicates(rangesSets [][]Range) (uniqRanges []Range) {
	// count
	var totalRanges int
	for _, set := range rangesSets {
		totalRanges += len(set)
	}
	// search
	uniqRanges = make([]Range, 0, totalRanges)
	for _, set := range rangesSets {
	candidate:
		for _, rangeCandidate := range set {
			for _, addedRange := range uniqRanges {
				if rangeCandidate == addedRange {
					continue candidate
				}
			}
			uniqRanges = append(uniqRanges, rangeCandidate)
		}
	}
	return
}
