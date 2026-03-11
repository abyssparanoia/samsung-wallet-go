package wallet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	// API paths
	pathUpdateCard = "/v1/wallet/card/update"
	pathCancelCard = "/v1/wallet/card/cancel"
	pathGetCard    = "/v1/wallet/card/get"

	// Samsung Wallet Server API base URL
	serverAPIBaseURL = "https://api-card.walletsvc.samsung.com"

	// Link types
	linkTypeDataTransmit = "data_transmit"
	linkTypeDataFetch    = "data_fetch"
)

// Client represents the Samsung Wallet client
type Client struct {
	config     *Config
	httpClient *http.Client
	jwtManager *JWTManager
	baseURL    string
}

// NewClient creates a new Samsung Wallet client
func NewClient(config *Config) (*Client, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	if config.PartnerID == "" {
		return nil, fmt.Errorf("partner ID is required")
	}

	if config.PartnerPrivateKey == "" {
		return nil, fmt.Errorf("partner private key is required")
	}

	if config.SamsungPublicKey == "" {
		return nil, fmt.Errorf("samsung public key is required")
	}

	if config.CertificateID == "" {
		return nil, fmt.Errorf("certificate ID is required")
	}

	// Initialize JWT manager with Samsung public key and partner private key
	jwtManager, err := NewJWTManager(
		config.PartnerPrivateKey,
		config.SamsungPublicKey,
		config.PartnerID,
		config.CertificateID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create JWT manager: %v", err)
	}

	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = "https://a.swallet.link" // Samsung Wallet ATW URL
	}

	return &Client{
		config:     config,
		jwtManager: jwtManager,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		baseURL:    baseURL,
	}, nil
}

// CreateATWLink creates an Add to Samsung Wallet link with legacy CardData
func (c *Client) CreateATWLink(cardID string, cardData CardData, linkType string, callbackURL ...string) (string, error) {
	if cardID == "" {
		return "", fmt.Errorf("card ID is required (obtain from Partners Portal when registering card type)")
	}

	if linkType == "" {
		linkType = linkTypeDataTransmit // default to data transmit
	}

	var callback string
	if len(callbackURL) > 0 {
		callback = callbackURL[0]
	}

	switch linkType {
	case linkTypeDataTransmit:
		return c.createDataTransmitLink(cardID, cardData, callback)
	case linkTypeDataFetch:
		return c.createDataFetchLink(cardID, cardData, callback)
	default:
		return "", fmt.Errorf("unsupported link type: %s", linkType)
	}
}

// CreateATWLinkFromWalletCard creates an Add to Samsung Wallet link with official WalletCard structure
func (c *Client) CreateATWLinkFromWalletCard(cardID string, walletCard WalletCard, linkType string, callbackURL ...string) (string, error) {
	if cardID == "" {
		return "", fmt.Errorf("card ID is required (obtain from Partners Portal when registering card type)")
	}

	if linkType == "" {
		linkType = linkTypeDataTransmit // default to data transmit
	}

	var callback string
	if len(callbackURL) > 0 {
		callback = callbackURL[0]
	}

	switch linkType {
	case linkTypeDataTransmit:
		return c.createDataTransmitLinkFromWalletCard(cardID, walletCard, callback)
	case linkTypeDataFetch:
		return c.createDataFetchLinkFromWalletCard(cardID, walletCard, callback)
	default:
		return "", fmt.Errorf("unsupported link type: %s", linkType)
	}
}

// createDataTransmitLink creates a data transmit link
func (c *Client) createDataTransmitLink(cardID string, cardData CardData, _ string) (string, error) {
	// Create CDATA token according to Samsung specification
	// This generates a JWT with Samsung-specific headers and 30-second expiration
	cdata, err := c.jwtManager.CreateDataTransmitToken(cardData)
	if err != nil {
		return "", fmt.Errorf("failed to create CDATA token: %v", err)
	}

	// Build ATW link according to Samsung Wallet API Guidelines
	// URL format: https://a.swallet.link/atw/v3/{cardId}#Clip?cdata={cdata}
	// cardId is the fixed identifier from Partners Portal (not individual card instance ID)
	atwURL := fmt.Sprintf("https://a.swallet.link/atw/v3/%s#Clip?cdata=%s", cardID, cdata)

	return atwURL, nil
}

