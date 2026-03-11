# Samsung Wallet Go SDK - プロジェクト概要

## 目的
Samsung Wallet "Add to Samsung Wallet" (ATW) サービス統合用の Go SDK。Samsung Walletにカードを追加するためのCDATA生成、リンク生成、サーバーAPIクライアント、コールバック処理、カード状態管理を提供する。

## サポートするカードタイプ
- Event Ticket (イベントチケット)
- Boarding Pass (搭乗券)
- Coupon (クーポン)
- Gift Card (ギフトカード)
- Loyalty Card (ロイヤリティカード)

## 主要機能
- **CDATA Generation**: Samsung公式仕様に準拠 (JWE暗号化 + JWS署名)
- **30秒トークン有効期限**: セキュリティ準拠のトークン有効期限
- **Samsung固有JWTヘッダー**: `cty`, `partnerId`, `ver`, `certificateId`, `utc` ヘッダー
- Add to Samsung Wallet リンク生成
- Server API クライアント
- コールバック処理
- カード状態管理

## テックスタック
- **言語**: Go 1.24.0
- **JWT/JWE**:
  - github.com/go-jose/go-jose/v3 (JWE暗号化)
  - github.com/golang-jwt/jwt/v5 (JWT署名)
- **ユーティリティ**:
  - github.com/google/uuid (UUID生成)
- **リンター**: golangci-lint

## コードベース構造
```
.
├── wallet/                 # メインSDKコード
│   ├── client.go          # クライアント実装
│   ├── types.go           # 型定義
│   ├── cards.go           # カードビルダー
│   └── jwt.go             # JWT/JWE処理
├── examples/               # サンプルコード
│   └── event-ticket/      # イベントチケットのサンプル
├── bin/                   # ビルド出力
├── secret/                # 秘密鍵（gitignore）
├── Makefile               # ビルド/テストコマンド
├── .golangci.yml          # リンター設定
└── go.mod                 # Go モジュール定義
```

## 環境変数
- `SAMSUNG_WALLET_PARTNER_ID`: パートナーID
- `SAMSUNG_WALLET_PARTNER_PRIVATE_KEY`: RSA秘密鍵
- `SAMSUNG_WALLET_SAMSUNG_PUBLIC_KEY`: Samsung公開鍵
- `SAMSUNG_WALLET_CERTIFICATE_ID`: 4桁英数字の証明書ID
- `SAMSUNG_WALLET_EVENT_TICKET_CARD_ID`: イベントチケットカードID
