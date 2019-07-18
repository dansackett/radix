package radix

import (
	"bufio"
	"errors"
	"os"
)

// Dictionary defines an interface which we can use to get words for the tree
type Dictionary interface {
	GetWords() ([]string, error)
}

// LinuxDictionary defines the words file in a Linux file system
type LinuxDictionary struct{}

// GetWords for the linux dictionary returns a slice of all of the words stored
// in the Linux dictionary file
func (d *LinuxDictionary) GetWords() ([]string, error) {
	f, err := os.Open("/usr/share/dict/words")
	defer f.Close()

	if err != nil {
		return nil, errors.New("Could not read dictionary file")
	}

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	var words []string

	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return []string{}, err
	}

	return words, nil
}
