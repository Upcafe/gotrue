FROM golang:1.13-alpine as build
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux

RUN apk add --no-cache make git

WORKDIR /go/src/github.com/netlify/gotrue

# Pulling dependencies
COPY ./Makefile ./go.* ./
RUN make deps

# Building stuff
COPY . /go/src/github.com/netlify/gotrue
RUN make build

FROM alpine:3.7
RUN adduser -D -u 1000 netlify

RUN apk add --no-cache ca-certificates
COPY --from=build /go/src/github.com/netlify/gotrue/gotrue /usr/local/bin/gotrue
COPY --from=build /go/src/github.com/netlify/gotrue/migrations /usr/local/etc/gotrue/migrations/

ENV GOTRUE_DB_MIGRATIONS_PATH /usr/local/etc/gotrue/migrations

ENV GOTRUE_JWT_SECRET "CHANGE-THIS! VERY IMPORTANT!"
ENV GOTRUE_JWT_EXP 3600
ENV GOTRUE_JWT_AUD api.netlify.com
ENV GOTRUE_DB_DRIVER mysql
ENV DATABASE_URL "root@tcp(127.0.0.1:3306)/gotrue_development?parseTime=true&multiStatements=true"
ENV GOTRUE_API_HOST localhost
ENV PORT 9999
ENV GOTRUE_SITE_URL https://gotrueintegration.netlify.app
ENV GOTRUE_SMTP_HOST smtp.mandrillapp.com
ENV GOTRUE_SMTP_PORT 587
ENV GOTRUE_SMTP_USER test@example.com
ENV GOTRUE_SMTP_PASS super-secret-password
ENV GOTRUE_SMTP_ADMIN_EMAIL admin@example.com
ENV GOTRUE_MAILER_SUBJECTS_CONFIRMATION "Welcome to GoTrue!"
ENV GOTRUE_MAILER_SUBJECTS_RECOVERY "Reset your GoTrue password!"
ENV GOTRUE_EXTERNAL_GOOGLE_REDIRECT_URI http://localhost:9999/callback
ENV GOTRUE_LOG_LEVEL DEBUG
ENV GOTRUE_OPERATOR_TOKEN super-secret-operator-token
ENV GOTRUE_WEBHOOK_URL http://register-lambda:3000/
ENV GOTRUE_WEBHOOK_SECRET test_secret
ENV GOTRUE_WEBHOOK_RETRIES 5
ENV GOTRUE_WEBHOOK_TIMEOUT_SEC 3
ENV GOTRUE_WEBHOOK_EVENTS validate,signup,login


USER netlify
CMD ["gotrue"]
