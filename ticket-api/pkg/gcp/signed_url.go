package gcp

import (
	"encoding/json"
	"time"

	config "ticket-api/config"
	"ticket-api/pkg/cache"

	"cloud.google.com/go/storage"
)

type SignedURL string

// UnmarshalJSON -> convert from json string
func (s *SignedURL) UnmarshalJSON(by []byte) error {
	str := ""
	_ = json.Unmarshal(by, &str)
	*s = SignedURL(str)
	return nil
}

// MarshalJSON -> convert to json string
func (s *SignedURL) MarshalJSON() ([]byte, error) {

	if *s == "" {
		signedURL := ""
		str := "\"" + signedURL + "\""
		return []byte(str), nil
	}

	var signedURL string
	var err error

	cache := cache.GetCacheManager()
	url, _ := cache.Get(string(*s))
	if url == nil {
		signedURL, err = s.getObjectSignedURL()
		if err != nil {
			return []byte("\"\""), nil
		}
	} else {
		signedURL = url.(string)
	}

	str := "\"" + signedURL + "\""
	return []byte(str), nil
}

// GetObjectSignedURL -> gets the signed url for the stored object
func (s *SignedURL) getObjectSignedURL() (string, error) {
	cache := cache.GetCacheManager()

	env := config.GetEnv()
	bucketName := env.StorageBucketName

	conf := config.GetJWTConfig()

	opts := &storage.SignedURLOptions{
		Scheme:         storage.SigningSchemeV4,
		Method:         "GET",
		GoogleAccessID: conf.Email,
		PrivateKey:     conf.PrivateKey,
		Expires:        time.Now().Add(25 * time.Minute),
	}
	if *s == "" {
		return "", nil
	}

	u, err := storage.SignedURL(bucketName, string(*s), opts)
	if err != nil {
		return "", err
	}

	cache.SetWithTTL(string(*s), u, 300, 25*time.Minute)

	return u, nil
}
