# event-manager 要件定義・仕様書

## 1. 概要
本ドキュメントは、イベント管理プラットフォーム「event-manager」における具体的な機能要件、非機能要件、および利用者区分（ロール）の詳細を規定するものです。コミュニティを軸とした動的な権限管理と、プラットフォーム管理者の権限分離を明確化しています。

---

## 2. 利用者区分と権限管理 (RBAC)
ロールはユーザーに固定されるのではなく、原則として「コミュニティ」や「イベント」という枠組みに対して動的に付与されます。

### 2.1 権限一覧
| 利用者区分 | 権限内容 | 制限・補足 |
| :--- | :--- | :--- |
| **Admin** | ユーザーの強制退会、コミュニティの強制解体、システム全般の設定、ログ管理。 | 一般イベントの内容閲覧・編集・操作権限は一切持たない。 |
| **コミュニティ** | コミュニティの作成・編集、メンバーの招待、所属メンバーへのイベント管理権限の自動付与。 | イベント作成時にコミュニティを紐付けることで、所属ユーザー全員が管理可能となる。 |
| **イベントスタッフ** | 担当イベントの受付（QR/BLE）、参加者管理、Q&A対応、案内メール送信。 | コミュニティメンバーとして自動付与されるか、個別に追加される。 |
| **一般ユーザー** | イベント閲覧・申込、Q&A投稿、デジタル名刺交換、参加履歴（マイページ）閲覧。 | 全てのロールにおける基底権限。全ユーザーに付与される。 |

### 2.2 権限の継承・共有ロジック
- **コミュニティによる共有**: イベントにコミュニティIDを紐付けることで、そのコミュニティのメンバー全員が、個別に設定することなく当該イベントの管理権限（スタッフ相当）を取得する。
- **後付け共有**: スタッフが個人で作成したイベントであっても、後からコミュニティIDを入力することで、コミュニティ全体へ管理権限を移譲できる。

---

## 3. 機能要件

### 3.1 イベント管理・運営
- **connpass URL自動判別**: URL入力時に自動でイベント情報を取得し、名簿を同期。
- **マルチモード受付**: BLE（近接検知）、内蔵カメラQRスキャン、connpass ID手動入力。
- **PWA・オフライン対応**: 圏外でもアプリを起動し、BLE信号発信やQR表示による受付を可能にする。
- **AWS SESメール配信**: 申込通知、リマインド、フォローアップメールの一括送信。
- **グループ・座席表示**: 受付完了時に割り当てられた「座席番号」や「チーム名」を即時表示。

### 3.2 インタラクション・体験
- **コミュニティ誘導ポップアップ**: 運営側で設定したタイミング（申込時/受付時）にSlack/Discord等の案内を表示。
- **デジタル名刺交換**: Eight、Prally、GitHub等を統合し、会場内の人とワンタップ交換。
- **AI Q&A**: 匿名投稿された質問をAI（Gemini）がカテゴリ分け・要約。
- **マイページ（参加履歴）**: 過去の参加実績、配布資料、Q&Aログの振り返り。
- **持ち物チェックリスト**: イベント別の必須アイテム確認機能。

---

## 4. 非機能要件
- **Display IDの一意性**: URL `/[Display_ID]` の一意性を保証。connpass連携時はIDを固定。
- **プライバシー**: 位置情報の座標データ自体は保持せず、近接判定結果のみを記録。
- **リアルタイム性**: SSE（Server-Sent Events）による低遅延な情報反映。

---

## 5. 予約語一覧
ユーザーのDisplay IDとして使用できない、システム予約済みの文字列です。

| カテゴリ | 予約語 |
| :--- | :--- |
| システム基本 | admin, login, logout, register, signup, auth, api, settings, config, dashboard |
| 機能・リソース | events, communities, users, groups, notifications, messages, search, static, media |
| ファイル・メタデータ | favicon.ico, robots.txt, sitemap.xml, manifest.json, .well-known |

---

## 6. APIエンドポイント設計

### 6.1 認証・ユーザー管理
| Method | Endpoint | Description | 備考 |
| :--- | :--- | :--- | :--- |
| POST | `/api/v1/auth/verify` | トークン検証・サインアップ | 初回ログイン時のユーザー作成 |
| GET | `/api/v1/users/me` | ログインユーザー情報取得 | 自分の基本プロフィール |
| PATCH | `/api/v1/users/me` | プロフィール/Display ID更新 | ID重複・予約語チェック実施 |
| GET | `/api/v1/users/me/history` | 自分のイベント参加履歴取得 | 過去のイベント資料へのリンク含む |
| GET | `/[display_id]` | 公開プロフィール取得 | 【公開】 名刺交換用ページ |

