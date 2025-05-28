FROM golang:1.24.3-alpine3.21 AS build

RUN addgroup -S app && adduser -S app -G app

ARG APPLICATION_DIRECTORY=applications/rest-api
COPY ${APPLICATION_DIRECTORY} /build/applications/rest-api

ARG LIBRARIES_DIRECTORY=libraries
COPY ${LIBRARIES_DIRECTORY} /build/libraries

WORKDIR /build/applications/rest-api
RUN go build -o /rest-api

FROM scratch

COPY --from=build /rest-api /rest-api

COPY --from=build /etc/passwd /etc/passwd
USER app

ENTRYPOINT ["/rest-api"]
