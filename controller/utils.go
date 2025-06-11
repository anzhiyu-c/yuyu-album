package controller

import "strings"

func joinTags(tags []string) string {
	return strings.Join(tags, ",")
}
