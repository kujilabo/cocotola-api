package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/xerrors"

	"github.com/kujilabo/cocotola-api/src/auth/service"
	"github.com/kujilabo/cocotola-api/src/lib/log"
)

type googleAuthClient struct {
	client       http.Client
	clientID     string
	clientSecret string
	redirectURI  string
	grantType    string
}

func NewGoogleAuthClient(clientID, clientSecret, redirectURI string, timeout time.Duration) service.GoogleAuthClient {
	return &googleAuthClient{
		client: http.Client{
			Timeout:   timeout,
			Transport: otelhttp.NewTransport(http.DefaultTransport),
		},
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURI:  redirectURI,
		grantType:    "authorization_code",
	}

}

func (c *googleAuthClient) RetrieveAccessToken(ctx context.Context, code string) (*service.GoogleAuthResponse, error) {
	ctx, span := tracer.Start(ctx, "googleAuthClient.RetrieveAccessToken")
	defer span.End()

	logger := log.FromContext(ctx)

	paramMap := map[string]string{
		"client_id":     c.clientID,
		"client_secret": c.clientSecret,
		"redirect_uri":  c.redirectURI,
		"grant_type":    c.grantType,
		"code":          code,
	}

	paramBytes, err := json.Marshal(paramMap)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://accounts.google.com/o/oauth2/token", bytes.NewBuffer(paramBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, xerrors.Errorf("failed to retrieve access token.err: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		logger.Debugf("status:%d", resp.StatusCode)
		logger.Debugf("Resp:%s", string(respBytes))
		return nil, errors.New(string(respBytes))
	}

	googleAuthResponse := service.GoogleAuthResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&googleAuthResponse); err != nil {
		return nil, err
	}
	logger.Infof("RetrieveAccessToken:%s", googleAuthResponse.AccessToken)

	return &googleAuthResponse, nil
}

func (c *googleAuthClient) RetrieveUserInfo(ctx context.Context, googleAuthResponse *service.GoogleAuthResponse) (*service.GoogleUserInfo, error) {
	ctx, span := tracer.Start(ctx, "googleAuthClient.RetrieveUserInfo")
	defer span.End()

	logger := log.FromContext(ctx)

	req, err := http.NewRequestWithContext(ctx, "GET", "https://www.googleapis.com/oauth2/v1/userinfo", http.NoBody)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("alt", "json")
	q.Add("access_token", googleAuthResponse.AccessToken)
	req.URL.RawQuery = q.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	logger.Debugf("access_token:%s", googleAuthResponse.AccessToken)
	logger.Debugf("status:%d", resp.StatusCode)

	googleUserInfo := service.GoogleUserInfo{}
	if err := json.NewDecoder(resp.Body).Decode(&googleUserInfo); err != nil {
		return nil, err
	}

	return &googleUserInfo, nil
}
