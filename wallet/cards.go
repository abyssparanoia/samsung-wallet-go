package wallet

import (
	"encoding/json"
	"fmt"
	"time"
)

// Official Samsung Wallet Event Ticket Builder

// EventTicketBuilder creates an event ticket according to Samsung Wallet API specifications
type EventTicketBuilder struct {
	walletCard WalletCard
	attributes TicketAttributes
}

// NewEventTicket creates a new event ticket builder using official Samsung Wallet structure
func NewEventTicket(refID, title string) *EventTicketBuilder {
	now := time.Now().UnixMilli()

	return &EventTicketBuilder{
		walletCard: WalletCard{
			Card: WalletCardBody{
				Type:    "ticket",
				SubType: string(TicketSubTypeEntrances), // Default to entrances, can be changed
				Data: []WalletCardData{
					{
						RefID:      refID,
						CreatedAt:  now,
						UpdatedAt:  now,
						Language:   "en",
						Attributes: make(WalletCardAttributes),
					},
				},
			},
		},
		attributes: TicketAttributes{
			Title: title,
		},
	}
}

// SetSubType sets the ticket subtype using defined constants
func (b *EventTicketBuilder) SetSubType(subType TicketSubType) *EventTicketBuilder {
	b.walletCard.Card.SubType = string(subType)
	return b
}

// SetSubTypeString sets the ticket subtype using string (for backward compatibility)
func (b *EventTicketBuilder) SetSubTypeString(subType string) *EventTicketBuilder {
	b.walletCard.Card.SubType = subType
	return b
}

// SetLanguage sets the primary language for the ticket
func (b *EventTicketBuilder) SetLanguage(language string) *EventTicketBuilder {
	if len(b.walletCard.Card.Data) > 0 {
		b.walletCard.Card.Data[0].Language = language
	}
	return b
}

// Required fields setters

// SetMainImage sets the main ticket image URL (required, max 512 kB)
func (b *EventTicketBuilder) SetMainImage(imageURL string) *EventTicketBuilder {
	b.attributes.MainImg = imageURL
	return b
}

// SetLogoImage sets the logo image URL (required, max 256 kB)
func (b *EventTicketBuilder) SetLogoImage(logoURL string) *EventTicketBuilder {
	b.attributes.LogoImage = logoURL
	return b
}

// SetLogoImages sets both dark and light mode logo images
func (b *EventTicketBuilder) SetLogoImages(lightURL, darkURL string) *EventTicketBuilder {
	b.attributes.LogoImage = lightURL
	b.attributes.LogoImageLightURL = lightURL
	b.attributes.LogoImageDarkURL = darkURL
	return b
}

// SetProviderName sets the ticket provider name (required, max 32 chars)
func (b *EventTicketBuilder) SetProviderName(providerName string) *EventTicketBuilder {
	b.attributes.ProviderName = providerName
	return b
}

// Optional fields setters

// SetSubtitle sets the subtitle field (max 32 chars)
func (b *EventTicketBuilder) SetSubtitle(subtitle string) *EventTicketBuilder {
	b.attributes.Subtitle1 = subtitle
	return b
}

// SetEventInfo sets event-related information
func (b *EventTicketBuilder) SetEventInfo(eventID, groupingID, orderID string) *EventTicketBuilder {
	if eventID != "" {
		b.attributes.EventID = eventID
	}
	if groupingID != "" {
		b.attributes.GroupingID = groupingID
	}
	if orderID != "" {
		b.attributes.OrderID = orderID
	}
	return b
}

// SetSeatInfo sets seat-related information
func (b *EventTicketBuilder) SetSeatInfo(seatClass, entrance, seatNumber string) *EventTicketBuilder {
	if seatClass != "" {
		b.attributes.SeatClass = seatClass
	}
	if entrance != "" {
		b.attributes.Entrance = entrance
	}
	if seatNumber != "" {
		b.attributes.SeatNumber = seatNumber
	}
	return b
}

