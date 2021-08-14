# Guide
This repo is a simple REST-API for token based authentication in Go forked from `github.com/netlify/gotrue`.

## How to recreate changes in this fork?

### Install Tooling
- vscode
- docker
- golang 
- snoop (windows only)
- brew (mac only)
- golang-migrate

### Create Docker Compose

#### Create a docker compose file.
This will serve as the development database using PHPAdmin as a simple UI.

```docker
services:

  db:
    image: mysql:8.0.17
    environment:
      MYSQL_ROOT_PASSWORD: my-secret-pw
      MYSQL_DATABASE: gotrue_development
    ports:
      - "3306:3306"
      
  phpmyadmin:
    image: phpmyadmin
    restart: always
    ports:
      - 8080:80
    environment:
      PMA_HOST: db
      PMA_PORT: 3306
      PMA_ARBITRARY: 1
    links:
      - db
    volumes:
      - my-datavolume:/var/lib/mysql

volumes:
  my-datavolume:
```

#### Run
```sh
docker compose up
```

### Migrations

#### Fix migrations
Replace `{{ index .Options "Namespace" }}`to replace this snippet with nothing.

#### Run with golang-migrate
```sh
migrate -path ./migrations -database "mysql://root:my-secret-pw@tcp(127.0.0.1:3306)/gotrue_development?parseTime=true&multiStatements=true" up
```

#### Check migration
Navigate to PHPMyAdmin on `http://localhost:8080` in the browser and login with:
```
Server:
user:       root
password:   my-secret-pw
```

If all is well you should see that the tables are created.

## Set env variables
Copy and paste `.env.example` and rename to `.env`.

#### Change the following lines
```env
GOTRUE_JWT_AUD=localhost
DATABASE_URL="root:my-secret-pw@tcp(127.0.0.1:3306)/gotrue_development?parseTime=true&multiStatements=true"
```

### Add
Add this line to skip email validation.
```env
GOTRUE_MAILER_AUTOCONFIRM=true
```

### Remove
Webhook not present so delete this line.
```sh
GOTRUE_WEBHOOK_URL=http://register-lambda:3000/
```

## Run API 
```sh
go run .
```