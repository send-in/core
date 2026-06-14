package openai

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
)

type TimezoneResult struct {
	Location string `json:"location"`
	Country  string `json:"country"`
	Timezone string `json:"timezone"`
}

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

func (c *ClientStruct) InferTimezone(location string) (*TimezoneResult, error) {
	response, err := c.Client.Responses.New(
		context.Background(),
		responses.ResponseNewParams{
			Model: openai.ChatModelGPT5Mini,
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
	var result TimezoneResult

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