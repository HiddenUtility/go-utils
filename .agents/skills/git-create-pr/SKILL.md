---
name: git-create-pr
description: >
  mainやdev以外の作業ブランチからGitHub Pull Requestを作成する。
  ユーザーからの直接呼び出しと、他のスキルからの呼び出しの両方に対応する。
  Gitの差分とコミット履歴を分析し、英語のタイトルと本文を作成して
  GitHub CLI (gh) でPull Requestを作成する。
---

# Git Create PR Skill

## 1. Purpose

現在の作業ブランチの変更内容を分析し、GitHub Pull Requestを作成してください。

GitHub CLI (`gh`) を使用してPull Requestを作成してください。

Pull Requestのタイトルと本文は必ず英語で作成してください。

このスキルは、以下の2つの呼び出し方法に対応してください。

1. ユーザーからの直接呼び出し
2. 他のスキルからの呼び出し

呼び出し元に応じて、必要な情報の取得方法と結果の返却方法を切り替えてください。

## 2. Execution Rules

- Git操作はターミナルから実行してください。
- GitHub操作には `gh` CLIを使用してください。
- 実際のGit差分とコミット履歴を根拠にしてください。
- Pull Requestを重複して作成しないでください。
- 既存のPull Requestを勝手に変更・クローズしないでください。
- ユーザーの変更を破棄してはいけません。
- エラーが発生した場合は処理を停止してください。
- Pull Requestのタイトルと本文は必ず英語で作成してください。
- 呼び出し元から受け取った情報を優先し、同じ内容をユーザーに再質問しないでください。
- 呼び出し元から指定されたbaseブランチを勝手に変更しないでください。
- baseブランチを特定できない場合は、推測で決定してはいけません。
- Pull Request作成後に、自動マージや自動クローズを行ってはいけません。
- 既存のPull Requestが見つかった場合は、新規作成せず、その情報を報告してください。

## 3. Determine Invocation Type

最初に、呼び出し方法を判定してください。

### 3.1 User Invocation

ユーザーが直接このスキルを呼び出した場合は、以下の情報を確認してください。

- 現在の作業ブランチ
- Pull Requestの作成先となるbaseブランチ（任意）
- 関連するIssue番号（任意）

baseブランチが指定されていない場合は、Gitの履歴などから特定してください。

関連Issue番号はブランチ名やコミット履歴から取得できる場合があります。取得できない場合は任意でユーザーに確認してください。

必要な情報が不足している場合のみ、ユーザーに質問してください。

### 3.2 Skill Invocation

他のスキルから呼び出された場合は、呼び出し元から渡された情報を使用してください。

受け取る情報:

- Task Description (任意)
  - 実施した作業の概要
- Head Branch (任意)
  - Pull Requestの作成元ブランチ
- Base Branch (任意)
  - Pull Requestの作成先ブランチ
- Issue Number (任意)
  - 関連するGitHub Issue番号
- Additional Instructions (任意)
  - タイトルや本文に関する追加指示

以下のルールに従ってください。

1. Head Branchが指定されている場合は、そのブランチを使用する。
2. Head Branchが指定されていない場合は、現在のブランチを使用する。
3. Base Branchが指定されている場合は、そのブランチを使用する。
4. Base Branchが指定されていない場合は、Gitの履歴などから特定する。
5. Issue Numberが指定されている場合は、その番号を使用する。
6. Task Descriptionが指定されている場合は、変更内容の理解に利用する。
7. Additional Instructionsが指定されている場合は、その指示に従う。
8. 呼び出し元から受け取った情報をユーザーに再質問しない。
9. 必要な情報が不足しており、安全に処理できない場合のみ、呼び出し元に確認する。
10. PR作成結果を呼び出し元に返す。

呼び出し元が不明な場合は、ユーザーからの直接呼び出しとして扱ってください。

## 4. Workflow

### Step 1: Check Current Branch

現在のブランチを確認してください。

```bash
git branch --show-current
```

以下のブランチにいる場合は処理を終了してください。

- main
- dev

それ以外のブランチであることを確認してください。

Skill InvocationでHead Branchが指定されている場合は、指定されたブランチと現在のブランチが一致することを確認してください。

一致しない場合は処理を停止し、呼び出し元に報告してください。

### Step 2: Check Git Status

現在の作業状態を確認してください。

```bash
git status --short
```

未コミットの変更が存在する場合は、処理を停止してください。

ユーザーに変更をコミットするよう依頼してください。

変更を勝手にコミット、stash、破棄してはいけません。

### Step 3: Check Remote Branch

リモート情報を更新してください。

```bash
git fetch origin
```

現在のブランチに対応するリモートブランチが存在することを確認してください。

存在しない場合は、ユーザーまたは呼び出し元にpushを依頼してください。

このスキル内で勝手にpushしてはいけません。

### Step 4: Determine Base Branch

Pull Requestの作成先となるbaseブランチを決定してください。

候補:

- main
- dev

以下の優先順位で決定してください。

1. 呼び出し元から指定されたBase Branch
2. ユーザーが指定したBase Branch
3. Gitのコミット履歴などから特定できる作業ブランチの作成元

作成元の確認には以下を利用してください。

