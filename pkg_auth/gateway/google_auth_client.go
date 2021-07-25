package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/kujilabo/cocotola-api/pkg_auth/domain"
	"github.com/kujilabo/cocotola-api/pkg_lib/log"
	"golang.org/x/xerrors"
)

type googleAuthClient struct {
	clientID     string
	clientSecret string
	redirectURI  string
	grantType    string
}

func NewGoogleAuthClient(clientID, clientSecret, redirectURI string) domain.GoogleAuthClient {
	return &googleAuthClient{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURI:  redirectURI,
		grantType:    "authorization_code",
	}

}

func (c *googleAuthClient) RetrieveAccessToken(ctx context.Context, code string) (*domain.GoogleAuthResponse, error) {
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

	req, err := http.NewRequest("POST", "https://accounts.google.com/o/oauth2/token", bytes.NewBuffer(paramBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, xerrors.Errorf("failed to retrieve access token.err: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		respBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		logger.Infof("%s", c.clientSecret)
		logger.Infof("%s", c.clientID)
		logger.Infof("%s", c.redirectURI)
		logger.Infof("%s", c.grantType)
		logger.Infof("%s", code)
		logger.Infof("status:%d", resp.StatusCode)
		logger.Infof("Resp:%s", string(respBytes))
		return nil, xerrors.New(string(respBytes))
	}

	googleAuthResponse := domain.GoogleAuthResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&googleAuthResponse); err != nil {
		return nil, err
	}
	logger.Infof("RetrieveAccessToken:%s", googleAuthResponse.AccessToken)

	return &googleAuthResponse, nil
}

func (c *googleAuthClient) RetrieveUserInfo(ctx context.Context, googleAuthResponse *domain.GoogleAuthResponse) (*domain.GoogleUserInfo, error) {
	logger := log.FromContext(ctx)

	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v1/userinfo", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("alt", "json")
	q.Add("access_token", googleAuthResponse.AccessToken)
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	logger.Infof("access_token:%s", googleAuthResponse.AccessToken)
	logger.Infof("status:%d", resp.StatusCode)

	googleUserInfo := domain.GoogleUserInfo{}
	if err := json.NewDecoder(resp.Body).Decode(&googleUserInfo); err != nil {
		return nil, err
	}

	return &googleUserInfo, nil
}
