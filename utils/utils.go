package utils

import (
	"html/template"
	"strings"
)

func Indent(amount int, html template.HTML) template.HTML {
	// Convert template.HTML back to string for processing
	str := string(html)

	pad := strings.Repeat(" ", amount)
	// Indent each line by adding the padding
	indentedStr := pad + strings.ReplaceAll(str, "\n", "\n"+pad)

	// Return the indented string as template.HTML
	return template.HTML(indentedStr)
}

func SafeJSON(input string) (template.HTML, error) {
	return template.HTML(input), nil
}
