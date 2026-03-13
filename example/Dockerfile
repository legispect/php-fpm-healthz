FROM gcr.io/google.com/cloudsdktool/cloud-sdk:slim as build-tool
WORKDIR /src
RUN git clone --depth=1 --branch main https://github.com/legispect/php-fpm-healthz .
WORKDIR /usr/local/bin
RUN curl -o /php-fpm-healthcheck https://raw.githubusercontent.com/renatomefi/php-fpm-healthcheck/master/php-fpm-healthcheck 

FROM golang:latest as go-build
WORKDIR /usr/src/app

COPY --from=build-tool /src/go.* ./
RUN go mod download
COPY --from=build-tool /src/*.go ./
RUN go build -o /health-server .

FROM php:fpm as app

RUN apt-get update && \
    apt-get install libfcgi-bin --no-install-recommends -qqy && \
    rm -rf /var/lib/apt/lists/*

COPY --chmod=755 --from=go-build /health-server /usr/local/bin/health-server
COPY --chmod=755 --from=build-tool /php-fpm-healthcheck /usr/local/bin/php-fpm-healthcheck

COPY --chmod=755 docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh

ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["php-fpm"]
