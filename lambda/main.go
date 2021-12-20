package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"
)

type Response events.APIGatewayProxyResponse

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (*Response, error) {
	b64Decoded, err := base64.StdEncoding.DecodeString(request.Body)
	if err != nil {
		return nil, errors.Wrap(err, "could not base64 decode")
	}

	values, err := url.ParseQuery(string(b64Decoded))
	if err != nil {
		return nil, errors.Wrap(err, "could not parse request")
	}
	log.Printf("%+v", values)

	if values.Get("text") == "" {
		return &Response{StatusCode: 200, Body: "you need to send me an lpe emote name (for example lpehihi)"}, nil
	}

	emoteName := values.Get("text")

	gifEmotes := []string{"lperiot", "crocofat", "crocobro"}
	allowedStickers := []string{"lpehihi", "lpehype", "lpekill", "lperee", "lpethink", "lperiot", "crocofat", "crocobro"}

	found := false
	for _, allowed := range allowedStickers {
		if emoteName == allowed {
			found = true
			break
		}

	}
	if !found {
		return &Response{StatusCode: 200, Body: "sorry, only lpe emotes are allowed"}, nil
	}

	var buf bytes.Buffer

	body, err := json.Marshal(map[string]interface{}{
		"response_type": "in_channel",
		"blocks": []interface{}{
			map[string]interface{}{
				"type": "image",
				"title": map[string]interface{}{
					"type":  "plain_text",
					"text":  emoteName,
					"emoji": true,
				},
				"image_url": fmt.Sprintf("https://in.fkn.space/i/stickers/%s", getFilename(gifEmotes, emoteName)),
				"alt_text":  fmt.Sprintf("%s emote", emoteName),
			},
		},
	})
	if err != nil {
		return &Response{StatusCode: 404}, err
	}
	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	return &resp, nil
}

func getFilename(gifEmotes []string, emoteName string) string {
	for _, gifEmote := range gifEmotes {
		if emoteName == gifEmote {
			return fmt.Sprintf("%s.gif", emoteName)
		}
	}
	return fmt.Sprintf("%s.png", emoteName)
}

func main() {
	lambda.Start(Handler)
}
