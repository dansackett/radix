package radix

import (
	"fmt"
	"log"
	"sort"
	"strings"
)

// Tree represents the actual Radix Tree implementation. It is a group of
// connected nodes which has been compacted unlike a Trie such that the number
// of nodes is decreased to increase memory efficiency.
type Tree struct {
	Root *Node
}

// InitTree creates a new Tree instance ready for insertion
func InitTree() *Tree {
	node := InitNode("*")
	node.IsRootNode = true

	return &Tree{
		Root: node,
	}
}

// InitTreeFromDict creates a new Tree instance loaded with the passed Dictionary
func InitTreeFromDict(dict Dictionary) *Tree {
	words, err := dict.GetWords()

	if err != nil {
		log.Fatal(err)
	}

	tree := InitTree()

	for _, word := range words {
		tree.InsertWord(word)
	}

	return tree
}

// findMatchedNodeMeta searches the tree until it finds a matched child node.
// It returns the metadata about that match including:
//
// - The matched node instance
// - Index in the current node where there is a match against the query
// - The rest of the query if anything is unmatched still
// - The prefix calculated up to the matched node (for storing on nodes with children for faster suggestions)
//
// This function serves as the main recursion algorithm for insertion,
// searching, and suggesting words.
func (t *Tree) findMatchedNodeMeta(query, prefix string, currentNode *Node) (*Node, int, string, string) {
	for _, childNode := range currentNode.Children {
		var matchedIdx int
		matchNotFound := true

		// Setup a "runner" which determines where in the node label and query we stop matching.
		for matchedIdx < len(query) && matchedIdx < len(childNode.Label) && query[matchedIdx] == childNode.Label[matchedIdx] {
			matchedIdx++
			matchNotFound = false
		}

		// Skip to next child node if we don't have a match on this branch
		if matchNotFound {
			continue
		}

		// If we matched the entirety of the label then we need to recurse
		// into the childNode's children
		if matchedIdx == len(childNode.Label) {
			return t.findMatchedNodeMeta(query[matchedIdx:], prefix+query[:matchedIdx], childNode)
		}

		return childNode, matchedIdx, query, prefix
	}

	return currentNode, -1, query, prefix
}

// InsertWord adds a new word to the Radix Trie.
func (t *Tree) InsertWord(word string) {
	matchedNode, matchedIdx, restWord, prefix := t.findMatchedNodeMeta(word, "", t.Root)

	// If we reached the end of the tree and couldn't find anything else to
	// match then we simply add the rest of the word as a new child node to the
	// final matched node instance.
	if matchedIdx == -1 {
		// This word is already in the tree
		if len(restWord) == 0 {
			return
		}
		newNode := InitNode(restWord)
		newNode.IsWordBoundary = true
		matchedNode.Prefix = prefix
		matchedNode.Children = append(matchedNode.Children, newNode)
		return
	}

	// We have a partial prefix match but it will require that we split
	// the current label into two children. We split the current label
	// and word to insert and make new child nodes with the "rest" of
	// each of those strings.
	cachedIsWord := matchedNode.IsWord()
	cachedLabel := matchedNode.Label

	matchedNode.Label = restWord[:matchedIdx]
	matchedNode.Prefix = restWord[:matchedIdx]
	matchedNode.IsWordBoundary = false

	// One important thing to remember is that we want to transfer the
	// existing child nodes to the "rest" of the label so the previous
	// tree remains in tact. We also want to ensure that this new split
	// child is marked as a WordBoundary based on what it was prior to
	// the split.
	restLabelNode := InitNode(cachedLabel[matchedIdx:])
	restLabelNode.Prefix = cachedLabel
	restLabelNode.IsWordBoundary = cachedIsWord
	restLabelNode.Children = matchedNode.Children

	matchedNode.Children = []*Node{restLabelNode}

	// if the word we're inserting is a prefix to the current label then we
	// don't need another branch
	if restWord[matchedIdx:] != "" {
		restWordNode := InitNode(restWord[matchedIdx:])
		restWordNode.IsWordBoundary = true
		matchedNode.Children = append(matchedNode.Children, restWordNode)
	}
}

// Search looks in the tree to see if it can find a word based on a query
func (t *Tree) Search(query string) bool {
	matchedNode, _, restQuery, _ := t.findMatchedNodeMeta(query, "", t.Root)
	return len(restQuery) == 0 && matchedNode.IsWord()
}

// GetSuggestions returns any children that would complete the given search
// query. This is useful for autocomplete.
func (t *Tree) GetSuggestions(query string) []string {
	var suggestions []string

	matchedNode, matchedIdx, restQuery, prefix := t.findMatchedNodeMeta(query, "", t.Root)

	// This is a sign that we don't have a matching node so no suggestions
	if matchedNode == nil || (matchedIdx < len(restQuery) && len(restQuery) > 0) {
		return suggestions
	}

	// If we have a partial prefix match on the current node then we need to
	// update the prefix for iteration. A leaf node will use the current prefix
	// and it's label to finish the suggestion while a node with children has
	// an updated prefix that we should use to iterate.
	if matchedIdx > 0 && matchedIdx < len(matchedNode.Label) && matchedIdx == len(restQuery) && len(restQuery) > 0 {
		if matchedNode.IsLeaf() {
			prefix = prefix + matchedNode.Label
		} else {
			prefix = matchedNode.Prefix
		}
	}

	ch := make(chan string)

	go func(ch chan string, node *Node, prefix string) {
		t.iter(ch, node, prefix)
		close(ch)
	}(ch, matchedNode, prefix)

	for suggestion := range ch {
		suggestions = append(suggestions, suggestion)
	}

	return suggestions
}

// GetSuggestionsForSlice does a concurrent pass on the tree gathering
// suggestions for each search query and aggregating them into a unique slice.
func (t *Tree) GetSuggestionsForSlice(queries []string) []string {
	queriesLeft := len(queries)

	ch := make(chan []string, len(queries))

	for _, query := range queries {
		go func(ch chan []string, tree *Tree, query string) {
			ch <- tree.GetSuggestions(query)
		}(ch, t, query)
	}

	// In order to keep suggestions unique we use a map for gathering the results
	suggestionsMap := make(map[string]bool)

	for matches := range ch {
		queriesLeft--

		for _, match := range matches {
			if ok := suggestionsMap[match]; !ok {
				suggestionsMap[match] = true
			}
		}

		if queriesLeft == 0 {
			break
		}
	}

	close(ch)

	// Our final results will be sized perfectly to optimize memory allocations
	results := make([]string, 0, len(suggestionsMap))

	for suggestion := range suggestionsMap {
		results = append(results, suggestion)
	}

	sort.Strings(results)

	return results
}

// Iter creates a channel for consuming the words that have been added to the
// tree allowing us to display them. It is driven by the helper recursive function.
func (t *Tree) Iter() <-chan string {
	ch := make(chan string)

	go func() {
		t.iter(ch, t.Root, "")
		close(ch)
	}()

	return ch
}

func (t *Tree) iter(out chan<- string, node *Node, currentWord string) {
	if node.IsWord() {
		out <- currentWord
	}

	for _, child := range node.Children {
		t.iter(out, child, currentWord+child.Label)
	}
}

// Debug prints the nodes in the tree showing their depth to check and see if
// the tree is behaving as expected.
func (t *Tree) Debug() {
	t.debug(t.Root, 0)
}

func (t *Tree) debug(node *Node, depth int) {
	if !node.IsRoot() {
		fmt.Println(strings.Repeat("-", depth), node.Label, "|| PREFIX:", node.Prefix)
	}

	for _, child := range node.Children {
		t.debug(child, depth+1)
	}
}
