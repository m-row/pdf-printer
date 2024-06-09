# stage 1: wkhtmltopdf-builder
FROM surnet/alpine-wkhtmltopdf:3.20.0-0.12.6-small as wkhtmltopdf-builder

# stage 2: go-builder
FROM golang:1.22.3-alpine3.18 AS go-builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download && \ 
    apk --update --no-cache add ca-certificates git
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
    go build \
    -ldflags="-s -w -X main.CommitCount=$(git rev-list --count HEAD) -X main.CommitDescribe=$(git describe --always)" \
    -o main ./

# stage 3: main container
FROM alpine:3.20.0

WORKDIR /

RUN apk add --no-cache libstdc++ \
    libx11 \
    libxrender \
    libxext \
    libssl3 \
    ca-certificates \
    fontconfig \
    freetype \
    ttf-dejavu \
    ttf-droid \
    ttf-freefont \
    ttf-liberation \
    && apk add --no-cache --virtual .build-deps msttcorefonts-installer \
    && update-ms-fonts \
    && fc-cache -f \
    && rm -rf /tmp/* \
    && apk del .build-deps # buildkit

COPY --from=wkhtmltopdf-builder ["/bin/wkhtmltopdf", "/bin/wkhtmltopdf"]
COPY --from=go-builder ["/build/main", "/"]
COPY --from=go-builder ["/etc/ssl/certs/ca-certificates.crt", "/etc/ssl/certs/ca-certificates.crt"]
COPY ./.env ./.env

EXPOSE ${PORT}

CMD [ "/main" ]
