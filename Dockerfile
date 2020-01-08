FROM golang:1.13-alpine3.11 AS backend-builder
RUN apk --no-cache add libvirt-dev build-base
COPY ./backend /app
RUN set -ex && cd /app && go build

FROM node:13-alpine3.11 AS frontend-builder
COPY ./frontend /app
RUN set -ex && cd /app && yarn && yarn build --modern

FROM alpine:3.11
RUN apk --no-cache add libvirt openssh-client cdrkit
COPY --from=backend-builder /app/backend /maimok
COPY --from=frontend-builder /app/dist/ /dist
COPY ./backend/templates /templates

CMD ["/maimok"]
