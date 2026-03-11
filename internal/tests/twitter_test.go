package tests

import (
	"fmt"
	"testing"

	"github.com/tekofx/crossposter/internal/logger"
	"github.com/tekofx/crossposter/internal/services/twitter"
)

var data = setupData()

func setupData() twitter.SignatureData {
	oauthToken := "370773112-GmHxMAgYyLbNEtIKZeRNFsMKPR9EyMZeS9weJAEb"
	oauthTokenSecret := "LswwdoUaIvS8ltyTt5jkRh4J50vUPVVHtR2YPi5kE"

	signatureData := twitter.SignatureData{
		HttpMethod:          "POST",
		Url:                 "https://api.x.com/1.1/statuses/update.json",
		OauthConsumerKey:    "xvz1evFS4wEEPTGEFPHBog",
		OauthConsumerSecret: "kAcSOqF21Fu85e7zjz7ZN2U4ZRhfV3WpwPAoE3Z7kBw",
		OauthToken:          &oauthToken,
		OauthTokenSecret:    &oauthTokenSecret,
		OauthNonce:          "kYjzVBB8Y0ZFabxSWbWovY3uYSQ2pTgmZeNu2VS4cg",
		OauthTimestamp:      1318622958,
		OauthVersion:        "1.0",
		PathParameters: map[string]string{
			"include_entities": "true",
		},
		BodyParameters: map[string]string{
			"status": "Hello Ladies + Gentlemen, a signed OAuth request!",
		},
	}

	return signatureData
}

func TestBaseUrl(t *testing.T) {
	expected := "POST&https%3A%2F%2Fapi.x.com%2F1.1%2Fstatuses%2Fupdate.json&include_entities%3Dtrue%26oauth_consumer_key%3Dxvz1evFS4wEEPTGEFPHBog%26oauth_nonce%3DkYjzVBB8Y0ZFabxSWbWovY3uYSQ2pTgmZeNu2VS4cg%26oauth_signature_method%3DHMAC-SHA1%26oauth_timestamp%3D1318622958%26oauth_token%3D370773112-GmHxMAgYyLbNEtIKZeRNFsMKPR9EyMZeS9weJAEb%26oauth_version%3D1.0%26status%3DHello%2520Ladies%2520%252B%2520Gentlemen%252C%2520a%2520signed%2520OAuth%2520request%2521"
	result := twitter.CreateSignatureBaseUrl(data)
	Assert(t, result == expected, fmt.Sprintf("\nExpected:%s\n\nGot:%s\n", expected, result))
}

func TestSignature(t *testing.T) {
	expected := "Ls93hJiZbQ3akF3HF3x1Bz8/zU4="
	result := twitter.CreateSignature(data)
	Assert(t, result == expected, fmt.Sprintf("\nExpected:%s\nGot:%s", expected, result))
}

func TestAuthHeader(t *testing.T) {
	expected := `OAuth oauth\_consumer\_key="xvz1evFS4wEEPTGEFPHBog", oauth\_nonce="kYjzVBB8Y0ZFabxSWbWovY3uYSQ2pTgmZeNu2VS4cg", oauth\_signature="tnnArxj06cWHq44gCs1OSKk%2FjLY%3D", oauth\_signature\_method="HMAC-SHA1", oauth\_timestamp="1318622958", oauth\_token="370773112-GmHxMAgYyLbNEtIKZeRNFsMKPR9EyMZeS9weJAEb", oauth_version="1.0"`
	result, err := twitter.CreateAuthHeader(data)
	if err != nil {
		logger.Error(err)
		t.FailNow()
	}

	Assert(t, result == &expected, fmt.Sprintf("\nExpected:\n%s\n\nGot:\n%s", expected, *result))

}
