package types

import (
	"encoding/json"
	"fmt"
	"strings"
)

type ReviewType uint8

var ErrInvalidReviewType = fmt.Errorf("invalid review type")

const (
	// enum values are in the order of the precedence
	// order of precedence is used in evaluator to coalesce multiple reviews into one
	Skip ReviewType = iota
	Approve
	Comment
	RequestChanges
	Dismiss
)

var (
	reviewTypeNames = map[uint8]string{
		0: "SKIP",
		1: "APPROVE",
		2: "COMMENT",
		3: "REQUEST_CHANGES",
		4: "DISMISS",
	}
	reviewTypeValues = reverseMap(reviewTypeNames)

	reviewStateNames = map[uint8]string{
		0: "SKIP",
		1: "APPROVED",
		2: "COMMENTED",
		3: "CHANGES_REQUESTED",
		4: "DISMISSED",
	}
	reviewStateValues = reverseMap(reviewStateNames)
)

func reverseMap(a map[uint8]string) map[string]uint8 {
	r := make(map[string]uint8, len(a))
	for k, v := range a {
		r[v] = k
	}
	return r
}

func (r ReviewType) String() string {
	return reviewTypeNames[uint8(r)]
}

func ParseReviewType(s string) (ReviewType, error) {
	s = strings.ToUpper(strings.TrimSpace(s))
	value, ok := reviewTypeValues[s]
	if !ok {
		return Skip, fmt.Errorf("%w: %v", ErrInvalidReviewType, s)
	}
	return ReviewType(value), nil
}

func ParseReviewState(s string) (ReviewType, error) {
	s = strings.ToUpper(strings.TrimSpace(s))
	value, ok := reviewStateValues[s]
	if !ok {
		return Skip, fmt.Errorf("%w: %v", ErrInvalidReviewType, s)
	}
	return ReviewType(value), nil
}

func (r ReviewType) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

func (r *ReviewType) UnmarshalJSON(data []byte) (err error) {
	var reviewType string
	if err := json.Unmarshal(data, &reviewType); err != nil {
		return err
	}
	if *r, err = ParseReviewType(reviewType); err != nil {
		return err
	}
	return nil
}

type Review struct {
	Type ReviewType `json:"type"`
	Body string     `json:"body"`
}

type Result struct {
	Track  bool
	Review Review
}
