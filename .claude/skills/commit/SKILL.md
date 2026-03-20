---
name: commit
description: Conventional Commits 規約に従ってコミットを作成する
---

# Commit Skill

変更内容を分析し、Conventional Commits 規約に従ったコミットメッセージを作成してコミットする。

## 手順

1. `git status` と `git diff --staged` および `git diff` を並列実行し、現在の変更状態を確認する
2. 変更内容を分析し、以下の Conventional Commits ルールに従ってコミットメッセージを作成する
3. 関連するファイルをステージングし、コミットを実行する
4. コミット後に `git status` で成功を確認する

## コミットメッセージのフォーマット

```
<type>(<scope>): <description>

<body>

<footer>
```

- **type**（必須）: 下記の type 一覧から選択
- **scope**（任意）: 変更対象のモジュールやディレクトリ名
- **description**（必須）: 変更の要約を日本語で簡潔に記述
- **body**（任意）: 変更の詳細な説明を日本語で記述。なぜその変更が必要だったかを書く
- **footer**（任意）: `Refs: #123` や `BREAKING CHANGE:` など

### type 一覧

| type | 用途 |
|---|---|
| `feat` | 新機能 |
| `fix` | バグ修正 |
| `docs` | ドキュメントのみの変更 |
| `style` | コードの意味に影響しない変更（空白、フォーマットなど） |
| `refactor` | バグ修正でも機能追加でもないコード変更 |
| `perf` | パフォーマンス改善 |
| `test` | テストの追加・修正 |
| `chore` | ビルドプロセスや補助ツールの変更 |
| `ci` | CI設定の変更 |
| `build` | ビルドシステムや外部依存の変更 |

### Breaking Change

破壊的変更がある場合は以下のいずれかで示す:
- type の後に `!` を付与: `feat!: 認証方式を変更`
- footer に `BREAKING CHANGE: 説明` を記述

### 例

```
feat(auth): ログイン画面にSSO認証を追加

SAML 2.0 ベースのSSO認証フローを実装。
既存のパスワード認証と併用可能。

Refs: #123
```

```
fix: トークン更新時のレースコンディションを修正
```

## Co-Authored-By

コミットメッセージの末尾に以下を付与する:

```
Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>
```

## 注意事項

- `.env` やクレデンシャルなど機密情報を含むファイルはコミットしない。該当ファイルがある場合はユーザーに警告する
- `git add -A` ではなく、関連ファイルを個別に指定してステージングする
- pre-commit hook が失敗した場合は問題を修正し、amend ではなく新しいコミットを作成する
- コミットメッセージは HEREDOC で渡す
