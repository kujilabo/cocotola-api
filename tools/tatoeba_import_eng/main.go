package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"

	"github.com/kujilabo/cocotola-api/src/app/config"
	authG "github.com/kujilabo/cocotola-api/src/auth/gateway"
	userD "github.com/kujilabo/cocotola-api/src/user/domain"
)

var timeoutImportMin = 20

func main() {
	ctx := context.Background()
	cfg, err := config.LoadConfig("local")
	if err != nil {
		panic(err)
	}

	signingKey := []byte(cfg.Auth.SigningKey)
	signingMethod := jwt.SigningMethodHS256
	authTokenManager := authG.NewAuthTokenManager(signingKey, signingMethod, time.Duration(cfg.Auth.AccessTokenTTLMin)*time.Minute, time.Duration(cfg.Auth.RefreshTokenTTLHour)*time.Hour)

	modelU, err := userD.NewModel(1, 1, time.Now(), time.Now(), 1, 1)
	if err != nil {
		panic(err)
	}
	appUserModel, err := userD.NewAppUserModel(modelU, userD.OrganizationID(1), "test", "Test", []string{"Owner"}, map[string]string{})
	if err != nil {
		panic(err)
	}

	modelO, err := userD.NewModel(1, 1, time.Now(), time.Now(), 1, 1)
	if err != nil {
		panic(err)
	}
	organizationModel, err := userD.NewOrganizationModel(modelO, "Test")
	if err != nil {
		panic(err)
	}
	tokenSet, err := authTokenManager.CreateTokenSet(ctx, appUserModel, organizationModel)
	if err != nil {
		panic(err)
	}

	url := "http://localhost:8080/plugin/tatoeba/sentence/import"
	fieldname := "file"
	filename := "eng_sentences_detailed.tsv"

	file, err := os.Open("../cocotola-data/datasource/tatoeba/" + filename)
	if err != nil {
		panic(err)
	}

	body := bytes.Buffer{}

	mw := multipart.NewWriter(&body)

	fw, err := mw.CreateFormFile(fieldname, filename)
	if err != nil {
		panic(err)
	}

	if _, err := io.Copy(fw, file); err != nil {
		panic(err)
	}

	if err = mw.Close(); err != nil {
		panic(err)
	}

	req, err := http.NewRequest(http.MethodPost, url, &body)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", "Bearer "+tokenSet.AccessToken)
	req.Header.Set("Content-Type", mw.FormDataContentType())

	client := http.Client{
		Timeout: time.Duration(timeoutImportMin) * time.Minute,
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Printf("status: %d\n", resp.StatusCode)
	fmt.Printf("body: %s\n", string(respBody))
}