// createDataFetchLink creates a data fetch link
func (c *Client) createDataFetchLink(cardID string, cardData CardData, _ string) (string, error) {
	// Data Fetch Link format: https://a.swallet.link/atw/v3/{certificateId}/{cardId}#Clip?pdata={pdata}
	// certificateId and cardId are fixed identifiers from Partners Portal
	// pdata is unique reference ID for this specific card instance

	// Generate a unique reference ID for this card instance
	// In a real implementation, this should be a secure, non-predictable ID
	refID := fmt.Sprintf("ref_%s_%d", cardData.CardID, time.Now().Unix())

	// Both certificateId and cardId must be obtained from Partners Portal
	if c.config.CertificateID == "" {
		return "", fmt.Errorf("certificate ID is required for data fetch links")
	}

	atwURL := fmt.Sprintf("https://a.swallet.link/atw/v3/%s/%s#Clip?pdata=%s",
		c.config.CertificateID, cardID, refID)

	return atwURL, nil
}

// createDataTransmitLinkFromWalletCard creates a data transmit link from WalletCard
func (c *Client) createDataTransmitLinkFromWalletCard(cardID string, walletCard WalletCard, _ string) (string, error) {
	// Create CDATA token according to Samsung specification
	cdata, err := c.jwtManager.CreateDataTransmitTokenFromWalletCard(walletCard)
	if err != nil {
		return "", fmt.Errorf("failed to create CDATA token: %v", err)
	}

	// Build ATW link according to Samsung Wallet API Guidelines
	atwURL := fmt.Sprintf("https://a.swallet.link/atw/v3/%s#Clip?cdata=%s", cardID, cdata)

	return atwURL, nil
}

// createDataFetchLinkFromWalletCard creates a data fetch link from WalletCard
func (c *Client) createDataFetchLinkFromWalletCard(cardID string, _ WalletCard, _ string) (string, error) {
	// Generate a unique reference ID for this card instance
	refID := fmt.Sprintf("ref_%d", time.Now().Unix())

	if c.config.CertificateID == "" {
		return "", fmt.Errorf("certificate ID is required for data fetch links")
	}

	atwURL := fmt.Sprintf("https://a.swallet.link/atw/v3/%s/%s#Clip?pdata=%s",
		c.config.CertificateID, cardID, refID)

	return atwURL, nil
}

// UpdateCard updates a wallet card
func (c *Client) UpdateCard(cardID string, cardData CardData, countryCode string) error {
	request := map[string]interface{}{
		"partner_id":   c.config.PartnerID,
		"card_id":      cardID,
		"card_data":    cardData,
		"country_code": countryCode,
	}

	_, err := c.makeAPIRequest("POST", fmt.Sprintf("%s/%s", pathUpdateCard, countryCode), request)
	return err
}

// CancelCard cancels wallet cards for a specific event
func (c *Client) CancelCard(eventID, reason string) error {
	request := CancelCardRequest{
		PartnerID: c.config.PartnerID,
		EventID:   eventID,
		Reason:    reason,
	}

	_, err := c.makeAPIRequest("POST", pathCancelCard, request)
	return err
}

// GetCardData retrieves card data
func (c *Client) GetCardData(cardID, countryCode string) (*CardData, error) {
	request := map[string]interface{}{
		"partner_id":   c.config.PartnerID,
		"card_id":      cardID,
		"country_code": countryCode,
	}

	response, err := c.makeAPIRequest("POST", fmt.Sprintf("%s/%s", pathGetCard, countryCode), request)
	if err != nil {
		return nil, err
	}

	var cardData CardData
	if err := json.Unmarshal(response, &cardData); err != nil {
		return nil, fmt.Errorf("failed to parse card data: %v", err)
	}

	return &cardData, nil
}

