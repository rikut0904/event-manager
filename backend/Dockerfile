FROM golang:1.24-alpine

WORKDIR /app

# 依存関係のインストール
COPY go.mod ./
RUN go mod download

# ソースコード全体をコピー（ここに service-account.json が含まれていればイメージに入ります）
COPY . .

# 開発用ツールの追加
RUN apk add --no-cache git

EXPOSE 8080

# コンテナ起動時に実行
CMD ["go", "run", "./cmd/api/main.go"]
