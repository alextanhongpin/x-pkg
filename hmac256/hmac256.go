package hmac256

import (
	"crypto/hmac"
	"crypto/sha256"
)

// https://console.bluemix.net/docs/services/cloud-object-storage/hmac/hmac-signature.html#constructing-an-hmac-signature
// https://dev.to/mathijspim/hmac-authentication-better-protection-for-your-api-4e0

func ValidMAC(message, messageMAC, key []byte) bool {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	// signature := base64.URLEncoding.EncodeToString(expectedMAC)
	return hmac.Equal(messageMAC, expectedMAC)
}
