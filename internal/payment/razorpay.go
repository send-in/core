package payment

import (
	logger "core/pkg/log"

	"time"
	"bytes"

	"net/http"
	"crypto/hmac"
	"crypto/sha256"

	"encoding/base64"
	"encoding/hex"
	"encoding/json"
)

func CreateOrder(req OrderRequest) (*OrderResponse, error) {
	body, _ := json.Marshal(req)
	request, _ := http.NewRequest(
		"POST",
		"https://api.razorpay.com/v1/orders",
		bytes.NewBuffer(body),
	)

	auth := base64.StdEncoding.EncodeToString(
		[]byte(
			razorpay.KeyID + ":" + razorpay.KeySecret),
	)

	request.Header.Set(
		"Authorization",
		"Basic "+auth,
	)

	request.Header.Set(
		"Content-Type",
		"application/json",
	)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	response, err := client.Do(request)
	if 	response.StatusCode < 200 ||
		response.StatusCode >= 300 {
			
		var body map[string]any
		_ = json.
			NewDecoder(response.Body).
			Decode(&body)

		logger.Error(
			"razorpay returned %d: %v",
			response.StatusCode,
			body,
		)

		return nil, err
	}

	if err != nil { return nil, err }
	defer response.Body.Close()

	var result OrderResponse
	err = json.
		  NewDecoder(response.Body).
		  Decode(&result)

	return &result, err
}

func VerifyWebhook(body []byte, signature string) bool {
	mac := hmac.New(sha256.New,
		[]byte(razorpay.WebhookSecret),
	)

	mac.Write(body)
	expected := hex.EncodeToString(
		mac.Sum(nil),
	)

	return hmac.Equal(
		[]byte(expected),
		[]byte(signature),
	)
}