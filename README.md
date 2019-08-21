# Mesistant

**I help you to cook.**

## Features

- \*User account system
- \*Recipe management system
- \*Device managemant system
- \*Web frontend

## Dependencies

- Go 1.12
- MongoDB(@localhost:27017)

## Contribution

### ブランチ構成

- master

  動作する状態のブランチ。

- dev

  開発するブランチ。

### Server(Go)

- MongoDB(Docker)

  `docker run -d -p 27017:27017 mongo`

- Run server

  `go run main.go user.go recipe.go`

- (Test)

  `go test -v -run`



## FrontEnd

- 上記サーバー実行
- `/web`以下が`/`以下として実行されるので、それに合わせて開発。
  - 例　`/web/index.html`を置くと、ブラウザで`/index.html`にアクセスして見れる。

## サーバー実験用Curlコマンドたち
### レシピアップロード
`curl -X POST -H "Content-type: application/json" -d "{\"name\": \"Super delicious meal\", \"user_id\": \"ユーザーIDをいれてね\", \"description\": \"Sample description\", \"Steps\": [{\"type\": \"heat\", \"description\": \"加熱するぜい\", \"heat_strength\": 100}, {\"type\": \"add\", \"description\": \"なんか入れてくれ\", \"add_grams\": 100}]}" localhost:8080/recipe/upload
