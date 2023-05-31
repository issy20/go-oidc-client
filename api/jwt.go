package main

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/issy20/go-oidc-client/domain/dto"
)

type JwtData struct {
	HeaderPayload string
	Header        map[string]interface{}
	HeaderRaw     string
	Payload       Claim
	PayloadRaw    string
	signature     []byte
}

type Claim struct {
	Audience          string `json:"aud"`
	ExpirationTime    int64  `json:"exp"`
	IssuedAt          int64  `json:"iat"`
	Issuer            string `json:"iss"`
	Subject           string `json:"sub"`
	Azp               string `json:"azp"`
	Nonce             string `json:"nonce"`
	PreferredUsername string `json:"preferred_username"`
}

type Key struct {
	Alg string `json:"alg"`
	E   string `json:"e"`
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	N   string `json:"n"`
	Use string `json:"use"`
}

type KeyResponse struct {
	Keys []Key `json:"keys"`
}

func fillB64Length(jwt string) (b64 string) {
	// JWT 文字列内の - を + に、_ を / に置換
	replace := strings.NewReplacer("-", "+", "_", "/")
	b64 = replace.Replace(jwt)
	log.Print(len(b64))

	// base64エンコードは常に文字列の長さが4の倍数である必要があるため、
	// 4の倍数でない場合、4で割った余りを取得し、余りの数だけ"="を追加する

	remainder := len(jwt) % 4
	if remainder != 0 {
		addLength := 4 - remainder
		b64 += strings.Repeat("=", addLength)
	}

	log.Print(len(b64))

	return b64
}

func decodeJWT(idToken string) JwtData {
	tmp := strings.Split(idToken, ".")
	jwtData := JwtData{
		HeaderPayload: fmt.Sprintf("%s.%s", tmp[0], tmp[1]),
		HeaderRaw:     tmp[0],
		PayloadRaw:    tmp[1],
	}

	header := fillB64Length(tmp[0])
	payLoad := fillB64Length(tmp[1])

	// Base64 エンコードされた JWTをデコード
	decHeader, err := base64.StdEncoding.DecodeString(header)
	if err != nil {
		log.Fatal("decHeader : ", err)
	}
	decPayload, err := base64.StdEncoding.DecodeString(payLoad)
	if err != nil {
		log.Fatal("decPayload :", err)
	}
	decSignature, err := base64.RawURLEncoding.DecodeString(tmp[2])
	if err != nil {
		log.Fatal("decSignature :", err)
	}

	jwtData.signature = decSignature

	//デコードされた JWT を JSON 形式から構造体に変換
	json.NewDecoder(bytes.NewReader(decHeader)).Decode(&jwtData.Header)
	json.NewDecoder(bytes.NewReader(decPayload)).Decode(&jwtData.Payload)

	log.Print(jwtData.Payload.Issuer)

	return jwtData
}

func verifyJwtSignature(jwtData JwtData, idToken string) error {
	pubkey := rsa.PublicKey{}

	var keyResponse KeyResponse
	req, err := http.NewRequest("GET", dto.KeyEndpoint, nil)
	if err != nil {
		return fmt.Errorf("http request err : %w", err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("http client err : %w", err)
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&keyResponse)
	log.Print(&keyResponse)

	for _, key := range keyResponse.Keys {
		number, _ := base64.RawURLEncoding.DecodeString(key.N)
		pubkey.N = new(big.Int).SetBytes(number)
		pubkey.E = 65537
	}

	hasher := sha256.New()
	hasher.Write([]byte(jwtData.HeaderPayload))

	err = rsa.VerifyPKCS1v15(&pubkey, crypto.SHA256, hasher.Sum(nil), jwtData.signature)
	if err != nil {
		return fmt.Errorf("verify err : %w", err)
	} else {
		log.Println("Verify success by VerifyPKCS1v15!!")
	}

	return nil
}

func verifyToken(jwtData JwtData, access_token string) error {
	if dto.Iss != jwtData.Payload.Issuer {
		return fmt.Errorf("iss not match")
	}

	if oidc.clientId != jwtData.Payload.Audience {
		return fmt.Errorf("client_id is not match")
	}

	if oidc.nonce != jwtData.Payload.Nonce {
		return fmt.Errorf("nonce is not match")
	}

	now := time.Now().Unix()
	log.Print("now:", now)
	log.Print("exp:", jwtData.Payload.ExpirationTime)
	if jwtData.Payload.ExpirationTime < now {
		return fmt.Errorf("token time limit expired")
	}

	return nil
}
