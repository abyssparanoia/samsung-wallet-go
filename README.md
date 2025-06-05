# Samsung Wallet Go SDK

Samsung Wallet "Add to Samsung Wallet" (ATW) service integration SDK for Go.

## Features

- **CDATA Generation**: Compliant with Samsung's official specification (JWE encryption + JWS signing)
- **30-Second Token Expiry**: Security-compliant token expiration
- **Samsung-Specific JWT Headers**: Proper `cty`, `partnerId`, `ver`, `certificateId`, and `utc` headers
- Add to Samsung Wallet link generation
- Server API client
- Callback handling
- Card state management
- Support for various card types (Boarding Pass, Event Ticket, Coupon, Gift Card, Loyalty, etc.)

## Installation

```bash
go get github.com/abyssparanoia/samsung-wallet-go
```

## Quick Start

```go
package main

import (
    "github.com/abyssparanoia/samsung-wallet-go/wallet"
)

func main() {
    client, err := wallet.NewClient(&wallet.Config{
        PartnerID:         "your-partner-id",         // Partner ID from Samsung Partners Portal
        PartnerPrivateKey: "your-partner-private-key", // Your RSA private key for JWT signing
        SamsungPublicKey:  "your-samsung-public-key",  // Samsung's public key/certificate for JWE encryption
        CertificateID:     "your-cert-id",            // 4 digit alphanumeric from Partners Portal
        Environment:       wallet.EnvironmentSandbox, // or wallet.EnvironmentProduction
    })
    if err != nil {
        panic(err)
    }

    // Each card type has its own cardID obtained from Partners Portal
    cardID := "your-boarding-pass-card-id" // Specific to boarding pass card type

    // Create card data
    cardData := wallet.CardData{
        PartnerID:   "your-partner-id",  // Changed from ServiceID to match Samsung naming
        CardType:    wallet.CardTypeBoardingPass,
        CardID:      "BP001",
        Name:        "Flight SE701",
        Description: "Seoul to San Francisco",
        // ... other fields
    }

    // Create Add to Wallet link (CDATA expires in 30 seconds)
    link, err := client.CreateATWLink(cardID, cardData, "data_transmit")
    if err != nil {
        panic(err)
    }

    fmt.Println("ATW Link:", link)
}
```

## Builder Pattern (Recommended)

For easier and more intuitive card creation, use the Builder pattern:

```go
package main

import (
    "time"
    "github.com/abyssparanoia/samsung-wallet-go/wallet"
)

func main() {
    client, err := wallet.NewClient(&wallet.Config{
        PartnerID:         "your-partner-id",
        PartnerPrivateKey: "your-partner-private-key",
        SamsungPublicKey:  "your-samsung-public-key",
        CertificateID:     "your-cert-id",
    })
    if err != nil {
        panic(err)
    }

    // Create a boarding pass using the Builder pattern
    boardingPass := client.NewBoardingPass("BP001", "John Doe").
        Flight("Korean Air", "KE123", "ICN", "NRT", time.Now().Add(24*time.Hour)).
        Seat("12A", "Group 1").
        Gate("A12").
        Barcode("QRCODE", "ABC123DEF456").
        BackgroundColor("#1E3A8A").
        TextColor("#FFFFFF").
        Build()

    // Create ATW link
    cardID := "your-boarding-pass-card-id"
    link, err := client.CreateATWLink(cardID, boardingPass, "data_transmit")
    if err != nil {
        panic(err)
    }

    fmt.Println("ATW Link:", link)

    // Create an event ticket using Builder
    eventTicket := client.NewEventTicket("ET001", "BTS Concert").
        Event("Seoul Olympic Stadium", "2024-12-25", time.Now().Add(7*24*time.Hour)).
        Ticket("Section A", "Row 5", "Seat 12").
        Barcode("QRCODE", "TICKET123456").
        BackgroundColor("#8B5CF6").
        TextColor("#FFFFFF").
        AddLocation("Seoul Olympic Stadium", "Seoul, South Korea", 37.5151, 127.1240).
        Build()

    // Create coupon using Builder
    coupon := client.NewCoupon("CP001", "20% Off Everything").
        Description("Get 20% off on all items in our store").
        Discount("percentage", "20%", "$50").
        ValidPeriod(time.Now(), time.Now().Add(30*24*time.Hour)).
        Barcode("Code128", "COUPON20OFF").
        BackgroundColor("#DC2626").
        TextColor("#FFFFFF").
        Build()

    // Create gift card using Builder
    giftCard := client.NewGiftCard("GC001", "Amazon").
        Balance("100.00", "USD").
        CardNumber("1234-5678-9012-3456").
        Barcode("Code128", "GC123456789").
        BackgroundColor("#F59E0B").
        TextColor("#1F2937").
        ExpiryDate(time.Now().Add(365*24*time.Hour)).
        Build()

    // Create loyalty card using Builder
    loyaltyCard := client.NewLoyaltyCard("LC001", "Coffee Shop Rewards").
        MembershipInfo("M12345678", "Gold").
        Points("2,500", "500").
        Barcode("QRCODE", "MEMBER12345678").
        BackgroundColor("#059669").
        TextColor("#FFFFFF").
        AddField("next_reward", "Next Reward", "Free Coffee (3,000 points)").
        Build()
}
```

