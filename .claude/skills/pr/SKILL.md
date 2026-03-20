---
name: pr
description: Pull Requestを作成してGitHubへ送信する
---

# Pull Request Skill

現在のブランチの変更内容を分析し、GitHub Pull Request を作成する。

## 手順

1. 以下のコマンドを並列実行し、現在の状態を把握する
   - `git status`: 未コミットの変更がないか確認
   - `git log main..HEAD --oneline`: メインブランチからの全コミットを確認
   - `git diff main...HEAD`: メインブランチからの全変更差分を確認
   - `git rev-list HEAD --not --remotes`: リモート未プッシュのコミットがあるか確認
   - `gh pr view --json url,title,body 2>/dev/null`: 既存PRの有無を確認
2. 未コミットの変更がある場合はユーザーに警告し、先にコミットするか確認する
3. リモート未プッシュのコミットがある場合はユーザーに警告し、先にプッシュするか確認する
4. 全コミットと差分を分析し、PRタイトルと本文を作成する（既存PRがある場合、タイトルは不要）
5. **既存PRがある場合**: `gh pr edit` でPR本文のみを更新する
6. **既存PRがない場合**: `gh pr create` で新規PRを作成する
7. PRのURLをユーザーに報告する

## PRタイトル

- 70文字以内に収める
- 日本語で簡潔に変更内容を要約する
- Conventional Commits の type プレフィックスがコミットに含まれる場合は、それに合わせたタイトルにする

## PR本文のフォーマット

```
## 概要
- 変更の目的・背景を箇条書きで1〜3行

## 変更内容
- 主な変更点を箇条書きで記述

## テスト計画
- [ ] テスト方法や確認手順のチェックリスト

🤖 Generated with [Claude Code](https://claude.com/claude-code)
```

## 既存PRの確認

手順ステップ1で `gh pr view --json url,title,body 2>/dev/null` を並列実行して確認する。

- 成功した場合: 既存PRがあるので `gh pr edit` で本文を更新する
- 失敗した場合: PRが存在しないので `gh pr create` で新規作成する

## 新規PR作成（`gh pr create`）

HEREDOC を使って本文を渡す:

```bash
gh pr create --title "PRタイトル" --body "$(cat <<'EOF'
## 概要
- ...

## 変更内容
- ...

## テスト計画
- [ ] ...

🤖 Generated with [Claude Code](https://claude.com/claude-code)
EOF
)"
```

## 既存PR更新（`gh pr edit`）

既存PRがある場合は本文のみを更新する。タイトルは変更しない。

```bash
gh pr edit --body "$(cat <<'EOF'
## 概要
- ...

## 変更内容
- ...

## テスト計画
- [ ] ...

🤖 Generated with [Claude Code](https://claude.com/claude-code)
EOF
)"
```

## 注意事項

- ベースブランチはデフォルトで `main` を使用する。ユーザーが指定した場合はそれに従う
- 最新コミットだけでなく、ブランチ上の全コミットの変更を分析してPR内容を作成する
- ドラフトPRを作成したい場合はユーザーが `/pr --draft` のように指定できる。引数に `--draft` が含まれる場合は `gh pr create` に `--draft` を追加する
