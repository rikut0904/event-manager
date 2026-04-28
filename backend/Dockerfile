FROM golang:1.26-alpine

WORKDIR /app

# 依存関係のインストール
COPY go.mod ./
RUN go mod download

# ソースコード全体をコピー
COPY . .

# 開発用ツールの追加
RUN apk add --no-cache git

EXPOSE 8080

# コンテナ起動時に実行
CMD ["go", "run", "./cmd/api/main.go"]
