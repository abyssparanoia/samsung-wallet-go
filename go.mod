module github.com/abyssparanoia/samsung-wallet-go

godebug x509negativeserial=1

go 1.24.0

require (
	github.com/go-jose/go-jose/v3 v3.0.4
	github.com/golang-jwt/jwt/v5 v5.3.0
	github.com/google/uuid v1.6.0
)

require (
	github.com/stretchr/testify v1.10.0 // indirect
	golang.org/x/crypto v0.39.0 // indirect
)
