package collection

import (
	"sort"
)

type trieNode struct {
	IsWord   bool
	Children map[rune]*trieNode
}

func (node trieNode) Population() uint {
	var sum uint

	for _, child := range node.Children {
		sum += child.Population()
	}

	if node.IsWord {
		sum++
	}

	return sum
}

func (node *trieNode) Navigate(word string) *trieNode {
	cursor := node
	for len(word) > 0 && cursor != nil {
		if next, ok := cursor.Children[rune(word[0])]; ok {
			cursor = next
			word = word[1:]
		} else {
			return nil
		}
	}
	return cursor
}

// Dictionary is a list of words. It is implemented as a Trie for memory efficiency.
type Dictionary struct {
	root *trieNode
	size int64
}

// Add inserts a word into the dictionary, and returns whether or not that word was a new word.
//
// Time complexity: O(m) where 'm' is the length of word.
func (dict *Dictionary) Add(word string) (wasAdded bool) {
	if dict.root == nil {
		dict.root = &trieNode{}
	}

	cursor := dict.root

	for len(word) > 0 {
		if cursor.Children == nil {
			cursor.Children = make(map[rune]*trieNode)
		}

		nextLetter := rune(word[0])

		next, ok := cursor.Children[nextLetter]
		if !ok {
			next = &trieNode{}
			cursor.Children[nextLetter] = next
		}
		cursor = next
		word = word[1:]
	}
	wasAdded = !cursor.IsWord
	if wasAdded {
		dict.size++
	}
	cursor.IsWord = true
	return
}

// Clear removes all items from the dictionary.
func (dict *Dictionary) Clear() {
	dict.root = nil
	dict.size = 0
}

// Contains searches the Dictionary to see if the specified word is present.
//
// Time complexity: O(m) where 'm' is the length of word.
func (dict Dictionary) Contains(word string) bool {
	if dict.root == nil {
		return false
	}
	targetNode := dict.root.Navigate(word)
	return targetNode != nil && targetNode.IsWord
}

// Remove ensures that `word` is not in the Dictionary. Returns whether or not an item was removed.
//
// Time complexity: O(m) where 'm' is the length of word.
func (dict *Dictionary) Remove(word string) (wasRemoved bool) {
	lastPos := len(word) - 1
	parent := dict.root.Navigate(word[:lastPos])
	if parent == nil {
		return
	}

	lastLetter := rune(word[lastPos])

	subject, ok := parent.Children[lastLetter]
	if !ok {
		return
	}

	wasRemoved = subject.IsWord

	if wasRemoved {
		dict.size--
	}

	subject.IsWord = false
	if subject.Population() == 0 {
		delete(parent.Children, lastLetter)
	}
	return
}

// Size reports the number of words there are in the Dictionary.
//
// Time complexity: O(1)
func (dict Dictionary) Size() int64 {
	return dict.size
}

// Enumerate lists each word in the Dictionary alphabetically.
func (dict Dictionary) Enumerate(cancel <-chan struct{}) Enumerator {
	if dict.root == nil {
		return Empty.Enumerate(cancel)
	}
	return dict.root.Enumerate(cancel)
}

func (node trieNode) Enumerate(cancel <-chan struct{}) Enumerator {
	var enumerateHelper func(trieNode, string)

	results := make(chan interface{})

	enumerateHelper = func(subject trieNode, prefix string) {
		if subject.IsWord {
			select {
			case results <- prefix:
			case <-cancel:
				return
			}
		}

		alphabetizedChildren := []rune{}
		for letter := range subject.Children {
			alphabetizedChildren = append(alphabetizedChildren, letter)
		}
		sort.Slice(alphabetizedChildren, func(i, j int) bool {
			return alphabetizedChildren[i] < alphabetizedChildren[j]
		})

		for _, letter := range alphabetizedChildren {
			enumerateHelper(*subject.Children[letter], prefix+string(letter))
		}
	}

	go func() {
		defer close(results)
		enumerateHelper(node, "")
	}()

	return results
}
