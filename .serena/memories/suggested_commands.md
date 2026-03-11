# 開発コマンド一覧

## ビルド関連
| コマンド | 説明 |
|---------|------|
| `make build` | プロジェクトをビルド |
| `make build-example` | サンプルコードをビルド |
| `make run-example` | サンプルコードを実行 |

## テスト関連
| コマンド | 説明 |
|---------|------|
| `make test` | テスト実行 |
| `make test-coverage` | カバレッジ付きテスト実行 (coverage.html 生成) |
| `make test-race` | レースコンディション検出付きテスト |
| `make benchmark` | ベンチマークテスト |
| `make integration-test` | 統合テスト (API資格情報が必要) |

## コード品質
| コマンド | 説明 |
|---------|------|
| `make fmt` | コードフォーマット (go fmt) |
| `make vet` | go vet 実行 |
| `make lint` | golangci-lint 実行 |
| `make check` | 全チェック実行 (fmt, vet, lint, test) |

## 依存関係管理
| コマンド | 説明 |
|---------|------|
| `make deps` | 依存関係ダウンロード |
| `make mod` | go mod tidy 実行 |
| `make mod-verify` | モジュール整合性検証 |

## 開発環境
| コマンド | 説明 |
|---------|------|
| `make dev-setup` | 開発環境セットアップ (deps + install-tools) |
| `make install-tools` | 開発ツールインストール (golangci-lint) |
| `make generate-test-keys` | テスト用RSA鍵ペア生成 |
| `make watch` | ファイル変更監視 & 自動テスト (fswatch必要) |

## ドキュメント
| コマンド | 説明 |
|---------|------|
| `make doc` | ドキュメント生成 |
| `make doc-server` | ドキュメントサーバー起動 (:6060) |

## クリーンアップ
| コマンド | 説明 |
|---------|------|
| `make clean` | ビルド成果物削除 |

## プロファイリング
| コマンド | 説明 |
|---------|------|
| `make profile-cpu` | CPU プロファイル |
| `make profile-mem` | メモリプロファイル |

## リリース
| コマンド | 説明 |
|---------|------|
| `make tag VERSION=v1.0.0` | Git タグ作成 & プッシュ |

## システムユーティリティ (macOS/Darwin)
- `git`: バージョン管理
- `ls`: ファイル一覧
- `cd`: ディレクトリ移動
- `grep`: テキスト検索
- `find`: ファイル検索
- `openssl`: 鍵/証明書変換
