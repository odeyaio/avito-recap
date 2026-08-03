FROM golang:1.26.5-alpine AS build

WORKDIR /src

COPY apps/backend/go.mod apps/backend/go.sum ./
RUN go mod download

COPY apps/backend/ ./
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/avito-recap ./cmd/avito-recap

FROM alpine:3.22

RUN apk add --no-cache ca-certificates \
    && addgroup -S app \
    && adduser -S -G app app

COPY --from=build /out/avito-recap /usr/local/bin/avito-recap

USER app
EXPOSE 8080

ENTRYPOINT ["avito-recap"]
