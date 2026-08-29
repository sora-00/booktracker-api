# entity ここが全て

- Me
    - ID = firebase.ID
    - createdAt
    - updatedAt
    - **役割**
        - 初回ログイン時に UID で作成
        - 以後、Book / Log と紐づけるためのユーザーIDとして利用
- Book
    - User.ID（= Me.ID / Firebase UID）
    - `ID`
    - `Title`: string（更新可, 必須, 30字以内）
    - `Author`: string（更新可, 必須, 30字以内）
    - `Publish`: string（更新可, 必須, 20字以内）
    - `CoverURL`: string（更新可, 必須）
    - `PagesAll`: number
        - Log が 0 件のときのみユーザー更新可
        - `1 < PagesAll`
        - 4桁以内
    - `TargetReadDate`: timestamp
        - `PagesAll` と `TargetReadPagesPerDay` から計算可能
        - `now < TargetReadDate`
        - create 時の初期値: `now + 1ヶ月`
        - 更新可
    - `TargetReadPagesPerDay`: number
        - 自動計算:
            - `TargetReadPagesPerDay = PagesAll / (TargetReadDate - now)`
        - 逆にユーザー入力した場合:
            - `TargetReadDate = now + PagesAll / TargetReadPagesPerDay`
        - `0 < TargetReadPagesPerDay < PagesAll`
        - 4桁以内
        - 更新可
    - `Status`: `"unread" | "reading" | "read"`（更新可）
    - `createdAt`
    - `updatedAt`
    - `ReadPages`: number
        - `status != unread` のとき必須
        - `ReadPages <= PagesAll`
        - 4桁以内
        - 基本は Log 追加時に自動更新。それ以外の更新は不可。
    - `Background`: string（任意, 200字以内, 更新可）
- Log
    - `User.ID`
    - `Book.ID`
    - `ID`
    - `ReadDate`: timestamp
        - 同じ Book の「最新の ReadDate」より大きいことが原則
        - その条件を満たすなら更新可
    - `StartPage`: number
        - `0 < StartPage`
        - ひとつ前の Log の `EndPage <= StartPage`
        - `StartPage < EndPage`
        - 4桁以内
    - `EndPage`: number
        - `1 < EndPage`
        - `EndPage <= Book.PagesAll`
        - `StartPage < EndPage`
        - 4桁以内
    - `PagesRead = EndPage - StartPage + 1`
        - 4桁以内
    - `createdAt`
    - `updatedAt`
    - `note`: string（任意, 800字以内）

# ripository もはやクエリを書くだけのところ

### **Me リポジトリ**

- `Create`
    - アカウント作成する時。UID で作成。
- `Get`
    - ログインする時
- `Delete`
    - 退会する時

### **Book リポジトリ**

- `Create`
    - Book 新規作成（POST）
- `Update`
    - `me.ID == Book.UserID` で照合して PUT
- `GetAll`
    - 本の一覧
    - `me.ID == Book.UserID` でフィルタ
    - `createdAt` 降順
- `GetByID`
    - `me.ID == Book.UserID && Book.ID == :id`
    - `createdAt` 降順（単一だけどクエリ的にはそういう設計）
- `GetByStatus`
    - `me.ID == Book.UserID && Book.Status == :status`
    - `createdAt` 降順
- `DeleteByBookID`
    - `me.ID == Book.UserID && Book.ID == :id`
    - Book削除
    - Log 全削除は service で実施

### **Log リポジトリ**

- `Create`
- `Update`
- `GetAll`
    - `me.ID == Log.UserID`
    - `readAt` 昇順
- `GetByID`
    - `me.ID == Log.UserID && Log.ID == :id`
    - `readAt` 昇順
- `GetByBookID`
    - `me.ID == Log.UserID && Log.BookID == :bookId`
    - `readAt` 昇順
- `DeleteByBookID`
    - `me.ID == Log.UserID && Log.BookID == :bookId`
- `DeleteByLogID`
    - `me.ID == Log.UserID && Log.ID == :id`

# service repositoryを使ってこねこね

- `CreateBook`
    - 上記の Book バリデーション一式を実施
    - `TargetReadDate` / `TargetReadPagesPerDay` の自動計算
    - `createdAt` / `updatedAt` の自動セット
- `UpdateBook`
    - 更新可能フィールドのみ：
        - `Title`, `Author`, `Publish`, `CoverURL`, `PagesAll`,`TargetReadDate`, `TargetReadPagesPerDay`, `Status`, `Background`
    - `PagesAll` は Log が0件の場合のみ変更可
    - `ReadPages` は「log追加による自動更新以外不可」
- `DeleteBook`
    - `bookrepo.Delete` ＋ `logrepo.DeleteByBookID`