// HandleCallback handles the card state callback from Samsung Wallet
func (c *Client) HandleCallback(callbackData []byte) (*CardStateCallback, error) {
	var callback CardStateCallback
	if err := json.Unmarshal(callbackData, &callback); err != nil {
		return nil, fmt.Errorf("failed to parse callback data: %v", err)
	}

	// Validate the callback
	if callback.PartnerID != c.config.PartnerID {
		return nil, fmt.Errorf("invalid partner ID in callback")
	}

	return &callback, nil
}

// makeAPIRequest makes an HTTP request to Samsung Wallet API
func (c *Client) makeAPIRequest(method, path string, payload interface{}) ([]byte, error) {
	var body io.Reader
	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request payload: %v", err)
		}
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, c.baseURL+path, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Add authentication header if needed
	// This would typically be a JWT token or API key
	// Implementation depends on Samsung Wallet's authentication requirements

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %v", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode >= 400 {
		var apiError APIError
		if err := json.Unmarshal(responseBody, &apiError); err != nil {
			return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(responseBody))
		}
		return nil, &apiError
	}

	return responseBody, nil
}

// SetHTTPClient sets a custom HTTP client
func (c *Client) SetHTTPClient(client *http.Client) {
	c.httpClient = client
}

// GetJWTManager returns the JWT manager instance
func (c *Client) GetJWTManager() *JWTManager {
	return c.jwtManager
}

// SendUpdateNotification sends a card state update to Samsung Wallet Server API.
// This corresponds to the Update Notification API: POST /{cc2}/wltex/cards/{cardId}/updates
// Use this to update individual card states (e.g., DELETED, EXPIRED, SUSPENDED).
func (c *Client) SendUpdateNotification(cardID, cardType string, data []UpdateNotificationData, countryCode string) error {
	if cardID == "" {
		return fmt.Errorf("card ID is required")
	}
	if countryCode == "" {
		return fmt.Errorf("country code is required")
	}

	reqBody := UpdateNotificationRequest{
		Card: UpdateNotificationCard{
			Type: cardType,
			Data: data,
		},
	}

	return c.makeServerAPIRequest(
		http.MethodPost,
		fmt.Sprintf("%s/%s/wltex/cards/%s/updates", serverAPIBaseURL, countryCode, cardID),
		reqBody,
	)
}

// SendCancelNotification sends a cancel notification to Samsung Wallet Server API.
// This corresponds to the Cancel Notification API: POST /{cc2}/wltex/cards/{cardId}/cancels
// Use this to cancel all cards for a specific event (e.g., event cancellation).
func (c *Client) SendCancelNotification(cardID, cardType string, data []CancelNotificationData, countryCode string) error {
	if cardID == "" {
		return fmt.Errorf("card ID is required")
	}
	if countryCode == "" {
		return fmt.Errorf("country code is required")
	}

	reqBody := CancelNotificationRequest{
		Card: CancelNotificationCard{
			Type: cardType,
			Data: data,
		},
	}

	return c.makeServerAPIRequest(
		http.MethodPost,
		fmt.Sprintf("%s/%s/wltex/cards/%s/cancels", serverAPIBaseURL, countryCode, cardID),
		reqBody,
	)
}

// makeServerAPIRequest makes an authenticated HTTP request to Samsung Wallet Server API.
// Uses Bearer token (RS256 signed JWT) for authentication.
func (c *Client) makeServerAPIRequest(method, url string, payload interface{}) error {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request payload: %v", err)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// Generate Bearer token
	token, err := c.jwtManager.CreateServerAPIToken()
	if err != nil {
		return fmt.Errorf("failed to create server API token: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("x-smcs-partner-id", c.config.PartnerID)
	req.Header.Set("x-request-id", fmt.Sprintf("%d", time.Now().UnixNano()))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("API request failed with status %d: failed to read body: %v", resp.StatusCode, err)
		}
		var apiError APIError
		if err := json.Unmarshal(responseBody, &apiError); err != nil {
			return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(responseBody))
		}
		return &apiError
	}

	return nil
}

// Builder factory methods - automatically sets partnerID from client config

// NewEventTicket creates a new event ticket builder using official Samsung Wallet structure
func (c *Client) NewEventTicket(refID, title string) *EventTicketBuilder {
	return NewEventTicket(refID, title)
}
