package utils

import (
	"errors"
	neturl "net/url"
)

func EncodeUrl(url string, query map[string]string) (string, error) {
	if len(query) == 0 {
		return url, nil
	}

	parsedUrl, err := neturl.Parse(url)

	if err != nil {
		return "", errors.New("error parsing url")
	}

	values := parsedUrl.Query()

	for key, value := range query {
		values.Add(key, value)
	}

	parsedUrl.RawQuery = values.Encode()

	return parsedUrl.String(), nil
}
