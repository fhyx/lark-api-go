package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type tokenReq struct {
	CorpID     string `json:"app_id"`
	CorpSecret string `json:"app_secret"`
}

// Token ...
type Token struct {
	AppAccessToken    string `json:"app_access_token,omitempty"`
	TenantAccessToken string `json:"tenant_access_token,omitempty"`
	ExpiresIn         int64  `json:"expires,omitempty"` // seconds // 过期时间，单位为秒（两小时失效）
	Error
}

// TokenHolder ...
type TokenHolder struct {
	currToken  *Token
	uri        string
	method     string
	apiAuths   string
	corpID     string
	corpSecret string
	expiresAt  int64
}

var (
	errEmptyAuths = errors.New("empty auth string or corpID and corpSecret")
)

func NewTokenHolder(uri string) *TokenHolder {
	return &TokenHolder{
		uri:    uri,
		method: "POST",
	}
}

func (th *TokenHolder) SetAuth(auths string) {
	th.apiAuths = auths
}

func (th *TokenHolder) SetCorp(id, secret string) {
	th.corpID = id
	th.corpSecret = secret
}

func (th *TokenHolder) Expired() bool {
	return th.expiresAt < time.Now().Unix()
}

func (th *TokenHolder) Valid() bool {
	if th.currToken == nil {
		return false
	}
	return !th.Expired()
}

func (th *TokenHolder) GetAuthToken() (token string, err error) {
	if !th.Valid() {
		logger().Debugw("token is nil or expired, refreshing it")
		th.currToken, err = th.requestToken()
		if err != nil {
			return "", err
		}
		logger().Infow("th got new OK", "token", th.currToken)
		th.expiresAt = time.Now().Unix() + th.currToken.ExpiresIn
	}
	if th.currToken.TenantAccessToken != "" {
		token = th.currToken.TenantAccessToken
	} else {
		token = th.currToken.AppAccessToken
	}

	return
}

func (th *TokenHolder) requestToken() (token *Token, err error) {
	treq := &tokenReq{
		CorpID:     th.corpID,
		CorpSecret: th.corpSecret,
	}
	logger().Debugw("request token", "corpID", th.corpID)
	body, _ := json.Marshal(treq)
	req, err := http.NewRequest("POST", th.uri, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	hc := &http.Client{Transport: tr}
	var resp []byte
	resp, err = doRequest(hc, req)

	if err != nil {
		logger().Infow("th request fail", "err", err)
		return
	}

	obj := &Token{}
	err = parseResult(resp, obj)
	if err != nil {
		logger().Infow("parseResult fail", "err", err)
		return
	}
	token = obj

	return
}
