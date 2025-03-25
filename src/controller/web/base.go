package web

import "regexp"

type ControllerWebBase struct {
}

func (c *ControllerWebBase) stripHTMLTags(input string) string {
	re := regexp.MustCompile("<[^>]*>")
	return re.ReplaceAllString(input, "")
}
