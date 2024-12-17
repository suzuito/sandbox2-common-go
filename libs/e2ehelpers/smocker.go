package e2ehelpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/smocker-dev/smocker/server/types"
)

type SmockerClient struct {
	baseURL *url.URL
	client  *http.Client
}

func (t *SmockerClient) PostMocks(
	body types.Mocks,
	reset bool,
) error {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to json.Marshal: %w", err)
	}

	reqURL, _ := url.Parse(t.baseURL.String())
	reqURL.Path = "/mocks"
	query := reqURL.Query()
	query.Set("reset", strconv.FormatBool(reset))
	reqURL.RawQuery = query.Encode()

	req, err := http.NewRequest(http.MethodPost, reqURL.String(), bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to http.NewRequest: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := t.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to http request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		resBodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			resBodyBytes = []byte{}
		}
		return fmt.Errorf(
			"http error: status=%d body=%s",
			res.StatusCode, string(resBodyBytes),
		)
	}

	return nil
}

func NewSmockerClient(
	baseURL *url.URL,
	client *http.Client,
) *SmockerClient {
	return &SmockerClient{
		baseURL: baseURL,
		client:  client,
	}
}