// SetTicketInfo sets general ticket information
func (b *EventTicketBuilder) SetTicketInfo(reservationNumber, user, certification, grade string) *EventTicketBuilder {
	if reservationNumber != "" {
		b.attributes.ReservationNumber = reservationNumber
	}
	if user != "" {
		b.attributes.User = user
	}
	if certification != "" {
		b.attributes.Certification = certification
	}
	if grade != "" {
		b.attributes.Grade = grade
	}
	return b
}

// SetDates sets issue, start, and end dates
func (b *EventTicketBuilder) SetDates(issueDate, startDate, endDate *time.Time) *EventTicketBuilder {
	if issueDate != nil {
		b.attributes.IssueDate = issueDate.UnixMilli()
	}
	if startDate != nil {
		b.attributes.StartDate = startDate.UnixMilli()
	}
	if endDate != nil {
		b.attributes.EndDate = endDate.UnixMilli()
	}
	return b
}

// SetHolderInfo sets holder information including photo
func (b *EventTicketBuilder) SetHolderInfo(holderName, photoBase64, photoFormat string) *EventTicketBuilder {
	if holderName != "" {
		b.attributes.HolderName = holderName
	}
	if photoBase64 != "" {
		b.attributes.IDPhotoData = photoBase64
	}
	if photoFormat != "" {
		b.attributes.IDPhotoFormat = photoFormat
	}
	return b
}

// SetStyling sets visual styling options
func (b *EventTicketBuilder) SetStyling(bgColor, fontColor, blinkColor string) *EventTicketBuilder {
	if bgColor != "" {
		b.attributes.BGColor = bgColor
	}
	if fontColor != "" {
		b.attributes.FontColor = fontColor
	}
	if blinkColor != "" {
		b.attributes.BlinkColor = blinkColor
	}
	return b
}

// SetBarcode sets barcode information
func (b *EventTicketBuilder) SetBarcode(value, serialType, ptFormat, ptSubFormat string) *EventTicketBuilder {
	if value != "" {
		b.attributes.BarcodeValue = value
	}
	if serialType != "" {
		b.attributes.BarcodeSerialType = serialType
	}
	if ptFormat != "" {
		b.attributes.BarcodePTFormat = ptFormat
	}
	if ptSubFormat != "" {
		b.attributes.BarcodePTSubFormat = ptSubFormat
	}
	return b
}

// SetQRCode is a convenience method to set QR code barcode
func (b *EventTicketBuilder) SetQRCode(value string) *EventTicketBuilder {
	return b.SetBarcode(value, "QRCODE", "QRCODESERIAL", "QR_CODE")
}

// SetPersonInfo sets person information as JSON string
func (b *EventTicketBuilder) SetPersonInfo(personJSON string) *EventTicketBuilder {
	b.attributes.Person1 = personJSON
	return b
}

// SetPersonInfoFromStruct sets person information from a struct
func (b *EventTicketBuilder) SetPersonInfoFromStruct(persons []PersonInfo) *EventTicketBuilder {
	if len(persons) > 0 {
		personData := struct {
			Person []PersonInfo `json:"person"`
		}{
			Person: persons,
		}

		if jsonData, err := json.Marshal(personData); err == nil {
			b.attributes.Person1 = string(jsonData)
		}
	}
	return b
}

// PersonInfo represents person information for tickets
type PersonInfo struct {
	Category string `json:"category"`
	Count    int    `json:"count"`
}

// SetLocations sets location information as JSON string
func (b *EventTicketBuilder) SetLocations(locationsJSON string) *EventTicketBuilder {
	b.attributes.Locations = locationsJSON
	return b
}

// SetLocationsFromStruct sets location information from structs
func (b *EventTicketBuilder) SetLocationsFromStruct(locations []TicketLocation) *EventTicketBuilder {
	if jsonData, err := json.Marshal(locations); err == nil {
		b.attributes.Locations = string(jsonData)
	}
	return b
}

// TicketLocation represents location information for tickets
type TicketLocation struct {
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
	Address string  `json:"address"`
	Name    string  `json:"name"`
}

