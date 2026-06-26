package openai

import (
	timezone "core/pkg/timezone"

	"context"
	"encoding/json"
	"fmt"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
)

const prompt = `
You infer locations and timezones.
Return ONLY valid JSON.

{
	"location": "string",
	"country": "string",
	"timezone": "string"
}

Rules:
- timezone must be a valid IANA timezone.
- country must be the country name.
- never return markdown.
- never explain your reasoning.
`

func (c *ClientStruct) InferTimezone(location string) (*timezone.TimezoneResult, error) {
	response, err := c.Client.Responses.New(
		context.Background(),
		responses.ResponseNewParams{
			Model: openai.ChatModelGPT4_1Nano,
			Input: responses.ResponseNewParamsInputUnion{
				OfString: openai.String(
					prompt + "\n\nLocation:\n" + location,
				),
			},
		},
	)

	if err != nil {
		return nil, err
	}

	content := response.OutputText()
	var result timezone.TimezoneResult

	if 	err := json.Unmarshal([]byte(content), &result); 
		err != nil {
		return nil, fmt.Errorf(
			"failed to parse timezone response: %w\nresponse=%s",
			err,
			content,
		)
	}

	if result.Timezone == "" {
		return nil, fmt.Errorf(
			"openai returned empty timezone",
		)
	}

	return &result, nil
}