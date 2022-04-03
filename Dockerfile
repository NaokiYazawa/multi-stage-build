# Build stage
# golang:<version>-alpine は、Alpine Linux プロジェクトをベースにしている。
# イメージサイズを最小にするため、git、gcc、bash などは、Alpine-based のイメージには含まれていない。
FROM golang:1.16-alpine3.13 AS builder
# 作業ディレクトリの定義をする。今回は、app ディレクトリとした。
WORKDIR /app
# go.mod と go.sum を app ディレクトリにコピー
COPY go.mod go.sum ./
# 指定されたモジュールをダウンロードする。
RUN go mod download
# src ディレクトリの中身を app フォルダにコピーする
COPY . .
# 実行ファイルの作成
# -o はアウトプットの名前を指定。
# ビルドするファイル名を指定（今回は main.go）。
RUN go build -o main /app/main.go

# Run stage
# Goで作成したバイナリは Alpine Linux 上で動く。
# alpineLinux とは軽量でセキュアな Linux であり、とにかく軽量。
FROM alpine:3.13
# 作業ディレクトリの定義
WORKDIR /app
# Build stage からビルドされた main だけを Run stage にコピーする。
COPY --from=builder /app/main .
# ローカルの .env と .wait-for.sh をコンテナ側の app フォルダにコピーする
COPY .env .
COPY wait-for.sh .
# wait-for.sh の権限を変更
# x ･･･ 実行権限
RUN chmod +x wait-for.sh
# EXPOSE 命令は、実際にポートを公開するわけではない。
# これは、イメージを構築する人とコンテナを実行する人の間で、どのポートを公開するかについての一種の文書として機能する。
# 今回、docker-compose.yml において、api コンテナは 8080 ポートを解放するため「8080」とする。
EXPOSE 8080
# バイナリファイルの実行
CMD [ "/app/main" ]
