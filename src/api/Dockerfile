FROM golang:1.14 AS base

FROM base AS dev
RUN go get -u \
    github.com/mdempsky/gocode \
    github.com/uudashr/gopkgs/v2/cmd/gopkgs \
    github.com/ramya-rao-a/go-outline \
    github.com/acroca/go-symbols \
    golang.org/x/tools/cmd/guru \
    golang.org/x/tools/gopls \
    golang.org/x/tools/cmd/gorename \
    github.com/cweill/gotests/... \
    github.com/fatih/gomodifytags \
    github.com/josharian/impl \
    github.com/davidrjenni/reftools/cmd/fillstruct \
    github.com/haya14busa/goplay/cmd/goplay \
    github.com/godoctor/godoctor \
    github.com/go-delve/delve/cmd/dlv \
    github.com/stamblerre/gocode \
    github.com/rogpeppe/godef \
    github.com/sqs/goreturns \
    golang.org/x/lint/golint \
    golang.org/x/tools/gopls
RUN go get -u github.com/cosmtrek/air
CMD air

FROM base AS builder
ENV GIN_MODE=release
ENV PORT=80
COPY . /go/src/github.com/philippecarle/mood/api
WORKDIR /go/src/github.com/philippecarle/mood/api
RUN go build -tags=jsoniter -o /go/bin/mood

FROM scratch
COPY --from=builder /go/bin/mood /go/bin/mood