# コードスタイルと規約

## フォーマット
- **フォーマッター**: `go fmt` および `gofumpt`
- **インポート整理**: `goimports` (ローカルプレフィックス: `github.com/abyssparanoia/samsung-wallet-go`)
- **行の長さ**: 最大120文字

## 命名規則
- **パッケージ名**: 小文字、単一単語 (例: `wallet`)
- **エクスポート型/関数**: PascalCase (例: `NewClient`, `CardData`)
- **プライベート型/関数**: camelCase (例: `jwtManager`, `makeAPIRequest`)
- **定数**: PascalCase (例: `CardTypeEventTicket`)
- **略語**: 既知の略語は大文字維持 (例: `ID`, `URL`, `JWT`, `API`)

## 構造体定義
```go
type Config struct {
    PartnerID         string `json:"partner_id"`          // フィールドの説明
    PartnerPrivateKey string `json:"partner_private_key"` // インラインコメント
    BaseURL           string `json:"base_url,omitempty"`  // omitempty for optional
}
```
- JSONタグを使用
- コメントはフィールドの横にインラインで記述
- オプショナルフィールドには `omitempty` を付与

## エラーハンドリング
```go
if config == nil {
    return nil, fmt.Errorf("config cannot be nil")
}

if err != nil {
    return nil, fmt.Errorf("failed to create JWT manager: %v", err)
}
```
- `fmt.Errorf` でエラーメッセージを構築
- エラーメッセージは小文字で開始
- コンテキストを含める (例: "failed to create JWT manager:")

## リンター設定 (.golangci.yml)
有効なリンター:
- `errcheck`: エラーチェック
- `govet`: 静的解析
- `staticcheck`: 追加静的解析
- `gocyclo`: 循環複雑度 (最大15)
- `gosec`: セキュリティチェック
- `dupl`: 重複コード検出
- `goconst`: 定数化すべき値
- `gocognit`: 認知複雑度
- `misspell`: スペルチェック (US英語)
- `lll`: 行長チェック
- `prealloc`: スライス事前確保

除外設定:
- テストファイル (`_test.go`) は一部リンター除外
- `examples/` ディレクトリは errcheck, gosec 除外
- `err` 変数のシャドウイングは許可

## ドキュメント
- エクスポート型/関数にはGoDocコメントを付ける
- 関数コメントは関数名で始める: `// NewClient creates a new Samsung Wallet client`

## テスト
- テストファイル: `*_test.go`
- テスト関数: `func TestXxx(t *testing.T)`
- テーブル駆動テストを推奨
