# Compile
FROM golang:latest as go

WORKDIR /app
ADD . /app
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .


# Rebuild with only the golang binary
FROM scratch

WORKDIR /root
COPY --from=go /app/main /root/main
ENV GIN_MODE release
EXPOSE 8080
CMD ["/root/main"]
