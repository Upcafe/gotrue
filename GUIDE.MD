# Remove
`````` from `./migrations`

# Start mysql
docker-compose up --build

# set connection string in .env
```sh
root:my-secret-pw@tcp(127.0.0.1:3306)/gotrue_development?parseTime=true&multiStatements=true
```

# PHPMyAdmin 8080
Server:
user: root
password: my-secret-pw

# migrate
migrate -path ./migrations -database "mysql://root:my-secret-pw@tcp(127.0.0.1:3306)/gotrue_development?parseTime=true&multiStatements=true" up

# mailer 
.env GOTRUE_MAILER_AUTOCONFIRM=true

test