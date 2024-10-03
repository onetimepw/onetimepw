# builder
FROM golang:1.22-alpine AS builder
RUN apk add build-base
RUN apk add git --no-cache
ARG VERSION
ARG BUILD
ARG JOB_TOKEN
ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPATH=/

WORKDIR /build
COPY ["go.mod", "go.sum", "./"]
RUN go mod download
COPY ./ ./
RUN go build -trimpath \
    -ldflags "\
        -X 'github.com/onetimepw/onetimepw/build.Version=$VERSION' \
        -X 'github.com/onetimepw/onetimepw/build.Release=$BUILD' \
        -s -w" \
    -o /binary


#final
FROM alpine

RUN apk add --no-cache tzdata
ENV TZ=Europe/Moscow

WORKDIR /app
COPY --from=builder /binary /binary
COPY res/ /app/res/
RUN ls -la /app
EXPOSE 80

ENTRYPOINT [ "/binary" ]
