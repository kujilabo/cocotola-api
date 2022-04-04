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

	app "github.com/kujilabo/cocotola-api/pkg_app/domain"
	"github.com/kujilabo/cocotola-api/pkg_plugin/common/domain"
	"github.com/kujilabo/cocotola-api/pkg_plugin/common/service"
)

type translationClient struct {
	endpoint string
	client   http.Client
}

func NewTranslationClient(endpoint string, timeout time.Duration) (service.TranslationClient, error) {
	return &translationClient{
		client: http.Client{
			Timeout: timeout,
		},
	}, nil
}

func (c *translationClient) DictionaryLookup(ctx context.Context, fromLang, toLang app.Lang2, text string) ([]domain.Translation, error) {
	u, err := url.Parse(c.endpoint)
	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, "dictionary", "lookup")
	q := u.Query()
	q.Set("text", text)
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := c.errorHandle(resp.StatusCode); err != nil {
		return nil, err
	}

	response := []domain.Translation{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response, nil
}

func (c *translationClient) DictionaryLookupWithPos(ctx context.Context, fromLang, toLang app.Lang2, text string, pos domain.WordPos) (domain.Translation, error) {
	u, err := url.Parse(c.endpoint)
	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, "dictionary", "lookup")
	q := u.Query()
	q.Set("text", text)
	q.Set("pos", strconv.Itoa(int(pos)))
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := c.errorHandle(resp.StatusCode); err != nil {
		return nil, err
	}

	response := []domain.Translation{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	if len(response) == 0 {
		return nil, service.ErrTranslationNotFound
	}

	return response[0], nil
}

func (c *translationClient) FindTranslationsByFirstLetter(ctx context.Context, lang app.Lang2, firstLetter string) ([]domain.Translation, error) {
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

	req, err := http.NewRequest(http.MethodGet, u.String(), bytes.NewBuffer(paramBytes))
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := c.errorHandle(resp.StatusCode); err != nil {
		return nil, err
	}

	response := []domain.Translation{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response, nil
}

func (c *translationClient) FindTranslationByTextAndPos(ctx context.Context, lang app.Lang2, text string, pos domain.WordPos) (domain.Translation, error) {
	u, err := url.Parse(c.endpoint)
	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, "text", text, "pos", strconv.Itoa(int(pos)))

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := c.errorHandle(resp.StatusCode); err != nil {
		return nil, err
	}

	response := []domain.Translation{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	if len(response) == 0 {
		return nil, service.ErrTranslationNotFound
	}

	return response[0], nil
}

func (c *translationClient) FindTranslationsByText(ctx context.Context, lang app.Lang2, text string) ([]domain.Translation, error) {
	u, err := url.Parse(c.endpoint)
	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, "text", text)

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := c.errorHandle(resp.StatusCode); err != nil {
		return nil, err
	}

	response := []domain.Translation{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response, nil
}

func (c *translationClient) errorHandle(statusCode int) error {
	if statusCode == http.StatusOK {
		return nil
	}

	return errors.New(http.StatusText(statusCode))
}

func (c *translationClient) AddTranslation(ctx context.Context, param service.TranslationAddParameter) error {
	u, err := url.Parse(c.endpoint)
	if err != nil {
		return err
	}

	paramMap := map[string]interface{}{
		"lang":       param.GetLang(),
		"text":       param.GetText(),
		"pos":        param.GetPos(),
		"translated": param.GetTranslated(),
	}

	paramBytes, err := json.Marshal(paramMap)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(paramBytes))
	if err != nil {
		return err
	}

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

func (c *translationClient) UpdateTranslation(ctx context.Context, lang app.Lang2, text string, pos domain.WordPos, param service.TranslationUpdateParameter) error {
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

	req, err := http.NewRequest(http.MethodPut, u.String(), bytes.NewBuffer(paramBytes))
	if err != nil {
		return err
	}

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

func (c *translationClient) RemoveTranslation(ctx context.Context, lang app.Lang2, text string, pos domain.WordPos) error {
	u, err := url.Parse(c.endpoint)
	if err != nil {
		return err
	}

	u.Path = path.Join(u.Path, "text", text, "pos", strconv.Itoa(int(pos)))

	req, err := http.NewRequest(http.MethodDelete, u.String(), nil)
	if err != nil {
		return err
	}

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
