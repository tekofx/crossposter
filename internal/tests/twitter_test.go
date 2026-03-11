package tests

import (
	"testing"

	"github.com/tekofx/crossposter/internal/services/twitter"
)

func TestBaseUrl(t *testing.T) {

}

func TestSignature(t *testing.T) {

	token := "370773112-GmHxMAgYyLbNEtIKZeRNFsMKPR9EyMZeS9weJAEb"

	signatureData := twitter.SignatureData{
		Url:                 "https://api.x.com/1.1/statuses/update.json",
		OauthConsumerKey:    "xvz1evFS4wEEPTGEFPHBog",
		OauthConsumerSecret: "kAcSOqF21Fu85e7zjz7ZN2U4ZRhfV3WpwPAoE3Z7kBw",
		OauthToken:          &token,
		OauthNonce:          "kYjzVBB8Y0ZFabxSWbWovY3uYSQ2pTgmZeNu2VS4cg",
		OauthTimestamp:      1318622958,
		OauthVersion:        "1.0",
		OtherData: map[string]string{
			"include_entities": "true",
			"status":           "Hello Ladies + Gentlemen, a signed OAuth request!",
		},
	}

	result := twitter.CreateSignature(signatureData)
	Assert(t, result == "Ls93hJiZbQ3akF3HF3x1Bz8/zU4=", "Signature creation failed")

}
