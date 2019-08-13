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

