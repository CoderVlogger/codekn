FROM golang:latest as builder
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o flash -v ./cmd/flash.go

FROM scratch
COPY --from=builder /app/flash /
ENV FLASH_PORT=8080
CMD [ "/flash" ]
