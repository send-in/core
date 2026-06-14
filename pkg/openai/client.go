package openai

import (
	sdk "github.com/openai/openai-go/v3"
	option "github.com/openai/openai-go/v3/option"
)

type Config struct {
	APIKey string
}

type ClientStruct struct {
	sdk.Client
}

var Client *ClientStruct

func Configure(
	config *Config,
) {
	Client = &ClientStruct{
		Client: sdk.NewClient(
			option.WithAPIKey(
				config.APIKey,
			),
		),
	}
}