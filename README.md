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

You can initialize the client with your `ClientID` and `ClientSecret` provided by BOG using the constructor function:

```go
client := gobog.NewPaymentClient("your-client-id", "your-client-secret")
```

### 2. Request Payment

The client handles authentication automatically. You can use constructor functions to prepare your payment:

```go
// Create PurchaseUnits
basket := []gobog.PaymentItem{
    *gobog.NewPaymentItem(2, decimal.NewFromFloat(10.50), 123), // 2 items at 10.50 each, product ID 123
}
purchaseUnits := gobog.PurchaseUnits{
    Currency:    "GEL",
    TotalAmount: 2100, // 21.00 GEL
    Basket:      basket,
}

// Create RedirectURLs
redirectURLs := gobog.RedirectURLs{
    Success: "https://yoursite.com/success",
    Fail:    "https://yoursite.com/fail",
}

// Create the payment payload
payload := gobog.NewPaymentPayload(
    "https://yoursite.com/callback", // callback URL
    12345,                          // local order ID
    purchaseUnits,
    redirectURLs,
)

// Make the request
response, err := client.RequestPayment("en", *payload)
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

    "github.com/shopspring/decimal"
    gobog "github.com/dtchkoidze/go-bog"
)

func main() {
    // Initialize client with constructor
    client := gobog.NewPaymentClient("your-client-id", "your-client-secret")

    // Create payment items with constructor
    item1 := gobog.NewPaymentItem(2, decimal.NewFromFloat(10.50), 123)
    item2 := gobog.NewPaymentItem(1, decimal.NewFromFloat(5.00), 456)
    
    // Create purchase units
    purchaseUnits := gobog.PurchaseUnits{
        Currency:    "GEL",
        TotalAmount: 2600, // 26.00 GEL (2 x 10.50 + 1 x 5.00)
        Basket:      []gobog.PaymentItem{*item1, *item2},
    }
    
    // Create redirect URLs
    redirectURLs := gobog.RedirectURLs{
        Success: "https://yoursite.com/success",
        Fail:    "https://yoursite.com/fail",
    }
    
    // Create payment payload with constructor
    payload := gobog.NewPaymentPayload(
        "https://yoursite.com/callback",
        12345,
        purchaseUnits,
        redirectURLs,
    )

    // Request payment
    response, err := client.RequestPayment("en", *payload)
    if err != nil {
        log.Fatal("Payment request failed:", err)
    }
    fmt.Println("Payment Response:", response)

    // Get payment info
    info, err := client.RequestPaymentInfo(response.ID)
    if err != nil {
        log.Fatal("Failed to get payment info:", err)
    }
    fmt.Println("Payment Info:", info)
}
```

## Exported Functions

### Constructor Functions

#### `NewPaymentClient(clientID, clientSecret string) *PaymentClient`

Creates a new payment client with the given credentials.

#### `NewPaymentPayload(cbURL string, loID int, pus PurchaseUnits, rURLS RedirectURLs) *PaymentPayload`

Creates a new payment payload with all required fields.

#### `NewPaymentItem(q int, up decimal.Decimal, pID int) *PaymentItem`

Creates a new payment item with quantity, unit price, and product ID.

### Client Methods

#### `auth() error`

Handles internal token authentication (called automatically).

#### `RequestPayment(acceptLang string, payload PaymentPayload) (*PaymentResponse, error)`

Initiates a payment.

#### `RequestPaymentInfo(paymentID string) (*PaymentInfo, error)`

Fetches payment details by ID.

#### `VerifyCallbackSignature(rawBody []byte, signatureHeader string, publicKey *rsa.PublicKey) error`

Verifies callback payload integrity.

#### `LogPayload(payload *PaymentPayload)`

Logs the outgoing payment payload for debugging.

## Dependencies

### External Dependencies
* `github.com/shopspring/decimal` - Precise decimal arithmetic

### Built-in Go libraries used:
* `crypto/rsa`, `crypto/sha256`
* `encoding/json`, `encoding/base64`
* `net/http`, `io`, `log`
* `time`

## Support

If you encounter any issues or have questions about this package:

1. Check the [GitHub Issues](https://github.com/dtchkoidze/go-bog/issues) page to see if your problem has been reported
2. Open a new issue with details about the problem, including code samples and error messages
3. For urgent matters, contact the maintainer directly at [dtchkoiddze@gmail.com](mailto:dtchkoiddze@gmail.com)


## License

MIT License. See the LICENSE file.

