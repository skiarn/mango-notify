package models

import (
	"sort"
	"testing"
)

func TestSortByNumberASC(t *testing.T) {
	lines := Lines{
		Line{Number: 4, Content: "d", Sent: false},
		Line{Number: 2, Content: "b", Sent: false},
		Line{Number: 3, Content: "c", Sent: false},
		Line{Number: 1, Content: "a", Sent: false},
	}

	sort.Sort(lines)

	var itemNr int
	for i, l := range lines {
		if i == 0 {
			itemNr = l.Number
		} else {
			if itemNr > l.Number {
				t.Errorf("Expected:'%s' to be less but got:'%s'", itemNr, l.Number)
			}
			itemNr = l.Number
		}

	}

}
