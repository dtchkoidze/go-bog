# Bog Payment Client

This package provides a Go client for integrating with the BOG (Bank of Georgia) payment system. It includes functionality to authenticate, request payments, retrieve payment information, and verify callback signatures.

## Installation

To use this package, ensure you have Go installed on your machine. You can install it from the official site: [Go Installation](https://golang.org/doc/install).

To use this package, import it into your Go project:

```go
import "path/to/your/package"
```

## Usage

### 1. Authentication

Before making any payment requests, the client must authenticate with the BOG API. This is done using the `ClientID` and `ClientSecret` provided by the BOG API. RequestPayment will do this for you anyway.

```go
client := &BogPaymentClient{
    ClientID:     "your-client-id",
    ClientSecret: "your-client-secret",
}

token, err := client.authentificate()
if err != nil {
    log.Fatal("Authentication failed:", err)
}
fmt.Println("Access Token:", token)
```

### 2. Request Payment

Once authenticated, you can create a payment request using the `RequestPayment` function.

```go
payload := BogPaymentPayload{
    Amount:    1000, // Example: 1000 for 10.00 units
    Currency:  "GEL", // Example: Georgian Lari
    OrderID:   "123456789",
    Customer:  "John Doe",
}

response, err := client.RequestPayment("en", payload)
if err != nil {
    log.Fatal("Payment request failed:", err)
}

fmt.Println("Payment Response:", response)
```

### 3. Request Payment Info

After initiating a payment, you can request information about the payment using the `RequestPaymentInfo` function.

```go
paymentInfo, err := client.RequestPaymentInfo("payment-id")
if err != nil {
    log.Fatal("Failed to get payment info:", err)
}

fmt.Println("Payment Info:", paymentInfo)
```

### 4. Verify Callback Signature

When receiving a callback from the BOG API, you need to verify the integrity of the callback data by checking the signature. Use the `VerifyCallbackSignature` function to verify the signature.

```go
publicKey := // Load or parse your public key, find it on: https://api.bog.ge/docs/en/payments/standard-process/callback
signature := "base64-encoded-signature"
rawBody := []byte("callback-payload")

err := VerifyCallbackSignature(rawBody, signature, publicKey)
if err != nil {
    log.Fatal("Signature verification failed:", err)
}

fmt.Println("Signature verified successfully!")
```

### Example Code

```go

func main() {
	// Initialize client
	client := &BogPaymentClient{
		ClientID:     "your-client-id",
		ClientSecret: "your-client-secret",
	}

	// Authenticate client
	token, err := client.authentificate()
	if err != nil {
		log.Fatal("Authentication failed:", err)
	}
	fmt.Println("Access Token:", token)

	// Prepare payment payload
	payload := BogPaymentPayload{
		Amount:   1000,
		Currency: "GEL",
		OrderID:  "123456789",
		Customer: "John Doe",
	}

	// Request payment
	response, err := client.RequestPayment("en", payload)
	if err != nil {
		log.Fatal("Payment request failed:", err)
	}
	fmt.Println("Payment Response:", response)

	// Request payment information
	paymentInfo, err := client.RequestPaymentInfo("payment-id")
	if err != nil {
		log.Fatal("Failed to get payment info:", err)
	}
	fmt.Println("Payment Info:", paymentInfo)
}
```

## Functions

### `authentificate()`

Authenticates with the BOG API using the client credentials (`ClientID` and `ClientSecret`) and retrieves an access token.

### `RequestPayment(acceptLang string, payload BogPaymentPayload)`

Initiates a payment request to the BOG API. Requires an access token obtained from authentication and a payment payload.

### `RequestPaymentInfo(paymentID string)`

Retrieves information about a specific payment using the provided `paymentID`.

### `VerifyCallbackSignature(rawBody []byte, signatureHeader string, publicKey *rsa.PublicKey) error`

Verifies the signature of the callback received from the BOG API to ensure the integrity of the callback data.

### `LogPayload(payload *BogPaymentPayload)`

Logs the payment payload for debugging purposes.

## Error Handling

The functions in this package return errors in the event of failure. Ensure to handle errors when calling these functions.

## Dependencies

This package uses the following libraries:

* `crypto/rsa`
* `crypto/sha256`
* `encoding/json`
* `encoding/base64`
* `net/http`
* `log`

Ensure these libraries are available in your Go environment.

## License

This package is licensed under the MIT License. See the LICENSE file for more information.

---

For more details on how to interact with the BOG API, refer to the official documentation provided by the Bank of Georgia.