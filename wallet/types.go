package wallet

import "time"

// CardType represents the type of wallet card
type CardType string

const (
	CardTypeEventTicket CardType = "event_ticket"
)

// TicketSubType represents the subtype of event tickets according to Samsung Wallet API
type TicketSubType string

const (
	TicketSubTypePerformances TicketSubType = "performances" // Concerts, theater, shows
	TicketSubTypeSports       TicketSubType = "sports"       // Sports events, games
	TicketSubTypeMovies       TicketSubType = "movies"       // Movie tickets
	TicketSubTypeEntrances    TicketSubType = "entrances"    // General entrance tickets
	TicketSubTypeOthers       TicketSubType = "others"       // Other event types
)

// CardState represents the state of a wallet card
type CardState string

const (
	CardStateAdded    CardState = "ADDED"
	CardStateDeleted  CardState = "DELETED"
	CardStateCanceled CardState = "CANCELED"
)

// Config holds the configuration for Samsung Wallet client
type Config struct {
	PartnerID         string `json:"partner_id"`          // Partner ID from Samsung Partners Portal
	PartnerPrivateKey string `json:"partner_private_key"` // Partner's RSA private key for JWT signing
	SamsungPublicKey  string `json:"samsung_public_key"`  // Samsung's public key for JWE encryption
	CertificateID     string `json:"certificate_id"`      // 4 digit alphanumeric from Partners Portal
	BaseURL           string `json:"base_url,omitempty"`  // Optional: Custom base URL
}

// Samsung Wallet Official API Card Structures

// WalletCard represents the official Samsung Wallet card structure
type WalletCard struct {
	Card WalletCardBody `json:"card"`
}

// WalletCardBody represents the main card body
type WalletCardBody struct {
	Type    string           `json:"type"`
	SubType string           `json:"subType"`
	Data    []WalletCardData `json:"data"`
}

// WalletCardData represents individual card data
type WalletCardData struct {
	RefID        string                   `json:"refId"`
	CreatedAt    int64                    `json:"createdAt"`
	UpdatedAt    int64                    `json:"updatedAt"`
	Language     string                   `json:"language"`
	Attributes   WalletCardAttributes     `json:"attributes"`
	Localization []WalletCardLocalization `json:"localization,omitempty"`
}

// WalletCardLocalization represents card localization data
type WalletCardLocalization struct {
	Language   string                   `json:"language"`
	Attributes WalletCardLocalizedAttrs `json:"attributes"`
}

// WalletCardLocalizedAttrs represents localized attributes
type WalletCardLocalizedAttrs map[string]interface{}

// Ticket-specific structures

