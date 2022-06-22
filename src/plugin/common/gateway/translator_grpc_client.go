package gateway

import (
	"context"
	"time"

	"google.golang.org/grpc"

	appD "github.com/kujilabo/cocotola-api/src/app/domain"
	"github.com/kujilabo/cocotola-api/src/plugin/common/domain"
	"github.com/kujilabo/cocotola-api/src/plugin/common/service"
	pb "github.com/kujilabo/cocotola-api/src/proto"
)

type translatorGRPCClient struct {
	userClient  pb.TranslatorUserClient
	adminClient pb.TranslatorAdminClient
}

func NewTranslatorGRPCClient(conn *grpc.ClientConn, timeout time.Duration) service.TranslatorClient {
	userClient := pb.NewTranslatorUserClient(conn)
	adminClient := pb.NewTranslatorAdminClient(conn)
	return &translatorGRPCClient{
		userClient:  userClient,
		adminClient: adminClient,
	}
}

func (c *translatorGRPCClient) DictionaryLookup(ctx context.Context, fromLang, toLang appD.Lang2, text string) ([]domain.Translation, error) {
	cancelContext, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	param := pb.DictionaryLookupParameter{
		FromLang2: fromLang.String(),
		ToLang2:   toLang.String(),
		Text:      text,
	}

	resp, err := c.userClient.DictionaryLookup(cancelContext, &param)
	if err != nil {
		return nil, err
	}

	translationList := make([]domain.Translation, len(resp.Results))
	for i, r := range resp.Results {
		pos, err := domain.NewWordPos(int(r.Pos))
		if err != nil {
			return nil, err
		}

		lang2, err := appD.NewLang2(r.Lang2)
		if err != nil {
			return nil, err
		}

		m, err := domain.NewTranslation(r.Text, pos, lang2, r.Translated, r.Provider)
		if err != nil {
			return nil, err
		}

		translationList[i] = m
	}

	return translationList, nil
}

func (c *translatorGRPCClient) DictionaryLookupWithPos(ctx context.Context, fromLang, toLang appD.Lang2, text string, pos domain.WordPos) (domain.Translation, error) {
	return nil, nil
}
func (c *translatorGRPCClient) FindTranslationsByFirstLetter(ctx context.Context, lang2 appD.Lang2, firstLetter string) ([]domain.Translation, error) {
	return nil, nil
}
func (c *translatorGRPCClient) FindTranslationByTextAndPos(ctx context.Context, lang2 appD.Lang2, text string, pos domain.WordPos) (domain.Translation, error) {
	return nil, nil
}
func (c *translatorGRPCClient) FindTranslationsByText(ctx context.Context, lang2 appD.Lang2, text string) ([]domain.Translation, error) {
	return nil, nil
}
func (c *translatorGRPCClient) AddTranslation(ctx context.Context, param service.TranslationAddParameter) error {
	return nil
}
func (c *translatorGRPCClient) UpdateTranslation(ctx context.Context, lang2 appD.Lang2, text string, pos domain.WordPos, param service.TranslationUpdateParameter) error {
	return nil
}
func (c *translatorGRPCClient) RemoveTranslation(ctx context.Context, lang2 appD.Lang2, text string, pos domain.WordPos) error {
	return nil
}
