# Common IDへのデータ移管

既存のFirebase UIDをCommon IDのユーザーIDへ置き換え、イベントの作成者参照も同時に更新するためのコマンドです。

## 事前準備

Common ID側でこのアプリの `client_id` とAPIキーを発行し、次の環境変数を設定します。

```sh
export DATABASE_URL='postgres://...'
export APP_ORIGIN='http://localhost:3000'
export COMMON_ID_API_ORIGIN='http://localhost:18080'
export COMMON_ID_CLIENT_ID='...'
export COMMON_ID_API_KEY='...'
```

event-managerをDockerで起動し、Common IDをホスト側で起動する場合は、バックエンドコンテナからホストへ接続するため、`COMMON_ID_API_ORIGIN=http://host.docker.internal:18080` を使用します。バックエンドをホスト上で直接起動する場合は `http://localhost:18080` のままです。

入力を省略すると、アプリDBのユーザーIDをFirebase Authで照合して移管対象を作ります。その場合はFirebase Admin用の `FIREBASE_PROJECT_ID` と、`FIREBASE_SERVICE_ACCOUNT_JSON` または `FIREBASE_SERVICE_ACCOUNT_KEY` も必要です。

## 実行

最初は必ずドライランを実行します。DBは変更されません。

```sh
make migrate-users
```

Firebase Authへの接続を避ける場合は、次のJSONLを用意して指定できます。

```json
{"source_user_id":"firebase-uid-1","email":"user@example.com","email_verified":true}
```

```sh
make migrate-users MIGRATE_FLAGS='-input ./migration-users.jsonl'
```

内容を確認し、Common ID側でのユーザー移管が完了した後に本実行します。

```sh
make migrate-users MIGRATE_FLAGS='-dry-run=false -update-db=true -output ./common-id-mapping.json'
```

このシステムではメール確認を行わないため、Common ID側で未検証メールの移管を許可する場合は、次のように指定します。既存Common IDユーザーを再利用し、新規ユーザー作成やパスワードリセットメール送信は行いません。

```sh
make migrate-users MIGRATE_FLAGS='-dry-run=false -update-db=true -allow-unverified'
```

`-allow-unverified` の既定値は `false` です。Common ID側でも対応済みの環境でのみ指定してください。

`-update-db=true` は `-dry-run=false` と同時でなければ実行できません。実行時は、対象DBのバックアップを取得し、アプリを停止または書き込みを止めた状態にしてください。更新は1トランザクションで行われ、失敗時は `users.id` と `events.creator_id` の変更がロールバックされます。

移管APIの失敗、Common IDの重複、既存ユーザーIDとの衝突、移管元ユーザーの欠落がある場合はDB更新を開始しません。APIキー、FirebaseサービスアカウントJSON、移管結果ファイルはGitへコミットしないでください。
