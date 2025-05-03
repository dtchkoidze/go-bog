package gobog

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
	"time"
)

const (
	authURL           = "https://oauth2.bog.ge/auth/realms/bog/protocol/openid-connect/token"
	reqPaymentURL     = "https://api.bog.ge/payments/v1/ecommerce/orders"
	reqPaymentInfoURL = "https://api.bog.ge/payments/v1/receipt/"
)

// Authenticate and cache the token if missing or expired
func (b *PaymentClient) auth() error {
	if b.token == nil || time.Now().After(b.tokenExpiry) {
		token, expiry, err := b.authenticate()
		if err != nil {
			return err
		}
		b.token = &token
		b.tokenExpiry = time.Now().Add(time.Duration(expiry) * time.Second)
	}
	return nil
}

// authenticate performs OAuth2 client credentials flow
func (b *PaymentClient) authenticate() (string, int, error) {
	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	body := strings.NewReader(form.Encode())

	req, err := http.NewRequest("POST", authURL, body)
	if err != nil {
		return "", 0, err
	}

	auth := base64.StdEncoding.EncodeToString([]byte(b.ClientID + ":" + b.ClientSecret))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}

	if resp.StatusCode >= 400 {
		return "", 0, fmt.Errorf("auth failed: %s", resBody)
	}

	var authResponse struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}

	err = json.Unmarshal(resBody, &authResponse)
	if err != nil {
		return "", 0, err
	}

	return authResponse.AccessToken, authResponse.ExpiresIn, nil
}

// RequestPayment sends a new payment request
func (b *PaymentClient) RequestPayment(acceptLang string, payload PaymentPayload) (*PaymentResponse, error) {
	if err := b.auth(); err != nil {
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

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("payment request failed: %s", resBody)
	}

	var paymentResp PaymentResponse
	if err := json.Unmarshal(resBody, &paymentResp); err != nil {
		return nil, err
	}

	return &paymentResp, nil
}

// RequestPaymentInfo retrieves detailed payment information
func (b *PaymentClient) RequestPaymentInfo(paymentID string) (*PaymentInfo, error) {
	if err := b.auth(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", reqPaymentInfoURL+paymentID, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+*b.token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("failed to fetch payment info: %s", resBody)
	}

	var info PaymentInfo
	err = json.Unmarshal(resBody, &info)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

// LogPayload prints the payment request for debugging
func (b *PaymentClient) LogPayload(payload *PaymentPayload) {
	log.Printf("Sending payment request. Payload: %+v\n\n", payload)
}

// VerifyCallbackSignature ensures the incoming callback was signed properly
func VerifyCallbackSignature(rawBody []byte, signatureHeader string, publicKey *rsa.PublicKey) error {
	decodedSig, err := base64.StdEncoding.DecodeString(signatureHeader)
	if err != nil {
		return fmt.Errorf("invalid signature encoding: %w", err)
	}
	hash := sha256.Sum256(rawBody)
	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash[:], decodedSig)
}