### 6.2 コミュニティ管理
| Method | Endpoint | Description | 備考 |
| :--- | :--- | :--- | :--- |
| POST | `/api/v1/communities` | 新規コミュニティ作成 | 誰でも作成可能 |
| GET | `/api/v1/communities` | 所属コミュニティ一覧取得 | 自分がメンバーのグループ |
| GET | `/api/v1/communities/:id` | コミュニティ詳細取得 | 概要、オーナー情報等 |
| PATCH | `/api/v1/communities/:id` | コミュニティ情報更新 | オーナー・Adminのみ |
| DELETE | `/api/v1/communities/:id` | コミュニティ削除（解散） | オーナーのみ |
| GET | `/api/v1/communities/:id/members` | メンバー一覧取得 | 所属ユーザーのリスト |
| POST | `/api/v1/communities/:id/invite` | メンバー招待 | 招待制の追加 |
| DELETE | `/api/v1/communities/:id/members/:u_id` | メンバー除名 | オーナー・Adminのみ |

### 6.3 イベント管理
| Method | Endpoint | Description | 備考 |
| :--- | :--- | :--- | :--- |
| POST | `/api/v1/events` | イベント作成 | URL入力でconnpass自動判別 |
| GET | `/api/v1/events` | イベント検索・一覧取得 | 【公開】 公開中のイベント |
| GET | `/api/v1/events/:id` | イベント詳細取得 | 概要、持ち物リスト、座席情報 |
| PATCH | `/api/v1/events/:id` | イベント編集 | 共同管理コミュニティの紐付け等 |
| DELETE | `/api/v1/events/:id` | イベント削除 | 作成者またはコミュニティAdmin |
| POST | `/api/v1/events/:id/register` | イベント参加申し込み | 申し込み完了通知（SES）送信 |
| GET | `/api/v1/events/:id/attendees` | 参加者名簿・受付状況取得 | スタッフ・コミュニティメンバーのみ |

### 6.4 現場オペレーション（受付・双方向機能）
| Method | Endpoint | Description | 備考 |
| :--- | :--- | :--- | :--- |
| POST | `/api/v1/events/:id/checkin` | チェックイン実行 | BLE/QR/connpass ID照合 |
| POST | `/api/v1/events/:id/qa` | 質問投稿 | 匿名投稿の許可設定あり |
| GET | `/api/v1/events/:id/qa` | 質問一覧取得 | SSEによるリアルタイム同期を推奨 |
| POST | `/api/v1/events/:id/survey` | アンケート回答送信 | ライブ投票機能 |
| GET | `/api/v1/events/:id/summary` | AI Q&A要約の取得 | Geminiによる要約結果 |
| POST | `/api/v1/events/:id/summary` | AI Q&A要約の実行・更新 | スタッフのみ |

### 6.5 システムAdmin（プラットフォーム管理）
| Method | Endpoint | Description | 備考 |
| :--- | :--- | :--- | :--- |
| DELETE | `/api/v1/admin/users/:id` | ユーザーの強制退会 | システムAdminのみ |
| DELETE | `/api/v1/admin/communities/:id` | コミュニティの強制解体 | システムAdminのみ |
| GET | `/api/v1/admin/logs` | システム運営ログの取得 | アクセス・エラー・セキュリティログ |

---

## 7. データベース設計 (ER図)

### 7.1 ユーザー・プロフィール関連
- **users (ユーザー基本情報)**
  - `id (UUID, PK)`: システム内部の不変ID
  - `display_id (VARCHAR(50), UQ)`: URL用ID
  - `email (VARCHAR(255), UQ)`: 連絡先
  - `connpass_id (VARCHAR(100), Nullable)`: connpassアカウント連携ID
  - `role (TEXT)`: システム管理権限 (admin)
  - `created_at (TIMESTAMP)`, `updated_at (TIMESTAMP)`
- **profile_links**
  - `id (UUID, PK)`, `user_id (UUID, FK)`, `platform_name (VARCHAR(50))`, `url (TEXT)`

### 7.2 コミュニティ・権限関連
- **communities (コミュニティ)**
  - `id (UUID, PK)`, `name (VARCHAR(100))`, `description (TEXT)`, `owner_id (UUID, FK)`
- **community_members (コミュニティ所属)**
  - `community_id (UUID, PK, FK)`, `user_id (UUID, PK, FK)`, `role (VARCHAR(20))`

### 7.3 イベント・運営関連
- **events (イベント基本情報)**
  - `id (UUID, PK)`, `creator_id (UUID, FK)`, `community_id (UUID, FK)`, `title (VARCHAR(200))`, `source_url (TEXT)`, `popup_timing (VARCHAR(20))`, `start_time (TIMESTAMP)`
- **event_belongings (持ち物マスター)**
  - `id (UUID, PK)`, `event_id (UUID, FK)`, `item_name (VARCHAR(100))`, `is_required (BOOLEAN)`
- **registrations (申し込み・座席管理)**
  - `id (UUID, PK)`, `event_id (FK)`, `user_id (FK)`, `status (VARCHAR(20))`, `assigned_group (VARCHAR(50))`

### 7.4 ログ・相互作用関連
- **attendance_logs (受付記録)**
  - `id (UUID, PK)`, `registration_id (UUID, FK)`, `timestamp (TIMESTAMP)`, `method (VARCHAR(20))`
- **interactions (Q&A・アンケート)**
  - `id (UUID, PK)`, `event_id (UUID, FK)`, `user_id (UUID, FK, Nullable)`, `type (VARCHAR(20))`, `content (TEXT)`, `ai_summary (TEXT)`