### Builder Pattern Benefits

- **Fluent API**: Method chaining for readable, intuitive card creation
- **Type Safety**: Card-specific methods prevent invalid field combinations
- **Reduced Boilerplate**: No need to manually create CardData structures
- **Intelligent Defaults**: Automatically sets common values like timestamps
- **Card-Specific Methods**: Each card type has specialized methods (e.g., `Flight()`, `Seat()` for boarding passes)

## Configuration

### Required Credentials from Samsung Partners Portal

1. **PartnerID**: Unique identifier for your service (called `partnerId` in Samsung documentation)
2. **PartnerPrivateKey**: Your RSA private key for JWT signing
3. **SamsungPublicKey**: Samsung's public key/certificate for JWE encryption
4. **CertificateID**: 4-digit alphanumeric certificate identifier
5. **CardIDs**: Specific identifiers for each card type you register

### Environment Variables

```bash
export SAMSUNG_WALLET_PARTNER_ID="your-partner-id"
export SAMSUNG_WALLET_PARTNER_PRIVATE_KEY="-----BEGIN RSA PRIVATE KEY-----..."
export SAMSUNG_WALLET_SAMSUNG_PUBLIC_KEY="-----BEGIN CERTIFICATE-----..."
export SAMSUNG_WALLET_CERTIFICATE_ID="YMtt"
export SAMSUNG_WALLET_BOARDING_PASS_CARD_ID="your-boarding-pass-id"
export SAMSUNG_WALLET_EVENT_TICKET_CARD_ID="your-event-ticket-id"
# ... other card type IDs
```

## CDATA Generation

This SDK generates CDATA (Card Data) tokens according to Samsung's official specification:

1. **JWE Encryption**: Card data is encrypted with Samsung's public key
2. **JWS Signing**: The encrypted payload is signed with your private key
3. **Samsung Headers**: Includes required headers (`cty: "CARD"`, `partnerId`, `ver: "3"`, `certificateId`, `utc`)
4. **30-Second Expiry**: Tokens expire in 30 seconds for security

## Development

### Available Make Commands

```bash
# Show help
make help

# Build the project
make build

# Run tests
make test

# Run tests with coverage
make test-coverage

# Format code
make fmt

# Run linter
make lint

# Run all checks
make check

# Run example
make run-example

# Setup development environment
make dev-setup

# Tidy dependencies
make mod

# Clean up
make clean
```

### Development Setup

```bash
# Setup development environment
make dev-setup

# Generate test RSA key pair
make generate-test-keys

# Watch file changes and auto-run tests (requires fswatch)
make watch
```

## Credential Acquisition Process

1. **Register Samsung Developer Account**
2. **Apply to Samsung Wallet Partners Portal**: https://partner.walletsvc.samsung.com
3. **Submit Service Information**: Undergo approval process
4. **Obtain Credentials**: After approval, receive:
   - PartnerID (called `partnerId` in Samsung documentation)
   - RSA Private Key for JWT signing
   - Samsung Public Key/Certificate for JWE encryption
   - CertificateID (4-digit alphanumeric)
   - CardIDs for each card type you register

## Important Security Notes

- **CDATA Expiry**: All CDATA tokens expire in 30 seconds for security
- **Key Management**: Store private keys securely using environment variables
- **Production Environment**: Use `wallet.EnvironmentProduction` for live services

## Documentation

- [Samsung Wallet API Documentation](https://developer.samsung.com/wallet/api_new/addto/overview.html)
- [API Guidelines](https://developer.samsung.com/wallet/api_new/addto/guidelines.html)
- [Partners Portal](https://partner.walletsvc.samsung.com)

## License

MIT
