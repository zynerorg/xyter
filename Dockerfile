FROM golang:1.24.6 as builder


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

FROM scratch as bot

WORKDIR /app
COPY --from=builder --chown=100:100 /home/builder/app/build/bin/xyter-bot ./

CMD ["./xyter-bot"]

FROM scratch as api

WORKDIR /app
COPY --from=builder --chown=100:100 /home/builder/app/build/bin/xyter-api ./

CMD ["./xyter-api"]


FROM scratch as cli

WORKDIR /app
COPY --from=builder --chown=100:100 /home/builder/app/build/bin/xyter-cli ./

CMD ["./xyter-cli"]
