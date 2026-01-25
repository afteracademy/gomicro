package main

import (
	"io"
	"math/rand"
	"net/http"
	"regexp"

	"github.com/Kong/go-pdk"
	"github.com/Kong/go-pdk/server"
)

var healthRegex = regexp.MustCompile(`^/[^/]+/health$`)

type Config struct {
	ApiKeyVerificationURLs []string `json:"verification_urls"`
}

func New() any {
	return &Config{}
}

func (conf *Config) Access(kong *pdk.PDK) {

	// Skip health check paths
	path, err := kong.Request.GetPath()
	if err == nil && healthRegex.MatchString(path) {
		return
	}

	headers := map[string][]string{"Content-Type": {"application/json"}}

	apiKey, err := kong.Request.GetHeader("x-api-key")
	if err != nil {
		kong.Response.Exit(401, []byte(err.Error()), headers)
		return
	}

	if apiKey == "" {
		kong.Response.Exit(401, []byte(`{"code":"10001","status":"401","message":"x-api-key header is missing"}`), headers)
		return
	}

	var verificationURL string
	switch len(conf.ApiKeyVerificationURLs) {
	case 0:
		kong.Log.Err("verification_urls is missing in apikey_auth_plugin")
		kong.Response.Exit(500, []byte(`{"code":"10001","status":"500","message":"something went wrong"}`), headers)
		return
	case 1:
		verificationURL = conf.ApiKeyVerificationURLs[0]
	default:
		randomIndex := rand.Intn(len(conf.ApiKeyVerificationURLs))
		verificationURL = conf.ApiKeyVerificationURLs[randomIndex]
	}

	req, err := http.NewRequest("GET", verificationURL, nil)
	if err != nil {
		kong.Log.Err(err.Error())
		kong.Response.Exit(500, []byte(err.Error()), headers)
		return
	}

	req.Header.Set("x-api-key", apiKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		kong.Log.Err(err.Error())
		kong.Response.Exit(500, []byte(err.Error()), headers)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		kong.Log.Err(string(body))
		kong.Response.Exit(resp.StatusCode, body, headers)
		return
	}
}

func main() {
	server.StartServer(New, "0.1", 1000)
}
