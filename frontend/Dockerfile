FROM node:24

# ネイティブバイナリの実行に必要なライブラリを確認
RUN apt-get update && apt-get install -y \
    libc6 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# package-lock.jsonを使って依存関係を固定し、Linuxコンテナ用のoptional native依存もインストールする
COPY package.json package-lock.json ./

RUN npm ci --include=optional

COPY . .

EXPOSE 3000

# 明示的に Turbopack を無効化
ENV NEXT_TELEMETRY_DISABLED 1
CMD ["npm", "run", "dev"]
