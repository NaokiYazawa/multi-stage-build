## 本記事の目標

本記事の目標は、**Docker の Multi-stage build を使って、Go 言語で開発したプロジェクトのイメージサイズを小さくする** です。

## 本記事の構成

本記事は全部で 3 章から構成されています。以下が各章の内容です。

**第１章**：Multi-stage build とは  
**第２章**：Go 言語を用いた簡単なプロジェクトの作成  
**第３章**：Multi-stage build の実行  
**第４章**：まとめ

## 第 1 章　 Multi-stage build とは

Multi-stage build とは何でしょうか。  
**Multi** には、「多くの」・「多重の」・「複数の」といった意味があります。  
よって、Multi-stage build とは「**複数のステージを用いたビルド**」となります。

では、**複数のステージを用いる** とはどういうことでしょうか。  
通常 Docker イメージには、イメージのビルドに関わるライブラリなども含まれています。しかし、**本番環境** ではアプリケーションの実行に必要なもののみをビルドしたいですよね。  
そこで複数のステージを用いると、この悩みが解消されるのです。

ここで、ステージを 2 つ用意するとします。  
1 つ目のステージでは、アプリケーションのビルドを行い Docker イメージを作成します。  
2 つ目のステージでは、1 つ目のステージで作成したイメージの中から必要なものだけをコピーしてきます。

このようにステージを 2 つ用意することで、最終的な Docker イメージには必要なものだけが含まれるようになるのです。  
その結果、イメージサイズが小さくなり、本番環境の運用のパフォーマンスが向上します。

## 第 2 章　 Golang × PostgreSQL の環境構築

本章では、Docker を用いた Golang × PostgreSQL の環境構築を行なっていきます。  
（プロジェクト全体のコードは、https://github.com/NaokiYazawa/multi-stage-build をご覧ください。）

### docker-compose.yml の作成

まずは、プロジェクトのルートディレクトリに、**docker-compose.yml** を作成してください。  
今回は、データベースである PostgreSQL と、API として機能する Go の 2 つのコンテナを定義して実行します。  
**services.api.build** において、**Dockerfile** のあるディレクトリのパスを指定しています。

```yaml:docker-compose.yml
version: "3.8"

services:
  postgres:
    container_name: postgres
    image: postgres:12.8
    restart: always
    env_file:
      - .env
    environment:
      - POSTGRES_USER=${POSTGRES_USER}
      - POSTGRES_PASSWORD=${POSTGRES_PASSWORD}
      - POSTGRES_DB=${POSTGRES_DB}
    ports:
      - 5432:5432
    volumes:
      - db:/var/lib/postgresql/data
  api:
    container_name: api
    build:
      # 「.」は本docker-compose.ymlがあるディレクトリ（現在のディレクトリ）を指す
      # 今回は、Dockerfile をルートディレクトリに配置する
      context: .
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    environment:
      POSTGRES_HOST: "${POSTGRES_HOST}"
    # depends_on は起動順を制御するだけである。
    # したがって、postgres コンテナが起動してから api コンテナが起動するという保証はされない
    depends_on:
      - postgres
    # entrypoint を設定すると、Dockerfile の ENTRYPOINT で設定されたデフォルトのエントリポイントが上書きされ、イメージのデフォルトコマンドがクリアされる。
    # つまり、Dockerfile に CMD 命令があれば、それは無視される。
    # よって、docker-compose.yml においても実行するコマンドを明示的に指定する必要がある。
    entrypoint: ["/app/wait-for.sh", "postgres:5432", "--"]
    command: ["/app/main"]

volumes:
  db:
```

```dockerfile:Dockerfile
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
# ルートディレクトリの中身を app フォルダにコピーする
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
```

```dockerfile:Dockerfile
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
# wait-for.sh の権限を変更
# x ･･･ 実行権限
RUN chmod +x wait-for.sh
# EXPOSE 命令は、実際にポートを公開するわけではない。
# これは、イメージを構築する人とコンテナを実行する人の間で、どのポートを公開するかについての一種の文書として機能する。
# 今回、docker-compose.yml において、api コンテナは 8080 ポートを解放するため「8080」とする。
EXPOSE 8080
# バイナリファイルの実行
CMD [ "/app/main" ]
```

【Multi-stage build を使わなかった場合】
![image.png](https://qiita-image-store.s3.ap-northeast-1.amazonaws.com/0/2279509/fbd33d80-d91c-e0dc-7acf-5b64edfeb1b4.png)

【Multi-stage build を使った場合】
![image.png](https://qiita-image-store.s3.ap-northeast-1.amazonaws.com/0/2279509/513cacb5-d9ec-edce-0f4d-7b3a4b06d751.png)

今回、Multi-stage build を使うと image size が 20 分の 1 程度になりました。

## 第 3 章　まとめ

https://stackoverflow.com/questions/49449012/dot-and-colon-meaning
