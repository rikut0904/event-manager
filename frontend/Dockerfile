FROM node:24

# ネイティブバイナリの実行に必要なライブラリを確認
RUN apt-get update && apt-get install -y \
    libc6 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# キャッシュを無効化して確実に再インストールするため package-lock.json をコピーしない
COPY package.json ./

# 実行環境(Linux arm64)向けのバイナリを強制的にインストール
RUN npm install

COPY . .

EXPOSE 3000

# 明示的に Turbopack を無効化
ENV NEXT_TELEMETRY_DISABLED 1
CMD ["npm", "run", "dev"]