- `FindLogByBookStatus`
    - `bookrepo.GetByStatus` 結果に対し `logrepo.GetAll` を紐づけ
- `CreateLog` / `UpdateLog`
    - Log の各種バリデーション（ページ範囲、ReadDate の単調増加など）
    - Book の `ReadPages` を自動更新（合計 / 最終ページ など、ここは設計が必要）

# エンドポイント

### **Book**

- **GET `/api/books`**
    - 認証ユーザーの Book 全件
- **GET `/api/books/status/{status}`**
    - `status` = `unread | reading | read`
- **GET `/api/books/{id}`**
    - Book ID（数値）。バリデーションはリクエスト層
- **GET `/api/books/status/{status}/logs`**
    - ステータス別の本＋ログ
- **GET `/api/books/{bookId}/logs`**
    - 特定 Book のログ一覧
- **POST `/api/books`**
    - Body: CreateBook 用 JSON
- **PUT `/api/books/{id}`**
    - Body: UpdateBook 用 JSON
- **DELETE `/api/books/{id}`**
    - 削除時に関連 Log も削除（service 内で `logrepo.DeleteByBookID`）

### **Log**

- **GET `/api/logs`**
    - 認証ユーザーの Log 一覧（必要ならクエリで絞り込み）
- **GET `/api/logs/{id}`**
- **POST `/api/logs`**
    - Body: CreateLog JSON
- **PUT `/api/logs/{id}`**
    - Body: UpdateLog JSON
- **DELETE `/api/logs/{id}`**
- **DELETE `/api/books/{bookId}/logs`**
    - ある本に紐づくログをすべて削除

### **Me / 認証**

- **middleware**
    - Firebase ID トークン検証
    - `context` に `userID`（UID）を埋め込む
- **GET `/api/me`**
    - ログイン済みユーザーの Me 情報取得
- 初回ログイン時:
    - 認証済みUIDが Me に存在しなければ `Create` する

# バリデーション

- `app/domain/service/validation` あるいは `app/domain/service/book_validation.go` / `log_validation.go` のような形で
    - `ValidateBookForCreate`
    - `ValidateBookForUpdate`
    - `ValidateLogForCreate`
    - `ValidateLogForUpdate`
- 既存の `request/book.go` に書かれているバリデーションロジックをここに寄せつつ、仕様書に沿ったチェック（文字数、数値範囲、相関チェック）を集約していくイメージです。

# ルール

- 決まった値しか入らないものは、entityとかrespとかで型定義をちゃんとして、それを参照するようにする
- 更新 API は基本、更新後リソースを返す
- restful APIを意識する

ディレクトリ構成は以下の通りで、上層のフォルダ名とかを下層のファイル名や関数名、変数名に使いたくない。パスで意味が被るものを入れたくない

```jsx
.
├── main.go
├── domain/
│   ├── entity/
│   │   ├── user.go
│   │   ├── book.go
│   │   └── log.go
│   ├── repository/
│   │   ├── user.go
│   │   ├── book.go
│   │   └── log.go
│   └── service/
│       ├── user.go
│       ├── book.go
│       └── log.go
├── usecase/
│   ├── user.go
│   ├── book.go
│   ├── log.go
│   ├── request/
│   │   ├── book.go
│   │   └── log.go
│   └── response/
│       ├── book.go
│       └── log.go
└── controller/
    ├── user.go
    ├── book.go
    └── log.go
```

- FireStoreを使う
- FirestoreはコレクションをEntityに対応させる
- middlewareでFirebase Authのトークン検証を一括で行う
- エラーは独自エラー型を定義してHTTPステータスと対応させる
- レスポンスは正常系/エラー系で形式を統一する
    - 正常系: そのまま JSON 本体を返す
    - エラー系: { "error": "message" }
- Linuxのスモールイズビューティフルの精神で、一つの関数の関心ごとはなるべく少なくし、それを組み合わせて使っていく書き方にする。
- 簡潔さよりも、読みやすさ、拡張性の高さを優先して実装する。
- serviceはなるべく膨らませない
- usecaseが仲介的な立ち位置だから、そこは膨らんでもいい
- usecaseでreqを受けて、repoとかservice使って、それをresに渡してそれ返すって感じ


## ログイン機能

**フロント（Expo）**

1. `@react-native-google-signin/google-signin` でGoogleログイン
2. GoogleのIDトークンをFirebaseに渡してFirebaseのIDトークンを取得
3. GoのAPIにFirebaseのIDトークンをヘッダーに乗せてリクエスト

**Go**

1. `firebase.google.com/go` のAdmin SDKでIDトークンを検証
2. 初回ログイン時はMeテーブルにUIDでユーザー作成