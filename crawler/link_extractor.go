package crawler

import "regexp"

var (
	exclusiveRegex = regexp.MustCompile(`(?i)\.(?:jpg|jpeg|png|gif|ico|js)$`)
)