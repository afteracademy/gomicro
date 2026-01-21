# Stage 1: Build the Go plugin
FROM golang:1.25.6-alpine AS builder

WORKDIR /plugins
COPY . .
WORKDIR /plugins/apikey_auth_plugin
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /plugins/apikey-auth-plugin main.go

# Stage 2: Build the Kong Docker image
FROM kong:3.9.1

COPY --from=builder /plugins/apikey-auth-plugin /usr/local/bin/apikey-auth-plugin

USER root
RUN chmod +x /usr/local/bin/apikey-auth-plugin
USER kong

COPY kong.yml /etc/kong/kong.yml

ENV KONG_DATABASE=off
ENV KONG_DECLARATIVE_CONFIG=/etc/kong/kong.yml
ENV KONG_PROXY_ACCESS_LOG=/dev/stdout
ENV KONG_ADMIN_ACCESS_LOG=/dev/stdout
ENV KONG_PROXY_ERROR_LOG=/dev/stderr
ENV KONG_ADMIN_ERROR_LOG=/dev/stderr
ENV KONG_PLUGINS=bundled,apikey-auth-plugin
ENV KONG_PLUGINSERVER_NAMES=apikey-auth-plugin

EXPOSE 8000 8443 8001 8444

CMD ["sh", "-c", "export KONG_PLUGINSERVER_APIKEY_AUTH_PLUGIN_QUERY_CMD='/usr/local/bin/apikey-auth-plugin -dump' && kong start -c /etc/kong/kong.yml"]