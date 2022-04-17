package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
	"github.com/kujilabo/cocotola-api/pkg_plugin/common/service"
)

var timeoutImportMin = 30

type tatoebaSentenceFindParameter struct {
	PageNo   int    `json:"pageNo" binding:"required,gte=1"`
	PageSize int    `json:"pageSize" binding:"required,gte=1"`
	Keyword  string `json:"keyword"`
	Random   bool   `json:"random"`
}

type tatoebaSentenceResponse struct {
	SentenceNumber int       `json:"sentenceNumber"`
	Lang           string    `json:"lang"`
	Text           string    `json:"text"`
	Author         string    `json:"author"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

func (s *tatoebaSentenceResponse) toModel() (service.TatoebaSentence, error) {
	lang, err := domain.NewLang3(s.Lang)
	if err != nil {
		return nil, err
	}

	return service.NewTatoebaSentence(s.SentenceNumber, lang, s.Text, s.Author, s.UpdatedAt)
}

type tatoebaSentencePair struct {
	Src tatoebaSentenceResponse `json:"src"`
	Dst tatoebaSentenceResponse `json:"dst"`
}

func (p *tatoebaSentencePair) toModel() (service.TatoebaSentencePair, error) {
	src, err := p.Src.toModel()
	if err != nil {
		return nil, err
	}

	dst, err := p.Dst.toModel()
	if err != nil {
		return nil, err
	}

	return service.NewTatoebaSentencePair(src, dst)
}

type tatoebaSentenceFindResponse struct {
	TotalCount int64                 `json:"totalCount"`
	Results    []tatoebaSentencePair `json:"results"`
}

func (r *tatoebaSentenceFindResponse) toModel() (*service.TatoebaSentencePairSearchResult, error) {
	sentences := make([]service.TatoebaSentencePair, len(r.Results))
	for i, r := range r.Results {
		pair, err := r.toModel()
		if err != nil {
			return nil, err
		}
		sentences[i] = pair
	}
	return &service.TatoebaSentencePairSearchResult{
		TotalCount: int64(len(r.Results)),
		Results:    sentences,
	}, nil
}

type tatoebaClient struct {
	endpoint     string
	username     string
	password     string
	client       http.Client
	importClient http.Client
}

func NewTatoebaClient(endpoint, username, password string, timeout time.Duration) service.TatoebaClient {
	return &tatoebaClient{
		endpoint: endpoint,
		username: username,
		password: password,
		client: http.Client{
			Timeout:   timeout,
			Transport: otelhttp.NewTransport(http.DefaultTransport),
		},
		importClient: http.Client{
			Timeout:   time.Minute * time.Duration(timeoutImportMin),
			Transport: otelhttp.NewTransport(http.DefaultTransport),
		},
	}
}

func (c *tatoebaClient) FindSentencePairs(ctx context.Context, param service.TatoebaSentenceSearchCondition) (*service.TatoebaSentencePairSearchResult, error) {
	ctx, span := tracer.Start(ctx, "tatoebaClient.FindSentencePairs")
	defer span.End()

	u, err := url.Parse(c.endpoint)
	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, "v1", "user", "sentence_pair", "find")

	params := tatoebaSentenceFindParameter{
		PageNo:   param.GetPageNo(),
		PageSize: param.GetPageSize(),
		Keyword:  param.GetKeyword(),
		Random:   param.IsRandom(),
	}

	paramBytes, err := json.Marshal(params)
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

	response := tatoebaSentenceFindResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.toModel()
}

func (c *tatoebaClient) FindSentenceBySentenceNumber(ctx context.Context, sentenceNumber int) (service.TatoebaSentence, error) {
	ctx, span := tracer.Start(ctx, "tatoebaClient.FindSentenceBySentenceNumber")
	defer span.End()

	u, err := url.Parse(c.endpoint)
	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, "v1", "user", "sentence", strconv.Itoa(sentenceNumber))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), http.NoBody)
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

	response := tatoebaSentenceResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.toModel()

}

func (c *tatoebaClient) ImportSentences(ctx context.Context, reader io.Reader) error {
	ctx, span := tracer.Start(ctx, "tatoebaClient.ImportSentences")
	defer span.End()

	u, err := url.Parse(c.endpoint)
	if err != nil {
		return err
	}

	u.Path = path.Join(u.Path, "v1", "admin", "sentence", "import")
	body := bytes.Buffer{}
	mw := multipart.NewWriter(&body)

	fw, err := mw.CreateFormFile("file", "filename")
	if err != nil {
		return err
	}

	if _, err := io.Copy(fw, reader); err != nil {
		return err
	}

	if err := mw.Close(); err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), &body)
	if err != nil {
		return err
	}

	req.SetBasicAuth(c.username, c.password)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	resp, err := c.importClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := c.errorHandle(resp.StatusCode); err != nil {
		return err
	}

	return nil
}

func (c *tatoebaClient) ImportLinks(ctx context.Context, reader io.Reader) error {
	ctx, span := tracer.Start(ctx, "tatoebaClient.ImportLinks")
	defer span.End()

	u, err := url.Parse(c.endpoint)
	if err != nil {
		return err
	}

	u.Path = path.Join(u.Path, "v1", "admin", "sentence", "import")
	body := bytes.Buffer{}
	mw := multipart.NewWriter(&body)

	fw, err := mw.CreateFormFile("file", "filename")
	if err != nil {
		return err
	}

	if _, err := io.Copy(fw, reader); err != nil {
		return err
	}

	if err := mw.Close(); err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), &body)
	if err != nil {
		return err
	}

	req.SetBasicAuth(c.username, c.password)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	resp, err := c.importClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := c.errorHandle(resp.StatusCode); err != nil {
		return err
	}

	return nil
}

func (c *tatoebaClient) errorHandle(statusCode int) error {
	if statusCode == http.StatusOK {
		return nil
	}

	return errors.New(http.StatusText(statusCode))
}
