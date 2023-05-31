package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/issy20/go-oidc-client/domain/dto"
	"github.com/issy20/go-oidc-client/handler"
	"github.com/issy20/go-oidc-client/infra/persistence"
	"github.com/issy20/go-oidc-client/middleware"
	"github.com/issy20/go-oidc-client/usecase"
	"github.com/issy20/go-oidc-client/util"
)

var oidc struct {
	clientId     string
	clientSecret string
	state        string
	nonce        string
}

const (
	response_type = "code"
	redirect_uri  = "http://localhost:8080/callback"
	grant_type    = "authorization_code"
	scope         = "clips:edit user:read:follows user:read:email openid"
)

type Token struct {
	AccessToken  string   `json:"access_token"`
	ExpiresIn    int      `json:"expires_in"`
	IDToken      string   `json:"id_token"`
	RefreshToken string   `json:"refresh_token"`
	Scope        []string `json:"scope"`
	TokenType    string   `json:"token_type"`
}

func setUp() {
	secrets, err := util.ReadJson()
	if err != nil {
		log.Fatal(err)
	}
	oidc.clientId = secrets.ClientId
	oidc.clientSecret = secrets.ClientSecret
	oidc.state = "123"
	oidc.nonce = "abc"
}

func login(w http.ResponseWriter, req *http.Request) {
	v := url.Values{}
	v.Add("response_type", response_type)
	v.Add("client_id", oidc.clientId)
	v.Add("state", oidc.state)
	v.Add("redirect_uri", redirect_uri)
	v.Add("scope", scope)
	v.Add("nonce", oidc.nonce)

	log.Printf("http redirect to : %s", fmt.Sprintf("%s?%s", dto.AuthEndpoint, v.Encode()))
	http.Redirect(w, req, fmt.Sprintf("%s?%s", dto.AuthEndpoint, v.Encode()), 302)
}

func redirect(w http.ResponseWriter, req *http.Request) {
	log.Print("redirect")
	http.Redirect(w, req, "http://localhost:3000/user", 200)
}

func tokenRequest(query url.Values, w http.ResponseWriter) (*Token, error) {
	v := url.Values{}
	v.Add("client_id", oidc.clientId)
	v.Add("client_secret", oidc.clientSecret)
	v.Add("grant_type", grant_type)
	v.Add("code", query.Get("code"))
	v.Add("redirect_uri", redirect_uri)

	req, err := http.NewRequest("POST", dto.TokenEndpoint, strings.NewReader(v.Encode()))
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	var token *Token
	err = json.Unmarshal(body, &token)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	expiresIn := time.Duration(token.ExpiresIn) * time.Second
	expiresAt := time.Now().Add(expiresIn)

	cookie := &http.Cookie{
		Name:    "session",
		Value:   string(token.AccessToken),
		Path:    "/",
		Expires: expiresAt,
	}
	http.SetCookie(w, cookie)

	log.Printf("token response :%s\n", string(body))

	return token, fmt.Errorf("%w", err)
}

func callback(w http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	token, err := tokenRequest(query, w)
	if err != nil {
		log.Println(err)
	}
	log.Print(token.IDToken)

	jwtData := decodeJWT(token.IDToken)
	err = verifyJwtSignature(jwtData, token.IDToken)
	if err != nil {
		log.Fatal("verify JWT Signature err : %w", err)
	}

	err = verifyToken(jwtData, token.AccessToken)
	if err != nil {
		log.Fatal("verifyToken is err : %w", err)
	}

	req, err = http.NewRequest("GET", dto.MyChannelEndpoint, nil)
	if err != nil {
		log.Println(err)
	}
	log.Print(jwtData.Payload.PreferredUsername)
	log.Print(jwtData.Payload.Subject)

	expiresIn := time.Duration(token.ExpiresIn) * time.Second
	expiresAt := time.Now().Add(expiresIn)

	cookie := &http.Cookie{
		Name:    "sub",
		Value:   jwtData.Payload.Subject,
		Path:    "/",
		Expires: expiresAt,
	}

	http.SetCookie(w, cookie)

	http.Redirect(w, req, "http://localhost:3000/user", 302)
}

func main() {
	setUp()

	cr := persistence.NewChannelRepository()
	cu := usecase.NewChannelUsecase(cr)
	ch := handler.NewChannelHandler(cu)

	http.HandleFunc("/login", login)
	http.HandleFunc("/callback", callback)
	http.HandleFunc("/redirect", redirect)
	http.HandleFunc("/followed-channel", middleware.AuthMiddleware(ch.GetFollowedChannel))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
