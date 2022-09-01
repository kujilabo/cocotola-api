package controller_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/oj"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/kujilabo/cocotola-api/src/app/controller"
	appD "github.com/kujilabo/cocotola-api/src/app/domain"
	"github.com/kujilabo/cocotola-api/src/plugin/common/service"
	service_mock "github.com/kujilabo/cocotola-api/src/plugin/common/service/mock"
)

var anythingOfContext = mock.MatchedBy(func(_ context.Context) bool { return true })
var anythingOfTatoebaSentenceSearchCondition = mock.MatchedBy(func(_ service.TatoebaSentenceSearchCondition) bool { return true })

func mapGetString(m map[string]interface{}, key string) string {
	return m[key].(string)

}

func mapGetInt(m map[string]interface{}, key string) int {
	return interfaceInt64ToInt(m[key])
}

func interfaceInt64ToInt(i interface{}) int {
	return int(i.(int64))
}

func parseJSON(t *testing.T, b *bytes.Buffer) interface{} {
	respBytes, err := io.ReadAll(b)
	require.NoError(t, err)
	obj, err := oj.Parse(respBytes)
	require.NoError(t, err)
	return obj
}

func parseExpr(t *testing.T, v string) jp.Expr {
	expr, err := jp.ParseString(v)
	require.NoError(t, err)
	return expr
}

func initTatoebaRouter(tatoebaClient service.TatoebaClient) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	v1 := router.Group("v1")
	plugin := v1.Group("plugin")
	controller.InitTatoebaPluginRouter(plugin, tatoebaClient)
	return router
}

func testNewTatoebaSentence(sentenceaNumber int, lang2 appD.Lang2, text string, author string) service.TatoebaSentence {
	sentence := new(service_mock.TatoebaSentence)
	sentence.On("GetSentenceNumber").Return(sentenceaNumber)
	sentence.On("GetLang2").Return(lang2)
	sentence.On("GetText").Return(text)
	sentence.On("GetAuthor").Return(author)
	sentence.On("GetUpdatedAt").Return(time.Now())
	return sentence
}

func Test_FindSentencePairs_OK(t *testing.T) {
	// given
	tatoebaClient := new(service_mock.TatoebaClient)
	// -
	src := testNewTatoebaSentence(1, appD.Lang2EN, "test1", "author1")
	dst := testNewTatoebaSentence(2, appD.Lang2EN, "test2", "author2")
	pair := new(service_mock.TatoebaSentencePair)
	pair.On("GetSrc").Return(src)
	pair.On("GetDst").Return(dst)
	results := &service.TatoebaSentencePairSearchResult{
		TotalCount: 1,
		Results:    []service.TatoebaSentencePair{pair},
	}
	tatoebaClient.On("FindSentencePairs", anythingOfContext, anythingOfTatoebaSentenceSearchCondition).Return(results, nil)

	r := initTatoebaRouter(tatoebaClient)

	// when
	// parameter is valid
	// - keyword: apple
	// - pageNo: 1
	// - pageSize: 1
	body, err := json.Marshal(gin.H{"keyword": "apple", "pageNo": 1, "pageSize": 10})
	require.NoError(t, err)
	req, err := http.NewRequest(http.MethodPost, "/v1/plugin/tatoeba/find", bytes.NewBuffer(body))
	require.NoError(t, err)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// then
	resultsExpr := parseExpr(t, "$.results[*]")
	totalCountExpr := parseExpr(t, "$.totalCount")

	// bytes, _ := io.ReadAll(w.Body)
	// t.Logf("resp: %s", string(bytes))
	// - check the status code
	assert.Equal(t, http.StatusOK, w.Code)
	jsonObj := parseJSON(t, w.Body)

	{
		results := resultsExpr.Get(jsonObj)
		assert.Equal(t, 1, len(results))

		results0 := results[0].(map[string]interface{})
		results0src := results0["src"].(map[string]interface{})
		assert.Equal(t, mapGetInt(results0src, "sentenceNumber"), 1)
		assert.Equal(t, mapGetString(results0src, "lang2"), "en")
		assert.Equal(t, mapGetString(results0src, "text"), "test1")
		assert.Equal(t, mapGetString(results0src, "author"), "author1")

		totalCount := totalCountExpr.Get(jsonObj)
		assert.Equal(t, 1, interfaceInt64ToInt(totalCount[0]))
	}
}
