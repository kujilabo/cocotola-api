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

func Test_FindSentencePairs_OK(t *testing.T) {
	// given
	tatoebaClient := new(service_mock.TatoebaClient)
	// -
	src := new(service_mock.TatoebaSentence)
	dst := new(service_mock.TatoebaSentence)
	src.On("GetSentenceNumber").Return(1)
	src.On("GetLang2").Return(appD.Lang2EN)
	src.On("GetText").Return("test1")
	src.On("GetAuthor").Return("author1")
	src.On("GetUpdatedAt").Return(time.Now())
	dst.On("GetSentenceNumber").Return(2)
	dst.On("GetLang2").Return(appD.Lang2JA)
	dst.On("GetText").Return("test2")
	dst.On("GetAuthor").Return("author2")
	dst.On("GetUpdatedAt").Return(time.Now())
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
	// - keyword: apple
	body, err := json.Marshal(gin.H{"keyword": "apple", "pageNo": 1, "pageSize": 10})
	require.NoError(t, err)
	req, err := http.NewRequest(http.MethodPost, "/v1/plugin/tatoeba/find", bytes.NewBuffer(body))
	req.SetBasicAuth("user", "pass")
	require.NoError(t, err)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// then
	resultsExpr := parseExpr(t, "$.results[*]")
	totalCountExpr := parseExpr(t, "$.totalCount")

	// bytes, _ := io.ReadAll(w.Body)
	// fmt.Println(string(bytes))
	// t.Logf("resp: %s", string(bytes))
	// - check the status code
	assert.Equal(t, http.StatusOK, w.Code)
	jsonObj := parseJSON(t, w.Body)

	jsonResults := resultsExpr.Get(jsonObj)
	assert.Equal(t, 1, len(jsonResults))

	jsonTotalCount := totalCountExpr.Get(jsonObj)
	assert.Equal(t, 1, int(jsonTotalCount[0].(int64)))
}
