# echo_lambda_sample

### start docker

```
docker compose -f 'docker-compose.yml' up -d --build
```

### start lambda server

```
# ローカルにGoインストール済み。API用とMySQL用のDockerコンテナ起動済み。
cd src
go mod tidy
go build

# how to install sam https://docs.aws.amazon.com/ja_jp/serverless-application-model/latest/developerguide/install-sam-cli.html
sam build --use-container

sam local start-api --docker-network api_default

# 高速化オプションあり
sam local start-api --docker-network api_default --warm-containers EAGER

# access http://localhost:3000/{function_name}/health
```
