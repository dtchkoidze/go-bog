package main

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const (
	authURL           = "https://oauth2.bog.ge/auth/realms/bog/protocol/openid-connect/token"
	reqPaymentURL     = "https://api.bog.ge/payments/v1/ecommerce/orders"
	reqPaymentInfoURL = "https://api.bog.ge/payments/v1/receipt/"
)

func (b *BogPaymentClient) authentificate() (string, error) {
	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	body := strings.NewReader(form.Encode())

	req, err := http.NewRequest("POST", authURL, body)
	if err != nil {
		return "", err
	}

	auth := base64.StdEncoding.EncodeToString([]byte(b.ClientID + ":" + b.ClientSecret))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var authResponse struct {
		AccessToken string `json:"access_token"`
	}

	err = json.Unmarshal(resBody, &authResponse)
	if err != nil {
		return "", err
	}

	return authResponse.AccessToken, nil
}

func (b *BogPaymentClient) RequestPayment(acceptLang string, payload BogPaymentPayload) (*BogPaymentResponse, error) {
	err := b.auth()
	if err != nil {
		return nil, err
	}

	b.LogPayload(&payload)

	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", reqPaymentURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept-Language", acceptLang)
	req.Header.Set("Authorization", "Bearer "+*b.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var paymentResp BogPaymentResponse
	if err := json.Unmarshal(resBody, &paymentResp); err != nil {
		return nil, err
	}

	return &paymentResp, nil
}

func (b *BogPaymentClient) RequestPaymentInfo(paymentID string) (*BogPaymentInfo, error) {

	client := &http.Client{}

	req, err := http.NewRequest("GET", reqPaymentInfoURL+paymentID, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer"+*b.token)

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var info BogPaymentInfo

	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		return nil, err
	}

	return &info, nil

}

func (b *BogPaymentClient) auth() error {
	if b.token == nil {
		token, err := b.authentificate()
		if err != nil {
			return err
		}
		b.token = &token
	}
	return nil
}

func (b *BogPaymentClient) LogPayload(payload *BogPaymentPayload) {
	log.Printf("Sending payment request. Payload is: %+v\n \n", payload)
}

func VerifyCallbackSignature(rawBody []byte, signatureHeader string, publicKey *rsa.PublicKey) error {
	decodedSig, err := base64.StdEncoding.DecodeString(signatureHeader)
	if err != nil {
		return fmt.Errorf("invalid signature encoding: %w", err)
	}
	hash := sha256.Sum256(rawBody)
	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash[:], decodedSig)
}
