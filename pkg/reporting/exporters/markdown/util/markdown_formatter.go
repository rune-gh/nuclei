package util

import (
	"fmt"
	"strings"
)

type MarkdownFormatter struct{}

func (markdownFormatter MarkdownFormatter) MakeBold(text string) string {
	return MakeBold(text)
}

func (markdownFormatter MarkdownFormatter) CreateCodeBlock(title string, content string, language string) string {
	fence := codeFence(content)
	return fmt.Sprintf("\n%s\n%s%s\n%s\n%s\n", markdownFormatter.MakeBold(title), fence, language, content, fence)
}

func (markdownFormatter MarkdownFormatter) CreateTable(headers []string, rows [][]string) (string, error) {
	return CreateTable(headers, rows)
}

func (markdownFormatter MarkdownFormatter) CreateLink(title string, url string) string {
	return CreateLink(title, url)
}

func (markdownFormatter MarkdownFormatter) CreateHorizontalLine() string {
	return CreateHorizontalLine()
}

// codeFence returns a backtick fence for a fenced code block holding content.
//
// Markdown does not process backslash escapes inside a code block, so escaping
// the content would show the backslashes in the rendered report and change
// requests, responses and curl commands. Instead the fence is made longer than
// the longest run of backticks in the content, so no line of the content can
// close the block early (CommonMark and GitHub Flavored Markdown).
func codeFence(content string) string {
	longest, run := 0, 0
	for i := 0; i < len(content); i++ {
		if content[i] == '`' {
			run++
			if run > longest {
				longest = run
			}
			continue
		}
		run = 0
	}
	return strings.Repeat("`", max(3, longest+1))
}
