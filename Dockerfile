FROM golang:1.24.6 AS builder


RUN useradd -m -u 10001 builder \
 && mkdir -p /home/builder/app \
 && chown -R builder:builder /home/builder
USER builder

WORKDIR /home/builder/app

RUN go install github.com/magefile/mage@v1.15.0

COPY --chown=builder:builder go.mod go.sum ./

RUN go mod download

COPY --chown=builder:builder ./ ./
#TEMP! When cli and api is added, this will change to BuildAll!
RUN go run mage.go BuildBot

FROM alpine AS runner

WORKDIR /app
COPY --from=builder --chown=100:100 /home/builder/app/build/bin/* ./
