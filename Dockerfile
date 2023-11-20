# ステージ1: ビルド環境
FROM golang:1.21-alpine AS builder

# 作業ディレクトリを設定
WORKDIR /app

# ソースコードをコピー
COPY . .

# Go アプリケーションをビルド
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# ステージ2: 実行環境
FROM gcr.io/distroless/base-debian12

# 作業ディレクトリを設定
WORKDIR /

# ビルドしたバイナリをコピー
COPY --from=builder /app/server .

# デフォルトコマンドを設定（環境変数PORTを使用）
CMD ["./server"]
