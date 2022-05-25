package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/overlorddamygod/go-auth/configs"
	"github.com/overlorddamygod/go-auth/utils"
)

type GithubTokenResponse struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

type GithubUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type GithubOauth struct {
	clientId     string
	clientSecret string
}

var (
	githubOauthUrl    string = "https://github.com/login/oauth/authorize"
	githubTokenUrl    string = "https://github.com/login/oauth/access_token?client_id=%s&code=%s&client_secret=%s&redirect_uri=%s"
	githubUserInfoUrl string = "https://api.github.com/user"
)

func NewGithubOauth() (*GithubOauth, error) {
	if !IsProviderValid(configs.MainConfig, "github") {
		return nil, errors.New("invalid oauth provider")
	}

	return &GithubOauth{
		clientId:     configs.MainConfig.Oauth["github"].ClientID,
		clientSecret: configs.MainConfig.Oauth["github"].ClientSecret,
	}, nil
}

func (g *GithubOauth) GetOauthUrl(redirect_to string) (string, error) {
	gihubConfig := configs.MainConfig.Oauth["github"]

	redirect_uri := fmt.Sprintf("%s/api/v1/auth/authorize?redirect_to=%s&oauth_provider=github", configs.MainConfig.ApiUrl, redirect_to)

	return utils.EncodeUrl(githubOauthUrl, map[string]string{
		"client_id":    gihubConfig.ClientID,
		"redirect_uri": redirect_uri,
	})
}

func (g *GithubOauth) Do(code string) (*GithubTokenResponse, *GithubUser, error) {
	tokenRes, err := g.getAccessToken(code)

	if err != nil {
		return nil, nil, err
	}

	user, err := g.getUserInfo(tokenRes.AccessToken)

	if err != nil {
		return nil, nil, err
	}

	return tokenRes, user, nil
}

func (g *GithubOauth) getAccessToken(code string) (*GithubTokenResponse, error) {
	url, err := utils.EncodeUrl(githubTokenUrl, map[string]string{
		"client_id":     g.clientId,
		"code":          code,
		"client_secret": g.clientSecret,
		"redirect_uri":  configs.MainConfig.ApiUrl,
	})

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, nil)
	req.Header.Set("Accept", "application/json")
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var tokenRes GithubTokenResponse

	json.NewDecoder(resp.Body).Decode(&tokenRes)
	return &tokenRes, nil
}

func (g *GithubOauth) getUserInfo(accessToken string) (*GithubUser, error) {
	req, err := http.NewRequest("GET", githubUserInfoUrl, nil)
	req.Header.Set("Authorization", "bearer "+accessToken)

	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var user GithubUser

	json.NewDecoder(resp.Body).Decode(&user)

	return &user, nil
}
