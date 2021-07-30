package gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	app "github.com/kujilabo/cocotola-api/pkg_app/domain"
	"github.com/kujilabo/cocotola-api/pkg_plugin/common/domain"
)

type synthesizer struct {
	client *http.Client
	key    string
}

type synthesizeResponse struct {
	AudioContent string `json:"audioContent"`
}

func NewSynthesizer(key string, timeout time.Duration) domain.Synthesizer {
	return &synthesizer{
		client: &http.Client{
			Timeout: timeout,
		},
		key: key,
	}
}

func (s *synthesizer) Synthesize(lang app.Lang5, text string) (string, error) {
	type m map[string]interface{}

	values := m{
		"input": m{
			"text": text,
		},
		"voice": m{
			"languageCode": string(lang),
			"ssmlGender":   "FEMALE",
		},
		"audioConfig": m{
			"audioEncoding": "MP3",
			"speakingRate":  1,
		},
	}

	b, err := json.Marshal(values)
	if err != nil {
		return "", err
	}

	u, err := url.Parse("https://texttospeech.googleapis.com/v1/text:synthesize")
	if err != nil {
		return "", err
	}

	q := u.Query()
	q.Set("key", s.key)
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("POST", u.String(), bytes.NewReader(b))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%s", string(body))
	}

	response := synthesizeResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}

	return response.AudioContent, nil
}
