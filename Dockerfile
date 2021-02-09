FROM golang:1.15-alpine3.13 AS build
RUN apk add mariadb-dev postgresql-dev sqlite-dev git
RUN apk add gcc musl-dev

WORKDIR /src/optigit
COPY . .
RUN go build .

FROM alpine:3.13
COPY --from=build /src/optigit/optigit /usr/bin/optigit

CMD ["optigit"]
