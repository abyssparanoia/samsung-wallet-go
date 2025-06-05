# Samsung Wallet Go SDK - Event Ticket Example

This directory contains an example for Event Ticket cards using Samsung Wallet's official API specification.

## Directory Structure

```
examples/
└── event-ticket/     # Event ticket example with official Samsung Wallet API structure
```

## Environment Variables

Before running the example, set the following environment variables:

```bash
export SAMSUNG_WALLET_PARTNER_ID="your_partner_id"
export SAMSUNG_WALLET_PARTNER_PRIVATE_KEY="-----BEGIN PRIVATE KEY-----..."
export SAMSUNG_WALLET_SAMSUNG_PUBLIC_KEY="-----BEGIN CERTIFICATE-----..."
export SAMSUNG_WALLET_CERTIFICATE_ID="YMtt"

# Event Ticket Card ID (obtained from Samsung Partners Portal)
export SAMSUNG_WALLET_EVENT_TICKET_CARD_ID="your_event_ticket_card_id"
```

## Running the Example

```bash
cd examples/event-ticket
go run main.go
```

## Official Samsung Wallet API Features

The Event Ticket example demonstrates:

- **Official API compliance**: Follows Samsung Wallet's official Event Ticket specification
- **Required fields**: title, mainImg, logoImage, providerName (per Samsung specification)
- **Optional fields**: Comprehensive support for seat info, dates, styling, barcodes, etc.
- **Multi-language support**: Korean localization example included
- **JSON debugging**: BuildAsJSON() method for inspecting the final structure

## Event Ticket Subtypes

Samsung Wallet supports the following Event Ticket subtypes with type-safe constants:

- `wallet.TicketSubTypePerformances` - Concerts, theater, shows
- `wallet.TicketSubTypeSports` - Sports events, games
- `wallet.TicketSubTypeMovies` - Movie tickets
- `wallet.TicketSubTypeEntrances` - General entrance tickets (default)
- `wallet.TicketSubTypeOthers` - Other event types

You can also use strings with `SetSubTypeString()` for backward compatibility.

### Usage Examples

```go
// Type-safe approach (recommended)
eventTicket := client.NewEventTicket("TICKET001", "Concert").
    SetSubType(wallet.TicketSubTypePerformances)

// Sports event
sportsTicket := client.NewEventTicket("SPORTS001", "Basketball Game").
    SetSubType(wallet.TicketSubTypeSports)

// Movie ticket
movieTicket := client.NewEventTicket("MOVIE001", "Avengers").
    SetSubType(wallet.TicketSubTypeMovies)

// String approach (backward compatibility)
otherTicket := client.NewEventTicket("OTHER001", "Exhibition").
    SetSubTypeString("others")
```

## Key Features Demonstrated

### Required Fields

- Title and subtitle
- Main image and logo images (light/dark mode)
- Provider name

### Event Information

- Event ID, grouping ID, order ID
- Issue date, start date, end date
- Seat class, entrance, seat number
- Reservation number and user info

### Visual Styling

- Background color, font color, blink color
- Custom images and branding

### Barcode Support

- QR codes with proper Samsung formatting
- Barcode serial types and presentation formats

### Location Data

- GPS coordinates and addresses
- Multiple location support

### Multi-language Support

- Primary language settings
- Localization for Korean (example included)

## Example Output

The example will:

1. Create an Event Ticket using Samsung's official API structure
2. Set required fields (title, images, provider)
3. Configure optional fields (seat info, dates, styling)
4. Add Korean localization
5. Generate Samsung-compliant CDATA with JWE encryption
6. Display the final "Add to Wallet" link
7. Show detailed card information and JSON structure

## Security Notes

- **Samsung Public Key**: Required for JWE encryption (Samsung specification compliance)
- **Partner Private Key**: Required for JWS signing
- **Short-lived tokens**: CDATA expires in 30 seconds per Samsung security requirements
- **Environment variables**: Never commit credentials to source control

## Samsung Wallet API Compliance

This implementation follows Samsung Wallet's official Event Ticket API specification:

- Uses the official `WalletCard` structure
- Includes proper `TicketAttributes` with Samsung-defined field names
- Supports all Samsung-specified optional fields
- Generates Samsung-compliant CDATA tokens
