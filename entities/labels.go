package entities

import (
	"fmt"
	"strings"

	"github.com/steiler/vlanallocator/utils"
)

type Labels struct {
	labels map[string]string
}

// NewLabels instantiated a new labels instance
// if no predefined labels should be set, provide nil
func NewLabels(l map[string]string) *Labels {
	if l == nil {
		l = map[string]string{}
	}
	return &Labels{
		labels: l,
	}
}

func (l *Labels) SetLabel(k, v string) {
	l.labels[k] = v
}

func (l *Labels) GetLabel(k string) (string, error) {
	var err error = nil
	v, exists := l.labels[k]
	if !exists {
		err = fmt.Errorf("entry %q not found", k)
	}
	return v, err
}

func (l *Labels) Exists(k string) bool {
	_, exists := l.labels[k]
	return exists
}

func (l *Labels) String(indent int) string {
	resultArr := []string{}
	white := utils.GetWhitespaces(indent)
	for k, v := range l.labels {
		resultArr = append(resultArr, fmt.Sprintf("%s%s: %s", white, k, v))
	}
	return strings.Join(resultArr, "\n") + "\n"
}

func (l *Labels) StringOneLine() string {
	return utils.MapStringString2String(l.labels, ": ", ", ")
}

// Equals checks if the both labels (l and l2) contain the exact same labels
func (l *Labels) Equals(l2 *Labels) bool {

	// if the length of the two label maps differs, they can't be equal
	if len(l.labels) != len(l2.labels) {
		return false
	}

	// get key and value from the actualLIndex
	for k, v := range l.labels {
		// Try gettign the same label keys value from the otherlabelSlice
		l2Val, l2Exists := l2.labels[k]
		// if the key does not exist or the value differs return false
		if !l2Exists || l2Val != v {
			return false
		}
	}
	// if everything matches, return true
	return true
}

// Contains checks if all labels provided in l2 are contained and match labels in l
func (l *Labels) Contains(l2 *Labels) bool {
	// iterate over l2 labels
	for k, v := range l2.labels {
		// retrieve l labels
		lVal, exists := l.labels[k]
		// if they do not exist or val differs return false
		if !exists || lVal != v {
			return false
		}
	}
	// all checks succeeded, return true
	return true
}
