package util

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMarkDownHeaderCreation(t *testing.T) {
	testCases := []struct {
		headers       []string
		expectedValue string
	}{
		{nil, ""},
		{[]string{}, ""},
		{[]string{"one"}, "| one |\n| --- |\n"},
		{[]string{"one", "two"}, "| one | two |\n| --- | --- |\n"},
		{[]string{"one", "two", "three"}, "| one | two | three |\n| --- | --- | --- |\n"},
	}

	for _, currentTestCase := range testCases {
		t.Run(strings.Join(currentTestCase.headers, ","), func(t1 *testing.T) {
			require.Equal(t1, CreateTableHeader(currentTestCase.headers...), currentTestCase.expectedValue)
		})
	}
}

func TestCreateTemplateInfoTableTooManyColumns(t *testing.T) {
	table, err := CreateTable([]string{"one", "two"}, [][]string{
		{"a", "b", "c"},
		{"d"},
		{"e", "f", "g"},
		{"h", "i"},
	})

	require.NotNil(t, err)
	require.Empty(t, table)
}

func TestCreateTemplateInfoTable1Column(t *testing.T) {
	table, err := CreateTable([]string{"one"}, [][]string{{"a"}, {"b"}, {"c"}})

	expected := `| one |
| --- |
| a |
| b |
| c |
`

	require.Nil(t, err)
	require.Equal(t, expected, table)
}

func TestCreateTemplateInfoTable2Columns(t *testing.T) {
	table, err := CreateTable([]string{"one", "two"}, [][]string{
		{"a", "b"},
		{"c"},
		{"d", "e"},
	})

	expected := `| one | two |
| --- | --- |
| a | b |
| c |  |
| d | e |
`

	require.Nil(t, err)
	require.Equal(t, expected, table)
}

func TestCreateTemplateInfoTable3Columns(t *testing.T) {
	table, err := CreateTable([]string{"one", "two", "three"}, [][]string{
		{"a", "b", "c"},
		{"d"},
		{"e", "f", "g"},
		{"h", "i"},
	})

	expected := `| one | two | three |
| --- | --- | --- |
| a | b | c |
| d |  |  |
| e | f | g |
| h | i |  |
`

	require.Nil(t, err)
	require.Equal(t, expected, table)
}

func TestCreateCodeBlockKeepsContentVerbatim(t *testing.T) {
	testCases := []struct {
		name    string
		content string
		fence   string
	}{
		{"no special characters", "normal text", "```"},
		{"backslashes", `curl -d '{\"a\":1}' C:\Windows`, "```"},
		{"inline backticks", "text with `backticks` inside", "```"},
		{"triple backticks", "before\n```\nafter", "````"},
		{"longer run", "a ````` b", "``````"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			block := MarkdownFormatter{}.CreateCodeBlock("Response", tc.content, "http")
			require.Equal(t, "\n**Response**\n"+tc.fence+"http\n"+tc.content+"\n"+tc.fence+"\n", block)

			// No line of the content may close the block early.
			for _, line := range strings.Split(tc.content, "\n") {
				require.False(t, strings.HasPrefix(strings.TrimLeft(line, " "), tc.fence), "line %q closes the fence", line)
			}
		})
	}
}
