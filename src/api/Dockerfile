FROM golang:1.15 AS base
ARG BUILT_CMD
ENV BUILT_CMD ${BUILT_CMD}

FROM base AS dev
RUN go get \
    github.com/acroca/go-symbols \
    github.com/cweill/gotests/... \
    github.com/davidrjenni/reftools/cmd/fillstruct \
    github.com/fatih/gomodifytags \
    github.com/go-delve/delve/cmd/dlv \
    github.com/godoctor/godoctor \
    github.com/haya14busa/goplay/cmd/goplay \
    github.com/josharian/impl \
    github.com/mdempsky/gocode \
    github.com/mgechev/revive \
    github.com/ramya-rao-a/go-outline \
    github.com/rogpeppe/godef \
    github.com/sqs/goreturns \
    github.com/uudashr/gopkgs/v2/cmd/gopkgs \
    github.com/zmb3/gogetdoc \
    golang.org/x/lint/golint \
    golang.org/x/tools/cmd/goimports \
    golang.org/x/tools/cmd/gorename \
    golang.org/x/tools/cmd/guru \
    golang.org/x/tools/gopls \
    golang.org/x/tools/gopls \
    honnef.co/go/tools/... 
RUN go get -u github.com/cosmtrek/air
CMD air -c .${BUILT_CMD}.air.conf

FROM base AS builder
ARG BUILT_CMD
ENV GIN_MODE=release
ENV PORT=80
COPY . /go/src/github.com/philippecarle/moood/${BUILT_CMD}
WORKDIR /go/src/github.com/philippecarle/moood/${BUILT_CMD}
RUN go build -tags=jsoniter -o /go/bin/${BUILT_CMD}

FROM scratch
ARG BUILT_CMD
ENV BUILT_CMD ${BUILT_CMD}
COPY --from=builder /go/bin/${BUILT_CMD} /go/bin/${BUILT_CMD}