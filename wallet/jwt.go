package wallet

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// CDATAClaims represents the claims for Samsung Wallet CDATA JWT
type CDATAClaims struct {
	PartnerId     string `json:"partnerId"`
	Ver           string `json:"ver"`
	CertificateId string `json:"certificateId"`
	UTC           int64  `json:"utc"`
	jwt.RegisteredClaims
}

// JWTManager handles JWT operations for Samsung Wallet
type JWTManager struct {
	partnerPrivateKey *rsa.PrivateKey
	samsungPublicKey  *rsa.PublicKey
	partnerID         string // Changed from serviceID to match Samsung naming
	certificateID     string
}

// NewJWTManager creates a new JWT manager
func NewJWTManager(partnerPrivateKeyPEM string, samsungPublicKeyPEM string, serviceID string, certificateID string) (*JWTManager, error) {
	// Parse partner private key
	partnerPrivateKey, err := parsePrivateKey(partnerPrivateKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to parse partner private key: %v", err)
	}

	// Parse Samsung public key
	samsungPublicKey, err := parsePublicKey(samsungPublicKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Samsung public key: %v", err)
	}

	return &JWTManager{
		partnerPrivateKey: partnerPrivateKey,
		samsungPublicKey:  samsungPublicKey,
		partnerID:         serviceID, // Changed field name
		certificateID:     certificateID,
	}, nil
}

// CreateCDATA creates CDATA token for Samsung Wallet following the official specification
// This implements the two-step process: JWE encryption + JWS signing
func (j *JWTManager) CreateCDATA(cardData interface{}) (string, error) {
	// Step 1: Convert card data to JSON
	cardDataJSON, err := json.Marshal(cardData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal card data: %v", err)
	}

	// Step 2: JWE Encryption with Samsung public key
	// Create JWE encrypter using Samsung's public key
	encrypter, err := jose.NewEncrypter(
		jose.A128GCM,
		jose.Recipient{
			Algorithm: jose.RSA1_5,
			Key:       j.samsungPublicKey,
		},
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("failed to create JWE encrypter: %v", err)
	}

	// Encrypt the card data
	jweObject, err := encrypter.Encrypt(cardDataJSON)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt card data: %v", err)
	}

	// Serialize the JWE object
	jwePayload, err := jweObject.CompactSerialize()
	if err != nil {
		return "", fmt.Errorf("failed to serialize JWE object: %v", err)
	}

	// Step 3: JWS Signing with partner private key
	now := time.Now()
	utc := now.UnixMilli() // UTC timestamp in milliseconds

	// Create JWS signer using partner private key
	signer, err := jose.NewSigner(
		jose.SigningKey{
			Algorithm: jose.RS256,
			Key:       j.partnerPrivateKey,
		},
		(&jose.SignerOptions{}).WithType("JWT").WithHeader("cty", "CARD").
			WithHeader("partnerId", j.partnerID).
			WithHeader("ver", "3").
			WithHeader("certificateId", j.certificateID).
			WithHeader("utc", utc),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create JWS signer: %v", err)
	}

	// Sign the JWE payload with JWS
	// The payload should be the JWE encrypted data
	jwsObject, err := signer.Sign([]byte(jwePayload))
	if err != nil {
		return "", fmt.Errorf("failed to sign with JWS: %v", err)
	}

	// Serialize the final token
	tokenString, err := jwsObject.CompactSerialize()
	if err != nil {
		return "", fmt.Errorf("failed to serialize JWS object: %v", err)
	}

	return tokenString, nil
}

// CreateDataTransmitToken creates a token for data transmit link
func (j *JWTManager) CreateDataTransmitToken(cardData interface{}) (string, error) {
	// Data transmit uses CDATA format with 30-second expiration
	return j.CreateCDATA(cardData)
}

// CreateDataFetchToken creates a token for data fetch link
func (j *JWTManager) CreateDataFetchToken(cardData interface{}) (string, error) {
	// Data fetch also uses CDATA format with 30-second expiration
	return j.CreateCDATA(cardData)
}

// CreateDataTransmitTokenFromWalletCard creates a token for data transmit link from WalletCard
func (j *JWTManager) CreateDataTransmitTokenFromWalletCard(walletCard WalletCard) (string, error) {
	return j.CreateCDATA(walletCard)
}

// CreateDataFetchTokenFromWalletCard creates a token for data fetch link from WalletCard
func (j *JWTManager) CreateDataFetchTokenFromWalletCard(walletCard WalletCard) (string, error) {
	return j.CreateCDATA(walletCard)
}

