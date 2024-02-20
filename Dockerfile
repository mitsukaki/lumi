# Build the go binary
FROM golang:latest as goBuild

WORKDIR /go/src/app

COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build -v -o /go/bin/lumi lumi.go

# Bundle the javascript and build the web ui
FROM node:20 as nodeBuild

WORKDIR /app
COPY web ./web

WORKDIR /app/web
RUN npm install
RUN npx vite build . --outDir ./public --emptyOutDir

# Build final container
FROM gcr.io/distroless/base-debian12

COPY --from=goBuild /go/bin/lumi /
COPY --from=nodeBuild /app/web/public /public

EXPOSE 8080

CMD ["/lumi"]
