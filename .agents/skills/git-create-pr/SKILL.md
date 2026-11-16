---
name: git-create-pr
description: >
  mainやdev以外の作業ブランチからGitHub Pull Requestを作成する。
  Gitの差分とコミット履歴を分析し、英語のタイトルと本文を作成して
  GitHub CLI (gh) でPull Requestを作成する。
---

# Git Create PR Skill

## 1. Purpose

現在の作業ブランチの変更内容を分析し、GitHub Pull Requestを作成してください。

GitHub CLI (`gh`) を使用してPull Requestを作成してください。

Pull Requestのタイトルと本文は必ず英語で作成してください。

## 2. Execution Rules

- Git操作はターミナルから実行してください。
- GitHub操作には `gh` CLIを使用してください。
- 実際のGit差分とコミット履歴を根拠にしてください。
- Pull Requestを重複して作成しないでください。
- 既存のPull Requestを勝手に変更・クローズしないでください。
- ユーザーの変更を破棄してはいけません。
- エラーが発生した場合は処理を停止してください。

## 3. Workflow

### Step 1: Check Current Branch

現在のブランチを確認してください。

```bash
git branch --show-current
```

以下のブランチにいる場合は処理を終了してください。

- main
- dev

それ以外のブランチであることを確認してください。

### Step 2: Check Git Status

現在の作業状態を確認してください。

```bash
git status --short
```

未コミットの変更が存在する場合は、処理を停止してください。

ユーザーに変更をコミットするよう依頼してください。

### Step 3: Check Remote Branch

リモート情報を更新してください。

```bash
git fetch origin
```

現在のブランチに対応するリモートブランチが存在することを確認してください。

存在しない場合は、ユーザーにpushを依頼してください。

### Step 4: Determine Base Branch

Pull Requestの作成先となるベースブランチを決定してください。

候補:

- main
- dev

作業ブランチの作成元となったブランチを確認してください。

確認には以下を利用してください。

- Gitのコミット履歴
- merge-base
- リモートブランチ
- GitHubのブランチ情報

作成元を明確に特定できない場合は、推測で決定せずユーザーに確認してください。

Pull Requestの方向は以下です。

```text
features/* -> main
features/* -> dev
```

### Step 5: Check Existing Pull Request

同じheadブランチとbaseブランチのPull Requestが既に存在しないか確認してください。

```bash
gh pr list --head <head-branch> --base <base-branch> --state open
```

既に存在する場合は、新しいPull Requestを作成しないでください。

既存Pull RequestのURLと状態をユーザーに報告して終了してください。

### Step 6: Analyze Changes

以下のコマンドを実行して変更内容を確認してください。

```bash
git log --oneline <base-branch>..HEAD
git diff --stat <base-branch>...HEAD
git diff <base-branch>...HEAD
```

確認する項目:

- 変更したファイル
- 追加・削除したファイル
- 実装した機能
- 修正したバグ
- リファクタリングの内容
- テストの追加・変更
- ドキュメントの変更

変更内容を十分に理解できない場合は、推測でPRを作成しないでください。

### Step 7: Generate PR Title

変更内容に基づいて、英語のPRタイトルを作成してください。

形式:

```text
<type>(<scope>): <summary>
```

typeは以下から適切なものを選択してください。

| Type | Description |
|---|---|
| feat | 新機能 |
| fix | バグ修正 |
| docs | ドキュメント |
| refactor | リファクタリング |
| test | テスト |
| chore | 設定・雑務 |
| perf | パフォーマンス改善 |
| build | ビルド関連 |
| ci | CI/CD関連 |

scopeは変更対象を簡潔に表してください。

### Step 8: Generate PR Description

PR本文は英語で作成してください。

以下の形式を基本とします。

```markdown
## Summary

- Describe the main changes.
- Explain the purpose of the changes.

## Changes

- List the implemented changes.
- List important modifications.

## Testing

- Describe the tests performed.
- Mention tests that have not been performed, if applicable.

## Related Issues

- Closes #<issue-number>
```

ルール:

- 実際の変更内容だけを記述する。
- 実施していないテストを実施済みと記述しない。
- 関連Issueが確認できない場合はRelated Issuesを省略する。
- Issue番号はブランチ名などから取得できる場合に利用する。
- Issueとの関連性が不明な場合は推測しない。

### Step 9: Create Pull Request

以下の形式でPRを作成してください。

```bash
gh pr create \
  --base <base-branch> \
  --head <head-branch> \
  --title "<PR title>" \
  --body "<PR description>"
```

複数行の本文は、安全な引数渡しを利用してください。

必要に応じて一時ファイルを作成し、`--body-file` を利用しても構いません。

PR作成に失敗した場合は処理を停止してください。

### Step 10: Verify Pull Request

作成したPRを確認してください。

```bash
gh pr view --json number,title,url,state,baseRefName,headRefName
```

以下を確認してください。

- PRが作成されていること。
- タイトルが正しいこと。
- baseブランチが正しいこと。
- headブランチが正しいこと。
- PRのURLが取得できること。

## 4. Error Handling

GitまたはGitHub CLIの操作に失敗した場合:

1. 失敗したコマンドを表示する。
2. エラーメッセージを表示する。
3. 原因が分かる場合は説明する。
4. ユーザーに必要な操作を依頼する。
5. それ以上の操作を行わず終了する。

禁止事項:

- エラーを無視した処理の続行。
- ユーザーの変更を破棄する操作。
- 既存PRの勝手な編集・クローズ。
- ベースブランチを推測してPRを作成すること。
- PR作成失敗時の強引な再作成。

## 5. Completion Report

正常終了した場合、以下を報告してください。

- PR番号
- PRタイトル
- PR URL
- headブランチ
- baseブランチ
- 変更内容の概要
- テスト結果

