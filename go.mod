module github.com/abyssparanoia/samsung-wallet-go

godebug x509negativeserial=1

go 1.24.0

require (
	github.com/go-jose/go-jose/v3 v3.0.4
	github.com/go-jose/go-jose/v4 v4.1.1
	github.com/golang-jwt/jwt/v5 v5.3.0
	github.com/google/uuid v1.6.0
)

require golang.org/x/crypto v0.39.0 // indirect