// TicketAttributes represents Event Ticket attributes according to Samsung Wallet API
type TicketAttributes struct {
	// Required fields
	Title        string `json:"title"`        // Main title (max 32 chars)
	MainImg      string `json:"mainImg"`      // URL for main ticket image (max 512 kB)
	LogoImage    string `json:"logoImage"`    // Logo image URL (max 256 kB)
	ProviderName string `json:"providerName"` // Ticket provider name (max 32 chars)

	// Logo images for dark/light mode
	LogoImageDarkURL  string `json:"logoImage.darkUrl,omitempty"`  // Logo image URL in dark mode
	LogoImageLightURL string `json:"logoImage.lightUrl,omitempty"` // Logo image URL in light mode

	// Optional common fields
	Subtitle1         string `json:"subtitle1,omitempty"`         // Auxiliary field (max 32 chars)
	Category          string `json:"category,omitempty"`          // Ticket category (max 16 chars) - deprecated
	EventID           string `json:"eventId,omitempty"`           // Event identifier (max 32 chars)
	GroupingID        string `json:"groupingId,omitempty"`        // Grouping identifier (max 32 chars)
	OrderID           string `json:"orderId,omitempty"`           // Order identifier (max 32 chars)
	WideImage         string `json:"wideImage,omitempty"`         // Wide horizontal image URL (max 256 kB)
	ProviderViewLink  string `json:"providerViewLink,omitempty"`  // Link to additional provider info (max 512 chars)
	Classification    string `json:"classification,omitempty"`    // ONETIME, REGULAR, or ANNUAL (default: ONETIME)
	HolderName        string `json:"holderName,omitempty"`        // Name of card holder (max 64 chars)
	IDPhotoData       string `json:"idPhoto.data,omitempty"`      // Holder's photo Base64 (max 20k)
	IDPhotoFormat     string `json:"idPhoto.format,omitempty"`    // Image format (jpeg, png)
	IDPhotoStatus     string `json:"idPhoto.status,omitempty"`    // Status (UNCHANGED)
	Grade             string `json:"grade,omitempty"`             // Ticket grade (max 32 chars)
	SeatClass         string `json:"seatClass,omitempty"`         // Seat class (max 32 chars)
	Entrance          string `json:"entrance,omitempty"`          // Entrance gate (max 64 chars)
	SeatNumber        string `json:"seatNumber,omitempty"`        // Seat location (max 256 chars)
	SeatLayoutImage   string `json:"seatLayoutImage,omitempty"`   // Seat layout image URL (max 256 kB)
	IssueDate         int64  `json:"issueDate,omitempty"`         // Issue date (epoch timestamp)
	ReservationNumber string `json:"reservationNumber,omitempty"` // Reservation number (max 32 chars)
	User              string `json:"user,omitempty"`              // User name (max 32 chars)
	Certification     string `json:"certification,omitempty"`     // Certification (max 32 chars)
	StartDate         int64  `json:"startDate,omitempty"`         // Event start date (epoch timestamp)
	EndDate           int64  `json:"endDate,omitempty"`           // Event end date (epoch timestamp)
	Person1           string `json:"person1,omitempty"`           // Person info JSON string (max 512 chars)
	Locations         string `json:"locations,omitempty"`         // Locations JSON string (max 512 chars)
	NoticeDesc        string `json:"noticeDesc,omitempty"`        // Notice description (max 1024 chars)
	GroupInfo1        string `json:"groupInfo1,omitempty"`        // Group info 1 (max 32 chars)
	GroupInfo2        string `json:"groupInfo2,omitempty"`        // Group info 2 (max 32 chars)
	GroupInfo3        string `json:"groupInfo3,omitempty"`        // Group info 3 (max 32 chars)
	CSInfo            string `json:"csInfo,omitempty"`            // Customer service info JSON string (max 512 chars)
	AppLinkName       string `json:"appLinkName,omitempty"`       // App link name (max 32 chars)
	AppLinkLogo       string `json:"appLinkLogo,omitempty"`       // App link logo URL (max 256 kB)
	AppLinkData       string `json:"appLinkData,omitempty"`       // App link data (max 512 chars)

	// Styling fields
	BGColor    string `json:"bgColor,omitempty"`    // Background color
	FontColor  string `json:"fontColor,omitempty"`  // Font color (light/dark or hex)
	BlinkColor string `json:"blinkColor,omitempty"` // Blink color

	// Barcode fields
	BarcodeValue          string `json:"barcode.value,omitempty"`                // Barcode value (max 4096 chars)
	BarcodeSerialType     string `json:"barcode.serialType,omitempty"`           // Serial type (QRCODE, BARCODE, etc.)
	BarcodePTFormat       string `json:"barcode.ptFormat,omitempty"`             // Presentation format
	BarcodePTSubFormat    string `json:"barcode.ptSubFormat,omitempty"`          // Presentation sub-format
	BarcodeErrorCorrLevel string `json:"barcode.errorCorrectionLevel,omitempty"` // Error correction level (L/M/Q/H)
	BarcodeInterval       string `json:"barcode.interval,omitempty"`             // Update interval

	// Provisioning fields
	ProvisionData     string `json:"provision.data,omitempty"`     // Provisioning data (max 512 chars)
	ProvisionInterval string `json:"provision.interval,omitempty"` // Provisioning interval

	// Related coupon fields (i: 1~3)
	RelCoupon1Title                string `json:"relCoupon1.title,omitempty"`                // Related coupon 1 title
	RelCoupon1Subtitle             string `json:"relCoupon1.subtitle,omitempty"`             // Related coupon 1 subtitle
	RelCoupon1ProviderName         string `json:"relCoupon1.providerName,omitempty"`         // Related coupon 1 provider name
	RelCoupon1ImageFileSrc         string `json:"relCoupon1.imageFileSrc,omitempty"`         // Related coupon 1 image URL
	RelCoupon1NoticeDescription    string `json:"relCoupon1.noticeDescription,omitempty"`    // Related coupon 1 notice
	RelCoupon1NotificationTime     int64  `json:"relCoupon1.notificationTime,omitempty"`     // Related coupon 1 notification time
	RelCoupon1Value                string `json:"relCoupon1.value,omitempty"`                // Related coupon 1 value
	RelCoupon1SerialType           string `json:"relCoupon1.serialType,omitempty"`           // Related coupon 1 serial type
	RelCoupon1PTFormat             string `json:"relCoupon1.ptFormat,omitempty"`             // Related coupon 1 PT format
	RelCoupon1PTSubFormat          string `json:"relCoupon1.ptSubFormat,omitempty"`          // Related coupon 1 PT sub-format
	RelCoupon1ErrorCorrectionLevel string `json:"relCoupon1.errorCorrectionLevel,omitempty"` // Related coupon 1 error correction
	RelCoupon2Title                string `json:"relCoupon2.title,omitempty"`                // Related coupon 2 title
	RelCoupon2Subtitle             string `json:"relCoupon2.subtitle,omitempty"`             // Related coupon 2 subtitle
	RelCoupon2ProviderName         string `json:"relCoupon2.providerName,omitempty"`         // Related coupon 2 provider name
	RelCoupon2ImageFileSrc         string `json:"relCoupon2.imageFileSrc,omitempty"`         // Related coupon 2 image URL
	RelCoupon2NoticeDescription    string `json:"relCoupon2.noticeDescription,omitempty"`    // Related coupon 2 notice
	RelCoupon2NotificationTime     int64  `json:"relCoupon2.notificationTime,omitempty"`     // Related coupon 2 notification time
	RelCoupon2Value                string `json:"relCoupon2.value,omitempty"`                // Related coupon 2 value
	RelCoupon2SerialType           string `json:"relCoupon2.serialType,omitempty"`           // Related coupon 2 serial type
	RelCoupon2PTFormat             string `json:"relCoupon2.ptFormat,omitempty"`             // Related coupon 2 PT format
	RelCoupon2PTSubFormat          string `json:"relCoupon2.ptSubFormat,omitempty"`          // Related coupon 2 PT sub-format
	RelCoupon2ErrorCorrectionLevel string `json:"relCoupon2.errorCorrectionLevel,omitempty"` // Related coupon 2 error correction
	RelCoupon3Title                string `json:"relCoupon3.title,omitempty"`                // Related coupon 3 title
	RelCoupon3Subtitle             string `json:"relCoupon3.subtitle,omitempty"`             // Related coupon 3 subtitle
	RelCoupon3ProviderName         string `json:"relCoupon3.providerName,omitempty"`         // Related coupon 3 provider name
	RelCoupon3ImageFileSrc         string `json:"relCoupon3.imageFileSrc,omitempty"`         // Related coupon 3 image URL
	RelCoupon3NoticeDescription    string `json:"relCoupon3.noticeDescription,omitempty"`    // Related coupon 3 notice
	RelCoupon3NotificationTime     int64  `json:"relCoupon3.notificationTime,omitempty"`     // Related coupon 3 notification time
	RelCoupon3Value                string `json:"relCoupon3.value,omitempty"`                // Related coupon 3 value
	RelCoupon3SerialType           string `json:"relCoupon3.serialType,omitempty"`           // Related coupon 3 serial type
	RelCoupon3PTFormat             string `json:"relCoupon3.ptFormat,omitempty"`             // Related coupon 3 PT format
	RelCoupon3PTSubFormat          string `json:"relCoupon3.ptSubFormat,omitempty"`          // Related coupon 3 PT sub-format
	RelCoupon3ErrorCorrectionLevel string `json:"relCoupon3.errorCorrectionLevel,omitempty"` // Related coupon 3 error correction
}

