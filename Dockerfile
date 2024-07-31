FROM golang:1.22.5-alpine3.20 AS build

WORKDIR /tmp
RUN wget https://github.com/dhojayev/traderepublic-portfolio-downloader/archive/refs/heads/main.zip
RUN unzip main.zip

WORKDIR /tmp/traderepublic-portfolio-downloader-main
RUN go mod vendor
RUN go build -v -o /tmp/portfoliodownloader ./cmd/portfoliodownloader/public

FROM alpine
RUN apk -u add --no-cache tzdata
COPY --from=build /tmp/portfoliodownloader /opt/portfoliodownloader/portfoliodownloader

WORKDIR /opt/portfoliodownloader
ENTRYPOINT ["/opt/portfoliodownloader/portfoliodownloader", "-w", "--debug"]