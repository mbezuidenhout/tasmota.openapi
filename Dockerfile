FROM golang:1.17 AS build
WORKDIR /go/src
COPY src/go ./go
COPY src/main.go .
COPY src/go.mod .
COPY src/go.sum .

ENV CGO_ENABLED=0
RUN go get -d -v ./...

RUN go build -a -installsuffix cgo -o swagger .

FROM scratch AS runtime
COPY --from=build /go/src/swagger ./
EXPOSE 8080/tcp
ENTRYPOINT ["./swagger"]
