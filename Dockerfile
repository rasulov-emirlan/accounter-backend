#use small build environment
FROM golang:1.19-alpine3.17 as build 

WORKDIR /app

#get dependencies
COPY go.mod ./
RUN go mod download && go mod tidy

#copy and build the project
COPY . ./
RUN go build -o /apiserver ./cmd/apiserver/main.go

#Use deploy env
FROM alpine:latest

WORKDIR /

COPY --from=build /apiserver /apiserver

#optionally use --expose for dynamic config
EXPOSE 8080

ENTRYPOINT ["./apiserver"]