package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/inhzus/go-berater/config"
	"github.com/inhzus/go-berater/models"
	"gopkg.in/go-playground/assert.v1"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

func getEngine() *gin.Engine {

	_, filename, _, _ := runtime.Caller(0)
	paths := strings.Split(filename, "/")
	workDir := strings.Join(paths[:len(paths)-4], "/")
	_ = os.Chdir(workDir)
	models.Setup()
	config.Setup()
	engine := gin.Default()
	ApplyRoutes(engine.Group("/"))
	return engine
}

func getTestToken() string {
	req := httptest.NewRequest("GET", "/api/test/token/test", nil)

	w := httptest.NewRecorder()
	r := getEngine()
	r.ServeHTTP(w, req)

	var response map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	return "Bearer " + response["token"]
}

func TestTestToken(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/test/token/test", nil)

	w := httptest.NewRecorder()
	r := getEngine()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, err, nil)
	_, exist := response["token"]
	assert.Equal(t, exist, true)
}

func TestCheckToken(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/token", nil)

	r := getEngine()

	q := httptest.NewRecorder()
	r.ServeHTTP(q, req)
	assert.Equal(t, http.StatusUnauthorized, q.Code)

	w := httptest.NewRecorder()
	req.Header.Set("Authorization", getTestToken())
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSendCode(t *testing.T) {
	r := getEngine()
	token := getTestToken()
	req := httptest.NewRequest("POST", "/api/code", nil)
	req.Header.Set("Authorization", token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusBadRequest)

	bodyString := `{"phone": "0"}`
	req = httptest.NewRequest("POST", "/api/code", strings.NewReader(bodyString))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusInternalServerError)
}

func TestCheckCode(t *testing.T) {
	client.Del("test_code")

	r := getEngine()
	genCode := "9999"
	req := httptest.NewRequest("GET", "/api/code/"+genCode, nil)
	req.Header.Set("Authorization", getTestToken())

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusNotFound)

	client.HMSet("test_code", map[string]interface{}{
		"code":   genCode,
		"phone":  "0",
		"status": false,
	})
	q := httptest.NewRecorder()
	r.ServeHTTP(q, req)
	assert.Equal(t, q.Code, http.StatusOK)

	ret, _ := client.HGet("test_code", "status").Result()

	status, _ := strconv.ParseBool(ret)
	assert.Equal(t, status, true)
	client.Del("test_code")
}

func TestAddCandidate(t *testing.T) {
	r := getEngine()
	redisKey := "test_code"
	client.Del(redisKey)
	client.HMSet(redisKey, map[string]interface{}{
		"code":   "9999",
		"phone":  "0",
		"status": false,
	})
	if models.ExistCandidateById("test") {
		models.RemoveCandidateById("test")
	}

	token := getTestToken()

	req := httptest.NewRequest("POST", "/api/candidate", nil)
	req.Header.Set("Authorization", token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, w.Code, http.StatusBadRequest)

	bodyBytes := `{"phone": "0","name": "inh","province": "sx","city": "123456","score": "675.5","subject": "li"}`

	req = httptest.NewRequest("POST", "/api/candidate", strings.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, w.Code, http.StatusUnauthorized)

	client.HSet(redisKey, "status", true)

	req = httptest.NewRequest("POST", "/api/candidate", strings.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, w.Code, http.StatusOK)

	req = httptest.NewRequest("POST", "/api/candidate", nil)
	req.Header.Set("Authorization", token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, w.Code, http.StatusConflict)

	client.Del(redisKey)
	models.RemoveCandidateById("test")
}

func TestUpdateCandidate(t *testing.T) {
	r := getEngine()

	redisKey := "test_code"
	client.Del(redisKey)
	models.RemoveCandidateById("test")

	token := getTestToken()
	req := httptest.NewRequest("PATCH", "/api/candidate", nil)
	req.Header.Set("Authorization", token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, w.Code, http.StatusNotFound)

	var added models.Candidate
	addedString := `{"openid":"test","phone": "0","name": "inh","province": "sx","city": "123456","score": "675.5","subject": "li"}`
	err := json.Unmarshal([]byte(addedString), &added)
	assert.Equal(t, err, nil)
	_ = models.AddCandidate(&added)

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, w.Code, http.StatusBadRequest)

	bodyString := `{"phone":"1","subject":"wen"}`
	req = httptest.NewRequest("PATCH", "/api/candidate", strings.NewReader(bodyString))
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, w.Code, http.StatusUnauthorized)

	bodyString = `{"subject":"wen"}`
	req = httptest.NewRequest("PATCH", "/api/candidate", strings.NewReader(bodyString))
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, w.Code, http.StatusOK)
	assert.Equal(t, models.GetCandidateById("test").Subject, "wen")

	client.HMSet(redisKey, map[string]interface{}{
		"status": false,
		"phone": "1",
		"code": "9999",
	})
	bodyString = `{"subject":"wex","phone":"1"}`
	req = httptest.NewRequest("PATCH", "/api/candidate", strings.NewReader(bodyString))
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, w.Code, http.StatusUnauthorized)

	client.HSet(redisKey, "status", true)
	bodyString = `{"subject":"wex","phone":"1"}`
	req = httptest.NewRequest("PATCH", "/api/candidate", strings.NewReader(bodyString))
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, w.Code, http.StatusOK)

	bodyString = `{"subject":"wex","phone":"2"}`
	req = httptest.NewRequest("PATCH", "/api/candidate", strings.NewReader(bodyString))
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, w.Code, http.StatusUnauthorized)


	models.RemoveCandidateById("openid")
	client.Del(redisKey)
}
