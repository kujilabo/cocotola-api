package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/kujilabo/cocotola-api/src/app/domain"
	"github.com/kujilabo/cocotola-api/src/app/service"
)

type synthesizerClient struct {
	endpoint string
	username string
	password string
	client   http.Client
}

type audioResponse struct {
	ID      int    `json:"id"`
	Lang2   string `json:"lang2"`
	Text    string `json:"text"`
	Content string `json:"content"`
}

func (r *audioResponse) toModel() (service.Audio, error) {
	lang2, err := domain.NewLang2(r.Lang2)
	if err != nil {
		return nil, err
	}

	audioModel, err := domain.NewAudioModel(uint(r.ID), lang2, r.Text, r.Content)
	if err != nil {
		return nil, err
	}

	return service.NewAudio(audioModel)
}

func NewSynthesizerClient(endpoint, username, password string, timeout time.Duration) service.SynthesizerClient {
	return &synthesizerClient{
		endpoint: endpoint,
		username: username,
		password: password,
		client: http.Client{
			Timeout:   timeout,
			Transport: otelhttp.NewTransport(http.DefaultTransport),
		},
	}
}

func (c *synthesizerClient) Synthesize(ctx context.Context, lang2 domain.Lang2, text string) (service.Audio, error) {
	ctx, span := tracer.Start(ctx, "synthesizerClient.Synthesize")
	defer span.End()

	u, err := url.Parse(c.endpoint)
	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, "v1", "user", "synthesize")
	paramMap := map[string]string{
		"lang2": lang2.String(),
		"text":  text,
	}

	paramBytes, err := json.Marshal(paramMap)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewBuffer(paramBytes))
	if err != nil {
		return nil, err

	}

	req.SetBasicAuth(c.username, c.password)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := c.errorHandle(resp.StatusCode); err != nil {
		return nil, err
	}

	response := audioResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.toModel()
}

func (c *synthesizerClient) FindAudioByAudioID(ctx context.Context, audioID domain.AudioID) (service.Audio, error) {
	ctx, span := tracer.Start(ctx, "synthesizerClient.FindAudioByAudioID")
	defer span.End()

	u, err := url.Parse(c.endpoint)
	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, "v1", "user", "audio", strconv.Itoa(int(audioID)))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), http.NoBody)
	if err != nil {
		return nil, err

	}

	req.SetBasicAuth(c.username, c.password)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := c.errorHandle(resp.StatusCode); err != nil {
		return nil, err
	}

	response := audioResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.toModel()
}

func (c *synthesizerClient) errorHandle(statusCode int) error {
	if statusCode == http.StatusOK {
		return nil
	} else if statusCode == http.StatusNotFound {
		return service.ErrAudioNotFound
	}

	return errors.New(http.StatusText(statusCode))
}
