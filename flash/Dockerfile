FROM golang:latest as builder
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o flash -v ./cmd/flash.go

FROM alpine
COPY --from=builder /app/flash /
COPY --from=builder /app/resources /resources
ENV FLASH_PORT=8080
CMD [ "/flash" ]
