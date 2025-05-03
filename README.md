# Bog Payment Client

This package provides a Go client for integrating with the BOG (Bank of Georgia) payment system. It supports authentication, initiating payments, retrieving payment information, and verifying callback signatures.

## Installation

Ensure you have Go installed. Then, run:

```bash
go get github.com/dtchkoidze/go-bog
```

In your code:

```go
import gobog "github.com/dtchkoidze/go-bog"
```

## Usage

### 1. Initialize Client

You need to initialize the client with your `ClientID` and `ClientSecret` provided by BOG.

```go
client := &gobog.PaymentClient{
    ClientID:     "your-client-id",
    ClientSecret: "your-client-secret",
}
```

### 2. Request Payment

The client handles authentication automatically. You only need to prepare a payload and make the request:

```go
payload := gobog.PaymentPayload{
    Amount:   1000,       // e.g., 1000 for 10.00 GEL
    Currency: "GEL",
    OrderID:  "123456789",
    Customer: "John Doe",
}

response, err := client.RequestPayment("en", payload)
if err != nil {
    log.Fatal("Payment request failed:", err)
}

fmt.Println("Payment Response:", response)
```

### 3. Request Payment Info

Retrieve information about a previously initiated payment:

```go
info, err := client.RequestPaymentInfo("payment-id")
if err != nil {
    log.Fatal("Failed to get payment info:", err)
}

fmt.Println("Payment Info:", info)
```

### 4. Verify Callback Signature

Use this to verify callback payloads using the public key provided by BOG:

```go
// Load or parse your public key (from https://api.bog.ge/docs/en/payments/standard-process/callback)
publicKey := /* *rsa.PublicKey */
signature := "base64-encoded-signature"
rawBody := []byte("callback-payload")

err := gobog.VerifyCallbackSignature(rawBody, signature, publicKey)
if err != nil {
    log.Fatal("Signature verification failed:", err)
}

fmt.Println("Signature verified successfully!")
```

## Example

```go
package main

import (
    "fmt"
    "log"

    gobog "github.com/dtchkoidze/go-bog"
)

func main() {
    client := &gobog.PaymentClient{
        ClientID:     "your-client-id",
        ClientSecret: "your-client-secret",
    }

    payload := gobog.PaymentPayload{
        Amount:   1000,
        Currency: "GEL",
        OrderID:  "123456789",
        Customer: "John Doe",
    }

    response, err := client.RequestPayment("en", payload)
    if err != nil {
        log.Fatal("Payment request failed:", err)
    }
    fmt.Println("Payment Response:", response)

    info, err := client.RequestPaymentInfo("payment-id")
    if err != nil {
        log.Fatal("Failed to get payment info:", err)
    }
    fmt.Println("Payment Info:", info)
}
```

## Exported Functions

### `auth() error`

Handles internal token authentication (called automatically).

### `RequestPayment(acceptLang string, payload PaymentPayload) (*PaymentResponse, error)`

Initiates a payment.

### `RequestPaymentInfo(paymentID string) (*PaymentInfo, error)`

Fetches payment details by ID.

### `VerifyCallbackSignature(rawBody []byte, signatureHeader string, publicKey *rsa.PublicKey) error`

Verifies callback payload integrity.

### `LogPayload(payload *PaymentPayload)`

Logs the outgoing payment payload for debugging.

## Dependencies

Built-in Go libraries used:
* `crypto/rsa`, `crypto/sha256`
* `encoding/json`, `encoding/base64`
* `net/http`, `io`, `log`

## Support
If you encounter any issues or have questions about this package:

Check the GitHub Issues page to see if your problem has been reported
Open a new issue with details about the problem, including code samples and error messages
For urgent matters, contact the maintainer directly at dtchkoiddze@gmail.com

## License

MIT License. See the LICENSE file.