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
	pathATWLink    = "/v1/wallet/atw/link"
	pathUpdateCard = "/v1/wallet/card/update"
	pathCancelCard = "/v1/wallet/card/cancel"
	pathGetCard    = "/v1/wallet/card/get"
	pathCardState  = "/v1/wallet/card/state"
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
		linkType = "data_transmit" // default to data transmit
	}

	var callback string
	if len(callbackURL) > 0 {
		callback = callbackURL[0]
	}

	switch linkType {
	case "data_transmit":
		return c.createDataTransmitLink(cardID, cardData, callback)
	case "data_fetch":
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
		linkType = "data_transmit" // default to data transmit
	}

	var callback string
	if len(callbackURL) > 0 {
		callback = callbackURL[0]
	}

	switch linkType {
	case "data_transmit":
		return c.createDataTransmitLinkFromWalletCard(cardID, walletCard, callback)
	case "data_fetch":
		return c.createDataFetchLinkFromWalletCard(cardID, walletCard, callback)
	default:
		return "", fmt.Errorf("unsupported link type: %s", linkType)
	}
}

// createDataTransmitLink creates a data transmit link
func (c *Client) createDataTransmitLink(cardID string, cardData CardData, callbackURL string) (string, error) {
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
func (c *Client) createDataFetchLink(cardID string, cardData CardData, callbackURL string) (string, error) {
	// Data Fetch Link format: https://a.swallet.link/atw/v3/{certificateId}/{cardId}#Clip?pdata={pdata}
	// certificateId and cardId are fixed identifiers from Partners Portal
	// pdata is unique reference ID for this specific card instance

	// Generate a unique reference ID for this card instance
	// In a real implementation, this should be a secure, non-predictable ID
	refId := fmt.Sprintf("ref_%s_%d", cardData.CardID, time.Now().Unix())

	// Both certificateId and cardId must be obtained from Partners Portal
	if c.config.CertificateID == "" {
		return "", fmt.Errorf("certificate ID is required for data fetch links")
	}

	atwURL := fmt.Sprintf("https://a.swallet.link/atw/v3/%s/%s#Clip?pdata=%s",
		c.config.CertificateID, cardID, refId)

	return atwURL, nil
}

// createDataTransmitLinkFromWalletCard creates a data transmit link from WalletCard
func (c *Client) createDataTransmitLinkFromWalletCard(cardID string, walletCard WalletCard, callbackURL string) (string, error) {
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
func (c *Client) createDataFetchLinkFromWalletCard(cardID string, walletCard WalletCard, callbackURL string) (string, error) {
	// Generate a unique reference ID for this card instance
	refId := fmt.Sprintf("ref_%d", time.Now().Unix())

	if c.config.CertificateID == "" {
		return "", fmt.Errorf("certificate ID is required for data fetch links")
	}

	atwURL := fmt.Sprintf("https://a.swallet.link/atw/v3/%s/%s#Clip?pdata=%s",
		c.config.CertificateID, cardID, refId)

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
func (c *Client) CancelCard(eventID string, reason string) error {
	request := CancelCardRequest{
		PartnerID: c.config.PartnerID,
		EventID:   eventID,
		Reason:    reason,
	}

	_, err := c.makeAPIRequest("POST", pathCancelCard, request)
	return err
}

// GetCardData retrieves card data
func (c *Client) GetCardData(cardID string, countryCode string) (*CardData, error) {
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

// Builder factory methods - automatically sets partnerID from client config

// NewEventTicket creates a new event ticket builder using official Samsung Wallet structure
func (c *Client) NewEventTicket(refID, title string) *EventTicketBuilder {
	return NewEventTicket(refID, title)
}
