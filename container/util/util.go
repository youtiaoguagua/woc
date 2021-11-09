package util

import "strings"

func ShortId(longID string) string {
	prefixSep := strings.IndexRune(longID, ':')
	offset := 0
	length := 12
	if prefixSep >= 0 {
		if longID[0:prefixSep] == "sha256" {
			offset = prefixSep + 1
		} else {
			length += prefixSep + 1
		}
	}

	if len(longID) >= offset+length {
		return longID[offset : offset+length]
	}

	return longID
}
