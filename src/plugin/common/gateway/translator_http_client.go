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

	appD "github.com/kujilabo/cocotola-api/src/app/domain"
	"github.com/kujilabo/cocotola-api/src/plugin/common/domain"
	"github.com/kujilabo/cocotola-api/src/plugin/common/service"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type translationResponse struct {
	Text       string `json:"text"`
	Pos        int    `json:"pos"`
	Lang2      string `json:"lang2"`
	Translated string `json:"translated"`
	Provider   string `json:"provider"`
}

func (r *translationResponse) toModel() (domain.Translation, error) {
	pos, err := domain.NewWordPos(r.Pos)
	if err != nil {
		return nil, err
	}

	lang2, err := appD.NewLang2(r.Lang2)
	if err != nil {
		return nil, err
	}

	return domain.NewTranslation(r.Text, pos, lang2, r.Translated, r.Provider)
}

type translationFindResponse struct {
	Results []translationResponse `json:"results"`
}

func (r *translationFindResponse) toModel() ([]domain.Translation, error) {
	translationList := make([]domain.Translation, len(r.Results))
	for i, r := range r.Results {
		m, err := r.toModel()
		if err != nil {
			return nil, err
		}
		translationList[i] = m
	}

	return translationList, nil
}

type translatorHTTPClient struct {
	endpoint string
	username string
	password string
	client   http.Client
}

func NewTranslatorHTTPClient(endpoint, username, password string, timeout time.Duration) service.TranslatorClient {
	return &translatorHTTPClient{
		endpoint: endpoint,
		username: username,
		password: password,
		client: http.Client{
			Timeout:   timeout,
			Transport: otelhttp.NewTransport(http.DefaultTransport),
		},
	}
}

func (c *translatorHTTPClient) DictionaryLookup(ctx context.Context, fromLang, toLang appD.Lang2, text string) ([]domain.Translation, error) {
	ctx, span := tracer.Start(ctx, "translatorClient.DictionaryLookup")
	defer span.End()

	u, err := url.Parse(c.endpoint)
	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, "v1", "user", "dictionary", "lookup")
	q := u.Query()
	q.Set("text", text)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
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

	response := translationFindResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.toModel()
}

func (c *translatorHTTPClient) DictionaryLookupWithPos(ctx context.Context, fromLang, toLang appD.Lang2, text string, pos domain.WordPos) (domain.Translation, error) {
	ctx, span := tracer.Start(ctx, "translatorClient.DictionaryLookupWithPos")
	defer span.End()

	u, err := url.Parse(c.endpoint)
	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, "v1", "user", "dictionary", "lookup")
	q := u.Query()
	q.Set("text", text)
	q.Set("pos", strconv.Itoa(int(pos)))
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
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

	response := translationResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.toModel()
}

func (c *translatorHTTPClient) FindTranslationsByFirstLetter(ctx context.Context, lang2 appD.Lang2, firstLetter string) ([]domain.Translation, error) {
	ctx, span := tracer.Start(ctx, "translatorClient.FindTranslationsByFirstLetter")
	defer span.End()

	u, err := url.Parse(c.endpoint)
	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, "find")
	paramMap := map[string]string{
		"letter": firstLetter,
	}

	paramBytes, err := json.Marshal(paramMap)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), bytes.NewBuffer(paramBytes))
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

	response := translationFindResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.toModel()
}

func (c *translatorHTTPClient) FindTranslationByTextAndPos(ctx context.Context, lang2 appD.Lang2, text string, pos domain.WordPos) (domain.Translation, error) {
	ctx, span := tracer.Start(ctx, "translatorClient.FindTranslationByTextAndPos")
	defer span.End()

	u, err := url.Parse(c.endpoint)
	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, "text", text, "pos", strconv.Itoa(int(pos)))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
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

	response := translationResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.toModel()
}

func (c *translatorHTTPClient) FindTranslationsByText(ctx context.Context, lang2 appD.Lang2, text string) ([]domain.Translation, error) {
	ctx, span := tracer.Start(ctx, "translatorClient.FindTranslationsByText")
	defer span.End()

	u, err := url.Parse(c.endpoint)
	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, "text", text)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
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

	response := translationFindResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.toModel()
}

func (c *translatorHTTPClient) AddTranslation(ctx context.Context, param service.TranslationAddParameter) error {
	ctx, span := tracer.Start(ctx, "translatorClient.AddTranslation")
	defer span.End()

	u, err := url.Parse(c.endpoint)
	if err != nil {
		return err
	}

	paramMap := map[string]interface{}{
		"lang2":      param.GetLang2(),
		"text":       param.GetText(),
		"pos":        param.GetPos(),
		"translated": param.GetTranslated(),
	}

	paramBytes, err := json.Marshal(paramMap)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewBuffer(paramBytes))
	if err != nil {
		return err
	}

	req.SetBasicAuth(c.username, c.password)
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := c.errorHandle(resp.StatusCode); err != nil {
		return err
	}

	return nil
}

func (c *translatorHTTPClient) UpdateTranslation(ctx context.Context, lang2 appD.Lang2, text string, pos domain.WordPos, param service.TranslationUpdateParameter) error {
	ctx, span := tracer.Start(ctx, "translatorClient.UpdateTranslation")
	defer span.End()

	u, err := url.Parse(c.endpoint)
	if err != nil {
		return err
	}

	u.Path = path.Join(u.Path, "text", text, "pos", strconv.Itoa(int(pos)))
	paramMap := map[string]interface{}{
		"translated": param.GetTranslated(),
	}

	paramBytes, err := json.Marshal(paramMap)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, u.String(), bytes.NewBuffer(paramBytes))
	if err != nil {
		return err
	}

	req.SetBasicAuth(c.username, c.password)
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := c.errorHandle(resp.StatusCode); err != nil {
		return err
	}

	return nil
}

func (c *translatorHTTPClient) RemoveTranslation(ctx context.Context, lang2 appD.Lang2, text string, pos domain.WordPos) error {
	ctx, span := tracer.Start(ctx, "translatorClient.RemoveTranslation")
	defer span.End()

	u, err := url.Parse(c.endpoint)
	if err != nil {
		return err
	}

	u.Path = path.Join(u.Path, "text", text, "pos", strconv.Itoa(int(pos)))

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, u.String(), nil)
	if err != nil {
		return err
	}

	req.SetBasicAuth(c.username, c.password)
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := c.errorHandle(resp.StatusCode); err != nil {
		return err
	}

	return nil
}

func (c *translatorHTTPClient) errorHandle(statusCode int) error {
	if statusCode == http.StatusOK {
		return nil
	}

	return errors.New(http.StatusText(statusCode))
}
