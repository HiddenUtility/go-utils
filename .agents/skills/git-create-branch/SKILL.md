---
name: git-create-branch
description: >
  mainまたはdevブランチから作業用ブランチを作成し、リモートリポジトリにpushする。
  ユーザーからの直接呼び出しと、他のスキルからの呼び出しの両方に対応する。
  作業概要と任意のIssue番号をもとに、作業内容に適したプレフィックスと
  descriptionからブランチ名を生成する。
---

# Git Create Branch Skill

## 1. Purpose

mainまたはdevブランチから、新しい作業用ブランチを作成してください。

ブランチ名は作業内容に応じたプレフィックスを使用してください。

基本形式:

- Issue番号あり: `<prefix>/<issue-number>-<description>`
- Issue番号なし: `<prefix>/<description>`

プレフィックスは作業内容に応じて選択してください。

ブランチ作成後は、リモートリポジトリにもpushしてください。

このスキルは、以下の2つの呼び出し方法に対応してください。

1. ユーザーからの直接呼び出し
2. 他のスキルからの呼び出し

呼び出し元に応じて、必要な情報の取得方法と結果の返却方法を切り替えてください。

## 2. Execution Rules

- Gitコマンドはターミナルから実行してください。
- Git操作の結果を必ず確認してください。
- ユーザーの既存変更を破棄してはいけません。
- 既存ブランチを削除・上書きしてはいけません。
- ブランチ作成後はリモートにpushしてください。
- Git操作に失敗した場合は処理を停止してください。
- Issue番号は任意とし、未指定でも処理を続行してください。
- 呼び出し元から受け取った情報を優先し、同じ内容をユーザーに再質問しないでください。
- 呼び出し元が指定した作業概要を勝手に変更しないでください。
- ブランチのプレフィックスを固定しないでください。
- 作業内容に応じて適切なプレフィックスを選択してください。
- 呼び出し元のスキルからブランチ命名規則が指定された場合は、その指定を優先してください。
- プレフィックスの判断ができない場合は、ユーザーまたは呼び出し元に確認してください。
- Git操作に失敗した場合、独自の代替手段で処理を続行してはいけません。

## 3. Determine Invocation Type

最初に、呼び出し方法を判定してください。

### 3.1 User Invocation

ユーザーが直接このスキルを呼び出した場合は、以下の手順で情報を取得してください。

1. 作業概要を質問する。
2. 関連するGitHub Issue番号があれば教えてもらう。
3. Issue番号は任意とし、未指定でも処理を続行する。

作業概要が不明確な場合のみ、追加で質問してください。

### 3.2 Skill Invocation

他のスキルから呼び出された場合は、呼び出し元から渡された情報を使用してください。

受け取る情報:

- Task Description (必須)
  - 実施する作業の概要
- Issue Number (任意)
  - 関連するGitHub Issue番号
- Branch Prefix (任意)
  - 呼び出し元が指定するブランチプレフィックス
- Branch Naming Rule (任意)
  - 呼び出し元が指定するブランチ命名規則

以下のルールに従ってください。

1. Task Descriptionが提供されている場合は、その内容を使用する。
2. Issue Numberが提供されている場合は、その番号を使用する。
3. Issue Numberが提供されていない場合は、Issue番号なしの形式を使用する。
4. Branch Prefixが提供されている場合は、その指定を優先する。
5. Branch Naming Ruleが提供されている場合は、その指定を優先する。
6. Branch Prefixが提供されていない場合は、作業内容から適切なプレフィックスを選択する。
7. Task Descriptionが不足している場合のみ、呼び出し元に確認する。
8. 呼び出し元から受け取った情報をユーザーに再質問しない。
9. 呼び出し元に対して、ブランチ作成結果を返す。

呼び出し元が不明な場合は、ユーザーからの直接呼び出しとして扱ってください。

## 4. Workflow

### Step 1: Check Current Branch

現在のブランチを確認してください。

```bash
git branch --show-current
```

以下のブランチにいる場合のみ処理を続行してください。

- main
- dev

それ以外のブランチにいる場合は、現在のブランチ名を表示して終了してください。

### Step 2: Obtain Requirements

Section 3で判定した呼び出し方法に従って作業概要とIssue番号を取得してください。

Task Descriptionは必須です。

Issue Numberは任意です。

Branch PrefixとBranch Naming Ruleも、呼び出し元から指定されている場合のみ使用してください。

### Step 3: Generate Branch Name

作業概要を分析し、作業内容に最も適したプレフィックスを選択してください。

#### Prefix Rules

