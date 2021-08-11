package gateway

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/cognitiveservices/v3.0/translatortext"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/to"

	app "github.com/kujilabo/cocotola-api/pkg_app/domain"
	"github.com/kujilabo/cocotola-api/pkg_plugin/common/domain"
)

type azureTranslatorClient struct {
	client translatortext.TranslatorClient
}

type AzureDisplayTranslation struct {
	Pos        int
	Target     string
	Confidence float64
}

func NewAzureTranslatorClient(subscriptionKey string) domain.Translator {
	client := translatortext.NewTranslatorClient("https://api.cognitive.microsofttranslator.com")
	client.Authorizer = autorest.NewCognitiveServicesAuthorizer(subscriptionKey)
	return &azureTranslatorClient{
		client: client,
	}

}

func (c *azureTranslatorClient) DictionaryLookup(ctx context.Context, text string, fromLang, toLang app.Lang2) ([]domain.TranslationResult, error) {
	result, err := c.client.DictionaryLookup(context.Background(), fromLang.String(), toLang.String(), []translatortext.DictionaryLookupTextInput{{Text: to.StringPtr(text)}}, "")
	if err != nil {
		return nil, err
	}
	if result.Value == nil {
		return nil, nil
	}

	translations := make([]domain.TranslationResult, 0)
	for _, v := range *result.Value {
		if v.Translations == nil {
			continue
		}

		for _, t := range *v.Translations {
			pos, err := domain.ParsePos(c.pointerToString(t.PosTag))
			if err != nil {
				return nil, err
			}
			translations = append(translations, domain.TranslationResult{
				Pos:        pos,
				Target:     c.pointerToString(t.DisplayTarget),
				Confidence: c.pointerToFloat64(t.Confidence),
			})
		}
	}
	return translations, nil
}

func (c *azureTranslatorClient) pointerToString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func (c *azureTranslatorClient) pointerToFloat64(value *float64) float64 {
	if value == nil {
		return 0
	}
	return *value
}

// func (c *azureTranslatorClient) stringToPos(value string) int {
// 	switch value {
// 	case "ADJ":
// 		return 1
// 	case "ADV":
// 		return 2
// 	case "CONJ":
// 		return 3
// 	case "DET":
// 		return 4
// 	case "MODAL":
// 		return 5
// 	case "NOUN":
// 		return 6
// 	case "PREP":
// 		return 7
// 	case "PRON":
// 		return 8
// 	case "VERB":
// 		return 9
// 	default:
// 		return 99
// 	}
// }
