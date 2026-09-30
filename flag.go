package main

import (
	"strings"
)

type Strings []string

func (a *Strings) String() string {
	if a == nil {
		return ""
	}
	return strings.Join(*a, ",")
}

func (a *Strings) Set(s string) error {
	if s == "" {
		return nil
	}

	for v := range strings.SplitSeq(s, ",") {
		if v == "" {
			continue
		}
		*a = append(*a, v)
	}
	return nil
}
