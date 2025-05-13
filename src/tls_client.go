package src

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)
import "crypto/tls"

func BcjClient(keys []string) ([]bool, error) {
	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	param, _ := json.Marshal(keys)
	resp, err := client.Post("https://localhost:8848/cache/add", "application/json", bytes.NewBuffer(param))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result []bool
	b, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(b, &result)
	if len(keys) != len(result) {
		return nil, errors.New("keys and results do not match")
	}

	return result, nil
}