// SetNoticeDescription sets notice description (supports HTML, max 1024 chars)
func (b *EventTicketBuilder) SetNoticeDescription(noticeHTML string) *EventTicketBuilder {
	b.attributes.NoticeDesc = noticeHTML
	return b
}

// SetGroupInfo sets group information fields
func (b *EventTicketBuilder) SetGroupInfo(groupInfo1, groupInfo2, groupInfo3 string) *EventTicketBuilder {
	if groupInfo1 != "" {
		b.attributes.GroupInfo1 = groupInfo1
	}
	if groupInfo2 != "" {
		b.attributes.GroupInfo2 = groupInfo2
	}
	if groupInfo3 != "" {
		b.attributes.GroupInfo3 = groupInfo3
	}
	return b
}

// SetCustomerServiceInfo sets customer service information as JSON string
func (b *EventTicketBuilder) SetCustomerServiceInfo(csInfoJSON string) *EventTicketBuilder {
	b.attributes.CSInfo = csInfoJSON
	return b
}

// SetCustomerServiceInfoFromStruct sets customer service information from struct
func (b *EventTicketBuilder) SetCustomerServiceInfoFromStruct(csInfo CustomerServiceInfo) *EventTicketBuilder {
	if jsonData, err := json.Marshal(csInfo); err == nil {
		b.attributes.CSInfo = string(jsonData)
	}
	return b
}

// CustomerServiceInfo represents customer service contact information
type CustomerServiceInfo struct {
	Call      string `json:"call,omitempty"`
	Email     string `json:"email,omitempty"`
	Website   string `json:"website,omitempty"`
	Instagram string `json:"instagram,omitempty"`
	YouTube   string `json:"youtube,omitempty"`
	Facebook  string `json:"facebook,omitempty"`
}

// SetAppLink sets app link information
func (b *EventTicketBuilder) SetAppLink(appLinkName, appLinkLogo, appLinkData string) *EventTicketBuilder {
	if appLinkName != "" {
		b.attributes.AppLinkName = appLinkName
	}
	if appLinkLogo != "" {
		b.attributes.AppLinkLogo = appLinkLogo
	}
	if appLinkData != "" {
		b.attributes.AppLinkData = appLinkData
	}
	return b
}

// SetClassification sets ticket classification (ONETIME, REGULAR, ANNUAL)
func (b *EventTicketBuilder) SetClassification(classification string) *EventTicketBuilder {
	b.attributes.Classification = classification
	return b
}

// AddLocalization adds localized attributes for multi-language support
func (b *EventTicketBuilder) AddLocalization(language string, localizedAttrs map[string]interface{}) *EventTicketBuilder {
	if len(b.walletCard.Card.Data) > 0 {
		localization := WalletCardLocalization{
			Language:   language,
			Attributes: localizedAttrs,
		}
		b.walletCard.Card.Data[0].Localization = append(b.walletCard.Card.Data[0].Localization, localization)
	}
	return b
}

// Build returns the final WalletCard structure
func (b *EventTicketBuilder) Build() WalletCard {
	// Convert TicketAttributes to map for attributes
	if len(b.walletCard.Card.Data) > 0 {
		attributesMap := make(WalletCardAttributes)

		// Convert struct to map using JSON marshaling/unmarshaling
		if jsonData, err := json.Marshal(b.attributes); err == nil {
			var tempMap map[string]interface{}
			if err := json.Unmarshal(jsonData, &tempMap); err == nil {
				// Only include non-zero values
				for key, value := range tempMap {
					if value != nil && value != "" && value != int64(0) {
						attributesMap[key] = value
					}
				}
			}
		}

		b.walletCard.Card.Data[0].Attributes = attributesMap
	}

	return b.walletCard
}

// BuildAsJSON returns the wallet card as JSON string
func (b *EventTicketBuilder) BuildAsJSON() (string, error) {
	walletCard := b.Build()
	jsonData, err := json.MarshalIndent(walletCard, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal wallet card to JSON: %v", err)
	}
	return string(jsonData), nil
}
