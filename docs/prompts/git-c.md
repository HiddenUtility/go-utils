git-diff.log ファイルがあるか確認して，なければ終了してください．
git-diff.log フィルを読んで，変更内容を理解してください．
ターミナルに git commit 用のメッセージを**英語で**表示してください．
メッセージ形式は下記を参照してください．

## メッセージ形式
- **1行目（プレフィックス + 概要）:**
  - 形式: `<type>(<scope>): <簡潔な要約>`
  - 例: `docs(agents): AGENTS.mdのルールと運用方針を更新`
  - 主なtype: `feat` (新機能), `fix` (バグ修正), `docs` (ドキュメント), `refactor` (リファクタリング), `chore` (雑務・設定等)
- **2行目:** 空行
- **3行目以降（詳細）:**
  - 変更内容を箇条書きで具体的に記述する


```bash
echo "[git status --short]" > git-dff.log
git status --short >> git-dff.log
echo "[git diff --stat]" >> git-dff.log
git diff --stat >> git-dff.log
echo "[git diff]" >> git-dff.log
git diff >> git-dff.log

```
