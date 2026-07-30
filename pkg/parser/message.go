package parser

import (
	"encoding/json"
	"strconv"
	"strings"
)

const (
	formatBold          = 1
	formatItalic        = 2
	formatStrikethrough = 4
	formatUnderline     = 8
)

type Format int
func (f *Format) UnmarshalJSON(data []byte) error {
	if (string(data) == `""` || string(data) == "null") {
		*f = 0
		return nil
	}

	var value int
	if  err := json.Unmarshal(data, &value); 
		err != nil {
		return err
	}

	*f = Format(value)
	return nil
}

type lexicalDocument struct {
	Root lexicalNode `json:"root"`
}

type lexicalNode struct {
	Type     string        `json:"type"`
	Text     string        `json:"text"`
	Format   Format        `json:"format"`
	ListType string        `json:"listType"`
	Children []lexicalNode `json:"children"`
}

func ParseLexical(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", nil
	}

	if value[0] != '{' {
		return value, nil
	}

	var document lexicalDocument
	if  err := json.Unmarshal([]byte(value), &document); 
		err != nil {
		return "", err
	}

	var blocks []string
	for _, child := range document.Root.Children {
		if block := renderBlock(child); block != "" {
			blocks = append(blocks, block)
		}
	}

	return strings.TrimSpace(strings.Join(blocks, "\n")), nil
}

func renderBlock(node lexicalNode) string {
	switch node.Type {
		case "paragraph":
			return renderInline(node)
		case "list":
			var lines []string
			for i, item := range node.Children {
				text := strings.TrimSpace(renderInline(item))
				if text == "" { 
					continue 
				}
				if node.ListType == "number" {
					lines = append(
						lines,
						strconv.Itoa(i+1)+". "+text,
					)
				} else {
					lines = append(lines, "• "+text)
				}
			}

			return strings.Join(lines, "\n")
		default:
			return renderInline(node)
	}
}

func renderInline(node lexicalNode) string {
	switch node.Type {
		case "text":
			return applyFormatting(
				node.Text,
				int(node.Format),
			)
		default:
			var builder strings.Builder
			for _, child := range node.Children {
				builder.WriteString(
					renderInline(child),
				)
			}
			return builder.String()
	}
}

func applyFormatting(text string, format int) string {
	result := text
	if format&formatBold != 0 {
		result = toBold(result)
	}

	if format&formatItalic != 0 {
		result = toItalic(result)
	}

	if format&formatUnderline != 0 {
		result = toUnderline(result)
	}

	if format&formatStrikethrough != 0 {
		result = toStrike(result)
	}

	return result
}

func toBold(text string ) string {
	var result strings.Builder
	for _, r := range text {
		switch {
			case r >= '0' && r <= '9':
				result.WriteRune(
					rune(0x1D7CE + (r - '0')),
				)
			case r >= 'A' && r <= 'Z':
				result.WriteRune(
					rune(0x1D5D4 + (r - 'A')),
				)
			case r >= 'a' && r <= 'z':
				result.WriteRune(
					rune(0x1D5EE + (r - 'a')),
				)
			default:
				result.WriteRune(r)
		}
	}

	return result.String()
}

func toItalic(text string) string {
	var result strings.Builder
	for _, r := range text {
		switch {
			case r >= 'A' && r <= 'Z':
				result.WriteRune(
					rune(0x1D608 + (r - 'A')),
				)
			case r >= 'a' && r <= 'z':
				result.WriteRune(
					rune(0x1D622 + (r - 'a')),
				)
			default:
				result.WriteRune(r)
		}
	}

	return result.String()
}

func toUnderline(text string) string {
	var result strings.Builder
	for _, r := range text {
		result.WriteRune(r)
		result.WriteRune('\u0332')
	}

	return result.String()
}

func toStrike(text string) string {
	var result strings.Builder
	for _, r := range text {
		result.WriteRune(r)
		result.WriteRune('\u0336')
	}

	return result.String()
}