| Prefix | 使用条件 |
|---|---|
| features | 新機能・機能追加 |
| fix | バグ修正・不具合修正 |
| refactor | リファクタリング・内部構造の改善 |
| docs | ドキュメント・説明書の変更 |
| test | テストの追加・修正 |
| chore | 保守作業・依存関係の更新・その他の雑務 |
| hotfix | 緊急性の高い本番障害の修正 |

上記に該当しない場合は、作業内容に適したプレフィックスを選択してください。

複数の作業内容が含まれる場合は、主目的に最も近いプレフィックスを選択してください。

#### Branch Name Format

Issue Numberがある場合:

```text
<prefix>/<issue-number>-<description>
```

Issue Numberがない場合:

```text
<prefix>/<description>
```

#### Naming Rules

- 英語の小文字を使用する。
- 単語の区切りにはハイフン `-` を使用する。
- 空白や特殊文字を使用しない。
- descriptionは作業内容を簡潔に表現する。
- Issue Numberがある場合はdescriptionの前に配置する。
- 作業内容に適したプレフィックスを使用する。
- 既存ブランチと重複しないことを確認する。

#### Examples

```text
features/add-user-authentication
features/123-add-user-authentication

fix/login-validation
fix/456-login-validation

refactor/logger-implementation

docs/update-readme

test/add-unit-tests

chore/update-dependencies

hotfix/789-fix-payment-error
```

生成したブランチ名を確定してください。

ユーザーからの直接呼び出しの場合は、生成したブランチ名をユーザーに表示してください。

スキルからの呼び出しの場合は、呼び出し元への処理結果に含めてください。

### Step 4: Check Working Tree

現在の作業状態を確認してください。

```bash
git status --short
```

未コミットの変更がある場合は、安全のため処理を停止してください。

変更内容を破棄したり、勝手にstashしたりしてはいけません。

ユーザーにコミットまたはstashなどの操作を依頼してください。

### Step 5: Update Base Branch

現在のブランチを最新状態に更新してください。

```bash
git fetch origin
git pull --ff-only
```

失敗した場合は処理を停止してください。

### Step 6: Create Branch

生成したブランチ名で新しいブランチを作成してください。

```bash
git switch -c <branch-name>
```

作成後、現在のブランチ名を確認してください。

```bash
git branch --show-current
```

作成したブランチに切り替わっていることを確認してください。

### Step 7: Push Remote Branch

作成したブランチをリモートにpushしてください。

```bash
git push -u origin <branch-name>
```

失敗した場合は処理を停止してください。

### Step 8: Verify

以下を確認してください。

```bash
git status --short
git branch -vv
```

以下の条件を確認してください。

- 現在のブランチが作成したブランチであること。
- リモートブランチが存在すること。
- upstreamが正しく設定されていること。

すべて確認できた場合のみ正常終了としてください。

## 5. Error Handling

Git操作に失敗した場合:

1. 失敗したコマンドを表示する。
2. エラーメッセージを表示する。
3. 原因が分かる場合は説明する。
4. ユーザーに必要な操作を依頼する。
5. それ以上のGit操作を行わず終了する。
6. 呼び出し元がスキルの場合は、失敗した処理とエラー内容を呼び出し元にも返す。

禁止事項:

- `git reset --hard`
- `git checkout --`
- `git restore`
- 強制push
- ユーザーの変更を破棄する操作
- エラーを無視した処理の続行
- 既存ブランチの削除
- 既存ブランチへの勝手な切り替え

## 6. Result

処理結果には以下の情報を含めてください。

- Status: SUCCESS / FAILED
- Branch Name: 作成したブランチ名
- Branch Prefix: 使用したプレフィックス
- Base Branch: 作成元のブランチ名
- Issue Number: 指定されている場合のみ
- Remote Push: 成功 / 失敗
- Current Branch: 処理終了時のブランチ名
- Error: 失敗した場合のエラー内容

## 7. Completion Report

### User Invocation

正常終了した場合、以下をユーザーに報告してください。

- 作成したブランチ名
- 使用したプレフィックス
- 作成元のブランチ名
- Issue Number (指定されている場合のみ)
- リモートpushの結果
- 現在のブランチ名

失敗した場合は、失敗した処理とエラーメッセージを報告してください。

### Skill Invocation

呼び出し元のスキルに以下の結果を返してください。

- Status: SUCCESS / FAILED
- Branch Name
- Branch Prefix
- Base Branch
- Issue Number (指定されている場合のみ)
- Remote Push
- Current Branch
- Error (失敗した場合のみ)

成功した場合のみ、呼び出し元のスキルは後続処理に進んでください。

失敗した場合は、呼び出し元のスキルも後続処理を中断してください。