// WalletCardAttributes represents generic attributes that can be used for any card type
type WalletCardAttributes map[string]interface{}

// Legacy types for backward compatibility

// CardData represents the basic card data structure
type CardData struct {
	PartnerID       string            `json:"partner_id"` // Changed from service_id to match Samsung naming
	CardType        CardType          `json:"card_type"`
	CardID          string            `json:"card_id"`
	Name            string            `json:"name"`
	Description     string            `json:"description,omitempty"`
	Logo            string            `json:"logo,omitempty"`
	Images          []string          `json:"images,omitempty"`
	BackgroundColor string            `json:"background_color,omitempty"`
	TextColor       string            `json:"text_color,omitempty"`
	Fields          []CardField       `json:"fields,omitempty"`
	Barcode         *Barcode          `json:"barcode,omitempty"`
	Locations       []Location        `json:"locations,omitempty"`
	ExpiryDate      *time.Time        `json:"expiry_date,omitempty"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	Metadata        map[string]string `json:"metadata,omitempty"`
}

// CardField represents a field in the card
type CardField struct {
	Key         string `json:"key"`
	Label       string `json:"label"`
	Value       string `json:"value"`
	Type        string `json:"type,omitempty"`
	Format      string `json:"format,omitempty"`
	Placeholder string `json:"placeholder,omitempty"`
}

// Barcode represents barcode information
type Barcode struct {
	Type    string `json:"type"` // QR, PDF417, Code128, etc.
	Value   string `json:"value"`
	Format  string `json:"format,omitempty"`
	AltText string `json:"alt_text,omitempty"`
}

// Location represents a location for the card
type Location struct {
	Name      string  `json:"name"`
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
}

// ATWLinkRequest represents a request to create an Add to Wallet link
type ATWLinkRequest struct {
	PartnerID   string   `json:"partner_id"` // Changed from service_id
	CardData    CardData `json:"card_data"`
	LinkType    string   `json:"link_type"` // "data_transmit" or "data_fetch"
	CallbackURL string   `json:"callback_url,omitempty"`
}

// ATWLinkResponse represents the response from creating an ATW link
type ATWLinkResponse struct {
	Link   string `json:"link"`
	Token  string `json:"token,omitempty"`
	Status string `json:"status"`
}

// CardStateCallback represents the callback data from Samsung Wallet
type CardStateCallback struct {
	PartnerID   string    `json:"partner_id"` // Changed from service_id
	CardID      string    `json:"card_id"`
	Event       CardState `json:"event"`
	CountryCode string    `json:"country_code"`
	Timestamp   time.Time `json:"timestamp"`
	UserID      string    `json:"user_id,omitempty"`
	DeviceID    string    `json:"device_id,omitempty"`
}

// UpdateCardRequest represents a request to update a card
type UpdateCardRequest struct {
	PartnerID string   `json:"partner_id"` // Changed from service_id
	CardID    string   `json:"card_id"`
	CardData  CardData `json:"card_data"`
}

// CancelCardRequest represents a request to cancel cards
type CancelCardRequest struct {
	PartnerID string `json:"partner_id"` // Changed from service_id
	EventID   string `json:"event_id"`
	Reason    string `json:"reason,omitempty"`
}

// APIError represents an error response from the Samsung Wallet API
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e *APIError) Error() string {
	return e.Message
}
