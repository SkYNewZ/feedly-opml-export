FROM --platform=$BUILDPLATFORM golang:1.15

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG SHA

WORKDIR /go/src/github.com/SkYNewZ/feedly-opml-export
COPY go.* ./

# Download dependencies
RUN go mod download

COPY . .

# Get final architecture and build
# linux/amd64, linux/arm/v7, linux/arm/v6, linux/arm64
RUN export GOOS=$(echo $TARGETPLATFORM | cut -d "/" -f1) && \
  export GOARCH=$(echo $TARGETPLATFORM | cut -d "/" -f2) && \
  echo "GOOS=$GOOS GOARCH=$GOARCH" && \
  export ARM=$(echo $TARGETPLATFORM | cut -d "/" -f3 | sed -e 's/v//') && \
  if [ ! -z "$ARM" ]; then export GOARM=$ARM && echo "GOARM=$GOARM"; fi && \
  CGO_ENABLED=0 go build -ldflags="-s -w -X 'main.version=${SHA}'" -o /feedly-opml-export .

FROM --platform=$BUILDPLATFORM scratch

ARG BUILD_DATE

LABEL maintainer="Quentin Lemaire <quentin@lemairepro.fr>"
LABEL org.label-schema.schema-version="1.0"
LABEL org.label-schema.build-date=${BUILD_DATE}
LABEL org.label-schema.name="feedly-opml-export"
LABEL org.label-schema.description="Simple Go script to print OPML export from Feedly"
LABEL org.label-schema.url="https://github.com/SkYNewZ/feedly-opml-export"

COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=0 /feedly-opml-export /feedly-opml-export

USER 1000:1000
ENTRYPOINT [ "/feedly-opml-export" ]
