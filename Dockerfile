# Golang base image for builder
FROM golang:latest as goBuild

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -v -o build/lumi lumi.go

# Node base image for builder
FROM node:20 as nodeBuild

WORKDIR /app
COPY web ./web

WORKDIR /app/web
RUN npm install
RUN npx vite build . --outDir ./public --emptyOutDir

# Build final container
FROM gcr.io/distroless/static:latest

WORKDIR /app

COPY --from=goBuild /app/build .
COPY --from=nodeBuild /app/web/public ./public

EXPOSE 8080

CMD ["./lumi"]
