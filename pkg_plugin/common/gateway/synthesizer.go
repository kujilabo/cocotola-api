package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	app "github.com/kujilabo/cocotola-api/pkg_app/domain"
	"github.com/kujilabo/cocotola-api/pkg_plugin/common/service"
	"golang.org/x/xerrors"
)

type synthesizer struct {
	client *http.Client
	key    string
}

type synthesizeResponse struct {
	AudioContent string `json:"audioContent"`
}

func NewSynthesizer(key string, timeout time.Duration) service.Synthesizer {
	return &synthesizer{
		client: &http.Client{
			Timeout: timeout,
		},
		key: key,
	}
}

func (s *synthesizer) Synthesize(ctx context.Context, lang app.Lang5, text string) (string, error) {
	type m map[string]interface{}

	values := m{
		"input": m{
			"text": text,
		},
		"voice": m{
			"languageCode": lang.String(),
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

	req, err := http.NewRequestWithContext(ctx, "POST", u.String(), bytes.NewReader(b))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", xerrors.Errorf("%s", string(body))
	}

	response := synthesizeResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}

	return response.AudioContent, nil
}
