package radix

import (
	"testing"
)

func TestTreeInsertAtRoot(t *testing.T) {
	rt := InitTree()
	rt.InsertWord("test")
	rt.InsertWord("slow")

	if len(rt.Root.Children) != 2 {
		t.Errorf("Tree should have 3 children, found %d", len(rt.Root.Children))
	}

	var currentLabel string

	currentLabel = rt.Root.Children[0].Label
	if currentLabel != "test" {
		t.Errorf("Child should have label 'test' found '%s'", currentLabel)
	}

	currentLabel = rt.Root.Children[1].Label
	if currentLabel != "slow" {
		t.Errorf("Child should have label 'slow' found '%s'", currentLabel)
	}
}

func TestTreeInsertExtendedWord(t *testing.T) {
	rt := InitTree()
	rt.InsertWord("test")
	rt.InsertWord("slow")
	rt.InsertWord("slower")

	if len(rt.Root.Children) != 2 {
		t.Errorf("Tree should have 2 children, found %d", len(rt.Root.Children))
	}

	if len(rt.Root.Children[0].Children) != 0 {
		t.Errorf("Child should have no children, found %d", len(rt.Root.Children[0].Children))
	}

	if len(rt.Root.Children[1].Children) != 1 {
		t.Errorf("Child should have 1 child, found %d", len(rt.Root.Children[1].Children))
	}

	var currentLabel string

	currentLabel = rt.Root.Children[0].Label
	if currentLabel != "test" {
		t.Errorf("Child should have label 'test', found %s", currentLabel)
	}

	currentLabel = rt.Root.Children[1].Label
	if currentLabel != "slow" {
		t.Errorf("Child should have label 'slow', found %s", currentLabel)
	}

	currentLabel = rt.Root.Children[1].Children[0].Label
	if currentLabel != "er" {
		t.Errorf("Child should have label 'er', found %s", currentLabel)
	}
}

func TestTreeInsertPrefix(t *testing.T) {
	rt := InitTree()
	rt.InsertWord("tester")
	rt.InsertWord("test")

	if len(rt.Root.Children) != 1 {
		t.Errorf("Tree should have 1 child, found %d", len(rt.Root.Children))
	}

	if len(rt.Root.Children[0].Children) != 1 {
		t.Errorf("Child should have 1 child, found %d", len(rt.Root.Children[0].Children))
	}

	var currentLabel string

	currentLabel = rt.Root.Children[0].Label
	if currentLabel != "test" {
		t.Errorf("Child should have label 'test', found '%s'", currentLabel)
	}

	currentLabel = rt.Root.Children[0].Children[0].Label
	if currentLabel != "er" {
		t.Errorf("Child should have label 'er', found '%s'", currentLabel)
	}
}

func TestTreeInsertSplitNode(t *testing.T) {
	rt := InitTree()
	rt.InsertWord("test")
	rt.InsertWord("team")

	if len(rt.Root.Children) != 1 {
		t.Errorf("Tree should have 1 child, found %d", len(rt.Root.Children))
	}

	if len(rt.Root.Children[0].Children) != 2 {
		t.Errorf("Child should have 2 children, found %d", len(rt.Root.Children[0].Children))
	}

	var currentLabel string

	currentLabel = rt.Root.Children[0].Label
	if currentLabel != "te" {
		t.Errorf("Child should have label 'te', found '%s'", currentLabel)
	}

	currentLabel = rt.Root.Children[0].Children[0].Label
	if currentLabel != "st" {
		t.Errorf("Child should have label 'st', found %s", currentLabel)
	}

	currentLabel = rt.Root.Children[0].Children[1].Label
	if currentLabel != "am" {
		t.Errorf("Child should have label 'am', found %s", currentLabel)
	}
}

func TestTreeInsertSplitPatchNode(t *testing.T) {
	rt := InitTree()
	rt.InsertWord("test")
	rt.InsertWord("team")
	rt.InsertWord("toast")

	if len(rt.Root.Children) != 1 {
		t.Errorf("Tree should have 1 child, found %d", len(rt.Root.Children))
	}

	if len(rt.Root.Children[0].Children) != 2 {
		t.Errorf("Child should have 2 children, found %d", len(rt.Root.Children[0].Children))
	}

	if len(rt.Root.Children[0].Children[0].Children) != 2 {
		t.Errorf("Child should have 2 children, found %d", len(rt.Root.Children[0].Children[0].Children))
	}

	if len(rt.Root.Children[0].Children[1].Children) != 0 {
		t.Errorf("Child should have 0 children, found %d", len(rt.Root.Children[0].Children[1].Children))
	}

	var currentLabel string

	currentLabel = rt.Root.Children[0].Label
	if currentLabel != "t" {
		t.Errorf("Child should have label 't', found %s", currentLabel)
	}

	currentLabel = rt.Root.Children[0].Children[0].Label
	if currentLabel != "e" {
		t.Errorf("Child should have label 'e', found %s", currentLabel)
	}

	currentLabel = rt.Root.Children[0].Children[1].Label
	if currentLabel != "oast" {
		t.Errorf("Child should have label 'oast', found %s", currentLabel)
	}

	currentLabel = rt.Root.Children[0].Children[0].Children[0].Label
	if currentLabel != "st" {
		t.Errorf("Child should have label 'st', found %s", currentLabel)
	}

	currentLabel = rt.Root.Children[0].Children[0].Children[1].Label
	if currentLabel != "am" {
		t.Errorf("Child should have label 'am', found %s", currentLabel)
	}
}

func TestSearchTree(t *testing.T) {
	words := []string{
		"slow",
		"slower",
		"team",
		"test",
		"tester",
		"toast",
		"water",
	}

	rt := InitTree()

	for _, word := range words {
		rt.InsertWord(word)
	}

	for _, word := range words {
		if !rt.Search(word) {
			t.Errorf("`%s` not found in search when it should exist", word)
		}
	}

	nonTreeWords := []string{
		"s",
		"t",
		"",
		"toasted",
		"wafer",
		"tea",
		"slowe",
	}

	for _, word := range nonTreeWords {
		if rt.Search(word) {
			t.Errorf("`%s` found in search when it should NOT exist", word)
		}
	}
}

func TestGetSuggestionsForPrefix(t *testing.T) {
	words := []string{
		"slow",
		"slower",
		"team",
		"teamed",
		"teamedup",
		"test",
		"tester",
		"toast",
		"water",
	}

	rt := InitTree()

	for _, word := range words {
		rt.InsertWord(word)
	}

	expectedSuggestions := map[string][]string{
		"wate":  []string{"water"},
		"teste": []string{"tester"},
		"slowe": []string{"slower"},
		"test":  []string{"test", "tester"},
		"te":    []string{"team", "teamed", "teamedup", "test", "tester"},
		"t":     []string{"team", "teamed", "teamedup", "test", "tester", "toast"},
		"wa":    []string{"water"},
		"was":   []string{},
		"slowl": []string{},
		"teame": []string{"teamed", "teamedup"},
		"teas":  []string{},
	}

	for query, suggestions := range expectedSuggestions {
		for idx, suggestion := range rt.GetSuggestions(query) {
			if len(expectedSuggestions[query]) > 0 {
				if suggestion != expectedSuggestions[query][idx] {
					t.Errorf("Expected '%s' as suggestion for '%s', found '%s'", suggestions[idx], query, suggestion)
				}
			}
		}
	}
}
