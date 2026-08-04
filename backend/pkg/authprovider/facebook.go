package authprovider

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const facebookGraphBaseURL = "https://graph.facebook.com/v19.0"

type FacebookProvider struct {
	appID     string
	appSecret string
	client    *http.Client
}

func NewFacebookProvider(appID, appSecret string) *FacebookProvider {
	return &FacebookProvider{
		appID:     appID,
		appSecret: appSecret,
		client:    &http.Client{Timeout: 5 * time.Second},
	}
}

type facebookDebugTokenResponse struct {
	Data struct {
		AppID   string `json:"app_id"`
		IsValid bool   `json:"is_valid"`
		UserID  string `json:"user_id"`
		Error   *struct {
			Message string `json:"message"`
		} `json:"error"`
	} `json:"data"`
}

type facebookProfileResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Picture struct {
		Data struct {
			URL string `json:"url"`
		} `json:"data"`
	} `json:"picture"`
}

func (f *FacebookProvider) Verify(token string) (*UserInfo, error) {
	if token == "" {
		return nil, ErrInvalidToken
	}

	if err := f.verifyTokenOwnership(token); err != nil {
		return nil, err
	}

	profile, err := f.fetchProfile(token)
	if err != nil {
		return nil, err
	}

	if profile.ID == "" {
		return nil, ErrInvalidToken
	}

	return &UserInfo{
		ID:      profile.ID,
		Account: profile.Email,
		Name:    profile.Name,
		Avatar:  profile.Picture.Data.URL,
	}, nil
}

func (f *FacebookProvider) verifyTokenOwnership(token string) error {
	appAccessToken := f.appID + "|" + f.appSecret
 
	endpoint := fmt.Sprintf(
		"%s/debug_token?input_token=%s&access_token=%s",
		facebookGraphBaseURL,
		url.QueryEscape(token),
		url.QueryEscape(appAccessToken),
	)
 
	resp, err := f.client.Get(endpoint)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
 
	var result facebookDebugTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}
 
	if result.Data.Error != nil {
		return errors.New(result.Data.Error.Message)
	}
 
	if !result.Data.IsValid {
		return ErrInvalidToken
	}
 
	if result.Data.AppID != f.appID {
		return ErrInvalidToken
	}
 
	return nil
}

func (f *FacebookProvider) fetchProfile(token string) (*facebookProfileResponse, error) {
	endpoint := fmt.Sprintf(
		"%s/me?fields=id,name,email,picture&access_token=%s",
		facebookGraphBaseURL,
		url.QueryEscape(token),
	)
 
	resp, err := f.client.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
 
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("facebook graph api returned status %d", resp.StatusCode)
	}
 
	var profile facebookProfileResponse
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, err
	}
 
	return &profile, nil
}