// VerifyToken verifies a Samsung Wallet token
func (j *JWTManager) VerifyToken(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return &j.partnerPrivateKey.PublicKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// TokenInfo represents information extracted from a token
type TokenInfo struct {
	ServiceID     string    `json:"service_id"`
	CertificateID string    `json:"certificate_id"`
	Version       string    `json:"version"`
	IssuedAt      time.Time `json:"issued_at"`
	ExpiresAt     time.Time `json:"expires_at"`
	TokenID       string    `json:"token_id"`
	Valid         bool      `json:"valid"`
	UTC           int64     `json:"utc"`
}

// GetTokenInfo extracts information from a token without full verification
func (j *JWTManager) GetTokenInfo(tokenString string) (*TokenInfo, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims type")
	}

	info := &TokenInfo{}

	// Extract from headers
	if partnerId, ok := token.Header["partnerId"].(string); ok {
		info.ServiceID = partnerId
	}
	if certId, ok := token.Header["certificateId"].(string); ok {
		info.CertificateID = certId
	}
	if ver, ok := token.Header["ver"].(string); ok {
		info.Version = ver
	}
	if utc, ok := token.Header["utc"].(float64); ok {
		info.UTC = int64(utc)
	}

	// Extract from claims
	if iat, ok := claims["iat"].(float64); ok {
		info.IssuedAt = time.Unix(int64(iat), 0)
	}
	if exp, ok := claims["exp"].(float64); ok {
		info.ExpiresAt = time.Unix(int64(exp), 0)
		info.Valid = time.Now().Before(info.ExpiresAt)
	}
	if jti, ok := claims["jti"].(string); ok {
		info.TokenID = jti
	}

	return info, nil
}

// CreateCallbackToken creates a token for callback verification
func (j *JWTManager) CreateCallbackToken(callback CardStateCallback) (string, error) {
	now := time.Now()

	claims := jwt.MapClaims{
		"partner_id":   callback.PartnerID, // Changed from service_id
		"card_id":      callback.CardID,
		"event":        callback.Event,
		"country_code": callback.CountryCode,
		"timestamp":    callback.Timestamp.Unix(),
		"iat":          now.Unix(),
		"exp":          now.Add(time.Hour).Unix(), // 1 hour expiration for callbacks
		"jti":          uuid.New().String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(j.partnerPrivateKey)
}

// VerifyCallbackToken verifies a callback token and extracts the callback data
func (j *JWTManager) VerifyCallbackToken(tokenString string) (*CardStateCallback, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return &j.partnerPrivateKey.PublicKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		callback := &CardStateCallback{}

		if partnerID, ok := claims["partner_id"].(string); ok {
			callback.PartnerID = partnerID
		}
		if cardID, ok := claims["card_id"].(string); ok {
			callback.CardID = cardID
		}
		if event, ok := claims["event"].(string); ok {
			callback.Event = CardState(event)
		}
		if countryCode, ok := claims["country_code"].(string); ok {
			callback.CountryCode = countryCode
		}
		if timestamp, ok := claims["timestamp"].(float64); ok {
			callback.Timestamp = time.Unix(int64(timestamp), 0)
		}

		return callback, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// parsePrivateKey parses a PEM-encoded RSA private key
func parsePrivateKey(privateKeyPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil || (block.Type != "PRIVATE KEY" && block.Type != "RSA PRIVATE KEY") {
		return nil, fmt.Errorf("invalid private key format")
	}

	var privateKey *rsa.PrivateKey
	var err error

	if block.Type == "RSA PRIVATE KEY" {
		// PKCS#1 format
		privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	} else {
		// PKCS#8 format
		key, parseErr := x509.ParsePKCS8PrivateKey(block.Bytes)
		if parseErr != nil {
			return nil, parseErr
		}
		var ok bool
		privateKey, ok = key.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("not an RSA private key")
		}
	}

	return privateKey, err
}

// parsePublicKey parses a PEM-encoded RSA public key
func parsePublicKey(publicKeyPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return nil, fmt.Errorf("invalid public key format")
	}

	var publicKey *rsa.PublicKey
	var err error

	switch block.Type {
	case "PUBLIC KEY":
		// X.509 format
		key, parseErr := x509.ParsePKIXPublicKey(block.Bytes)
		if parseErr != nil {
			return nil, parseErr
		}
		var ok bool
		publicKey, ok = key.(*rsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("not an RSA public key")
		}
	case "RSA PUBLIC KEY":
		// PKCS#1 format
		publicKey, err = x509.ParsePKCS1PublicKey(block.Bytes)
	case "CERTIFICATE":
		// X.509 certificate
		cert, parseErr := x509.ParseCertificate(block.Bytes)
		if parseErr != nil {
			return nil, parseErr
		}
		var ok bool
		publicKey, ok = cert.PublicKey.(*rsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("certificate does not contain an RSA public key")
		}
	default:
		return nil, fmt.Errorf("unsupported public key type: %s", block.Type)
	}

	return publicKey, err
}
