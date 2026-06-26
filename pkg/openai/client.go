package openai

import (
	cfg "core/internal/config"
	sdk "github.com/openai/openai-go/v3"
	option "github.com/openai/openai-go/v3/option"
)

type ClientStruct struct {
	sdk.Client
}

var Client *ClientStruct

func Configure(
	config *cfg.OpenAIConfig,
) {
	Client = &ClientStruct{
		Client: sdk.NewClient(
			option.WithAPIKey(
				config.SecretKey,
			),
		),
	}
}