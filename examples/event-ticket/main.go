//go:debug x509negativeserial=1
package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/abyssparanoia/samsung-wallet-go/wallet"
)

func main() {
	// Configure Samsung Wallet client with credentials from Partners Portal
	client, err := wallet.NewClient(&wallet.Config{
		PartnerID:         getEnv("SAMSUNG_WALLET_PARTNER_ID", ""),          // Partner ID from Samsung Partners Portal
		PartnerPrivateKey: getEnv("SAMSUNG_WALLET_PARTNER_PRIVATE_KEY", ""), // Partner's RSA private key
		SamsungPublicKey:  getEnv("SAMSUNG_WALLET_SAMSUNG_PUBLIC_KEY", ""),  // Samsung's public key/certificate
		CertificateID:     getEnv("SAMSUNG_WALLET_CERTIFICATE_ID", ""),      // 4 digit alphanumeric from Partners Portal
	})

	if err != nil {
		log.Fatalf("Failed to create Samsung Wallet client: %v", err)
	}

	fmt.Println("=== Event Ticket Example (Official Samsung Wallet API) ===")

	// Each card type has a specific cardID obtained from Partners Portal
	cardID := getEnv("SAMSUNG_WALLET_EVENT_TICKET_CARD_ID", "event_ticket_id")

	startDate := time.Now().AddDate(0, 0, -1)
	endDate := startDate.AddDate(0, 0, 2)

	// Create event ticket using official Samsung Wallet structure
	eventTicket := client.NewEventTicket(generateRefID(), "アリーナツアー2024 in 東京ドーム").
		SetSubType(wallet.TicketSubTypePerformances). // Using defined constants for type safety
		SetProviderName("MOALA Ticket").
		SetMainImage("https://developer.samsung.com/sd2/images/media/content/v1/samsung-dev-con-img.jpg").
		SetLogoImages(
			"https://developer.samsung.com/sd2/images/media/logo/samsung-logo.png",      // light mode
			"https://developer.samsung.com/sd2/images/media/logo/samsung-logo-dark.png", // dark mode
		).
		SetSubtitle("VIP Access").
		SetEventInfo("EVENT001", "GROUP001", "ORDER001").
		SetSeatInfo("VIP Section", "Main Entrance", "A-12").
		SetTicketInfo("RES123456", "舟口翔梧", "R", "Premium").
		SetDates(
			timePtr(time.Now()), // issue date
			timePtr(startDate),  // start date
			timePtr(endDate),    // end date (start + 1 day)
		).
		SetHolderInfo("舟口翔梧", "", ""). // no photo for this example
		SetStyling("#1F2937", "#FFFFFF", "#3B82F6"). // bg, font, blink colors
		SetPersonInfoFromStruct([]wallet.PersonInfo{
			{Category: "Adult", Count: 1},
		}).
		SetLocationsFromStruct([]wallet.TicketLocation{
			{
				Lat:     35.4212,
				Lng:     139.454,
				Address: "東京都文京区後楽1丁目3番61号",
				Name:    "東京ドーム",
			},
		}).
		SetNoticeDescription("<ul><li>Please arrive 30 minutes before the event</li><li>Valid photo ID required</li></ul>").
		SetGroupInfo("Adult 1", "VIP", "Premium").
		SetCustomerServiceInfoFromStruct(wallet.CustomerServiceInfo{
			Call:    "+81-90-6198-6430",
			Email:   "shogo.funaguchi@playground.live",
			Website: "https://moala.playground.live/contact",
		}).
		SetProviderViewLink(`{"count":1,"info":[{"text":"チケット画面へ","type":"web","link":"https://ticket.moala.fun/bundles/"}]}`).
		SetAppLink("MOALA Ticket", "https://../applinklogo.png", "https://ticket.moala.fun/bundles").
		SetClassification("REGULAR")

	// Add Korean localization
	eventTicket.AddLocalization("ko", map[string]interface{}{
		"title":        "삼성 개발자 컨퍼런스 2024",
		"subtitle1":    "VIP 입장권",
		"holderName":   "홍길동",
		"providerName": "삼성전자",
		"user":         "홍길동",
		"groupInfo1":   "성인 1명",
		"groupInfo2":   "VIP",
		"groupInfo3":   "프리미엄",
	})

	// Build the wallet card
	walletCard := eventTicket.Build()

	// Create Add to Wallet link using new method
	link, err := client.CreateATWLinkFromWalletCard(cardID, walletCard, "data_transmit")
	if err != nil {
		log.Fatalf("Failed to create event ticket link: %v", err)
	}

	fmt.Printf("Event Ticket ATW Link: %s\n", link)
	fmt.Printf("Note: CDATA expires in 30 seconds for security\n")

	// Display card details
	fmt.Printf("\nCard Details:\n")
	if len(walletCard.Card.Data) > 0 {
		cardData := walletCard.Card.Data[0]
		fmt.Printf("- Card Type: %s\n", walletCard.Card.Type)
		fmt.Printf("- Sub Type: %s\n", walletCard.Card.SubType)
		fmt.Printf("- Ref ID: %s\n", cardData.RefID)
		fmt.Printf("- Language: %s\n", cardData.Language)
		fmt.Printf("- Attributes: %d fields\n", len(cardData.Attributes))

		if title, ok := cardData.Attributes["title"].(string); ok {
			fmt.Printf("- Title: %s\n", title)
		}
		if providerName, ok := cardData.Attributes["providerName"].(string); ok {
			fmt.Printf("- Provider: %s\n", providerName)
		}
		if seatNumber, ok := cardData.Attributes["seatNumber"].(string); ok {
			fmt.Printf("- Seat: %s\n", seatNumber)
		}
		if barcodeValue, ok := cardData.Attributes["barcode.value"].(string); ok {
			fmt.Printf("- Barcode: %s\n", barcodeValue)
		}

		if len(cardData.Localization) > 0 {
			fmt.Printf("- Localizations: %d languages\n", len(cardData.Localization))
		}
	}

	// Print JSON representation for debugging
	if jsonStr, err := eventTicket.BuildAsJSON(); err == nil {
		fmt.Printf("\nJSON Representation:\n%s\n", jsonStr)
	}

	fmt.Printf("\nThis event ticket is ready to be added to Samsung Wallet!\n")
	fmt.Printf("Open the link above on an Android device with Samsung Wallet installed.\n")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func timePtr(t time.Time) *time.Time {
	return &t
}

// generateRefID returns a new random reference ID each run.
func generateRefID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		// Fallback to timestamp if randomness fails
		return fmt.Sprintf("ref-%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(b)
}