- Gitのコミット履歴
- `git merge-base`
- リモートブランチ
- GitHubのブランチ情報

例:

```bash
git log --oneline --decorate --graph --all
git merge-base <candidate-base> <head-branch>
git branch -r
```

作業ブランチの作成元を明確に特定できない場合は、推測で決定せずユーザーまたは呼び出し元に確認してください。

Base Branchが指定された場合は、その指定を優先してください。

Pull Requestの方向は以下を基本とします。

```text
<working-branch> -> main
<working-branch> -> dev
```

ただし、呼び出し元またはユーザーが別のbaseブランチを明示的に指定した場合は、その指定に従ってください。

### Step 5: Check Existing Pull Request

同じheadブランチとbaseブランチのPull Requestが既に存在しないか確認してください。

```bash
gh pr list --head <head-branch> --base <base-branch> --state open
```

既に存在する場合は、新しいPull Requestを作成しないでください。

既存Pull Requestの以下の情報を取得してください。

- PR番号
- PRタイトル
- PR URL
- PRの状態
- headブランチ
- baseブランチ

既存Pull Requestの情報をユーザーまたは呼び出し元に返して終了してください。

既存Pull Requestを勝手に編集・クローズしてはいけません。

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

Task Descriptionが指定されている場合は、実際の差分と照合してください。

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

scopeが不要な場合は省略して構いません。

Additional Instructionsでタイトルの形式が指定されている場合は、その指定を優先してください。

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
- Issue番号は呼び出し元、ユーザー、ブランチ名、コミット履歴などから取得できる場合に利用する。
- Issueとの関連性が不明な場合は推測しない。
- Additional Instructionsがある場合は、その指示を反映する。
- PR本文は読みやすく簡潔にまとめる。

PR作成前に、生成したタイトルと本文を確認してください。

ユーザーからの直接呼び出しの場合は、作成するタイトルと本文をユーザーに提示してください。

スキルからの呼び出しの場合は、呼び出し元にタイトルと本文を提示してください。

### Step 9: Create Pull Request

以下の形式でPRを作成してください。

```bash
gh pr create \
  --base <base-branch> \
  --head <head-branch> \
  --title "<PR title>" \
  --body-file <body-file>
```

複数行の本文は、安全な引数渡しを利用してください。

必要に応じて一時ファイルを作成し、`--body-file` を利用してください。

一時ファイルを使用した場合は、PR作成後に不要な一時ファイルを削除してください。

PR作成に失敗した場合は処理を停止してください。

失敗後に同じコマンドを無条件に再実行してはいけません。

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

すべて確認できた場合のみ正常終了としてください。

## 5. Error Handling

GitまたはGitHub CLIの操作に失敗した場合:

1. 失敗したコマンドを表示する。
2. エラーメッセージを表示する。
3. 原因が分かる場合は説明する。
4. ユーザーに必要な操作を依頼する。
5. それ以上の操作を行わず終了する。
6. 呼び出し元がスキルの場合は、失敗した処理とエラー内容を呼び出し元にも返す。

禁止事項:

- エラーを無視した処理の続行。
- ユーザーの変更を破棄する操作。
- 既存PRの勝手な編集・クローズ。
- baseブランチを推測してPRを作成すること。
- PR作成失敗時の強引な再作成。
- このスキル内での勝手なcommitやpush。
- 自動マージ。
- PR作成後の自動クローズ。

## 6. Result

処理結果には以下の情報を含めてください。

- Status: SUCCESS / FAILED / EXISTING
- PR Number: 作成済みまたは既存のPR番号
- PR Title: PRタイトル
- PR URL: PRのURL
- Head Branch: PRの作成元ブランチ
- Base Branch: PRの作成先ブランチ
- Related Issue: 関連Issue番号（確認できた場合のみ）
- Changes Summary: 変更内容の概要
- Testing: テスト結果
- Error: 失敗した場合のエラー内容

Statusの意味:

| Status | Meaning |
|---|---|
| SUCCESS | 新しいPRの作成と確認に成功した |
| EXISTING | 同じhead/baseの既存PRを発見したため、新規作成しなかった |
| FAILED | PRを作成できなかった、または必要な情報を確認できなかった |

## 7. Completion Report

### User Invocation

正常終了した場合、以下をユーザーに報告してください。

- PR番号
- PRタイトル
- PR URL
- headブランチ
- baseブランチ
- 変更内容の概要
- テスト結果

既存PRが見つかった場合は、そのPRの情報を報告してください。

失敗した場合は、失敗した処理とエラーメッセージを報告してください。

### Skill Invocation

呼び出し元のスキルに以下の結果を返してください。

- Status: SUCCESS / FAILED / EXISTING
- PR Number
- PR Title
- PR URL
- Head Branch
- Base Branch
- Related Issue (確認できた場合のみ)
- Changes Summary
- Testing
- Error (失敗した場合のみ)

呼び出し元のスキルは、結果に応じて以下のように処理してください。

- SUCCESS: PR作成後の後続処理に進む。
- EXISTING: 既存PRの情報を使用し、新規PR作成を行わない。
- FAILED: 後続処理を中断し、エラー内容を報告する。