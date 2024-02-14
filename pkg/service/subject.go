package service

import (
	"fmt"
	"strings"
)

// Subject represents a NATS subject which can include wildcards.
type Subject string

// String converts the Subject back to its string representation.
func (s Subject) String() string {
	return string(s)
}

// Tokens splits the Subject into its constituent parts, divided by '.'.
func (s Subject) Tokens() []string {
	return strings.Split(string(s), ".")
}

// Validate checks if the Subject contains any characters that are not allowed.
func (s Subject) Validate() error {
	if strings.ContainsAny(s.String(), ",? \r\n\t$\b") {
		return fmt.Errorf("invalid subject: %s", s.String())
	}
	return nil
}

// Match tries to pattern-match the subject against another subject.
// It considers NATS wildcard rules where '*' matches any token at a level, and '>' matches all subsequent tokens.
func (s Subject) Match(subject Subject) bool {
	pTokens := s.Tokens()
	sTokens := subject.Tokens()

	for i, pToken := range pTokens {
		if pToken == ">" {
			return true
		}
		if i >= len(sTokens) {
			return false
		}
		if pToken != "*" && pToken != sTokens[i] {
			return false
		}
	}

	return len(pTokens) == len(sTokens) || pTokens[len(pTokens)-1] == "*"
}

// SymmetricMatch tries to pattern-match a subject against another subject in both directions.
// It considers a match if either subject matches the other according to NATS wildcard rules.
func (s Subject) SymmetricMatch(subject Subject) bool {
	return s.Match(subject) || subject.Match(s)
}

// SubjectMap maps subjects (as strings) to indices to an array that holds stream and consumer configs.
type SubjectMap map[Subject]int

// Add inserts a subject and its associated index into the map.
// It returns an error if the subject is invalid, but in the current implementation, it always succeeds.
func (m SubjectMap) Add(subject Subject, idx int) error {
	if err := subject.Validate(); err != nil {
		return err
	}
	m[subject] = idx
	return nil
}

// Search looks for a subject in the map that matches the given subject according to NATS pattern matching rules.
// Returns the matching subject, its associated index, and true if a match is found.
// If no match is found, it returns empty string, 0, and false.
func (m SubjectMap) Search(subject Subject) (Subject, int, bool) {
	for k, v := range m {
		if k.Match(subject) {
			return k, v, true
		}
	}
	return "", 0, false
}

// Get returns the index stored
func (m SubjectMap) Get(subject Subject) (int, bool) {
	v, ok := m[subject]
	return v, ok
}

// SymmetricSearch looks for a subject in the map that symmetrically matches the given subject.
// Symmetric matching means either the map's subject matches the given subject or vice versa.
// Returns the matching subject, its associated index, and true if a match is found.
// If no match is found, it returns empty string, 0, and false.
func (m SubjectMap) SymmetricSearch(subject Subject) (Subject, int, bool) {
	for k, v := range m {
		if k.SymmetricMatch(subject) {
			return k, v, true
		}
	}
	return "", 0, false
}
