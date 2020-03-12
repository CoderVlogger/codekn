FROM golang:alpine as builder
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o profx -v ./cmd/profx.go

FROM alpine
COPY --from=builder /app/profx /
#RUN apk --no-cache add ca-certificates
CMD [ "/profx" ]
