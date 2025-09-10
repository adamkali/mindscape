
FROM node:22-alpine as node_builder
WORKDIR /usr/src/web
COPY web/package.json ./

RUN npm install -g pnpm && pnpm install
COPY web/ ./
RUN pnpm run build

FROM golang:1.24-alpine as go_builder

WORKDIR /usr/src
COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o mindscape .

# We now attatch config/ diectory to the build
# but because .dockerignore is set to ignore it
# we will use a volume to mount it
VOLUME /usr/src/config

# Copy the executable to the final image
FROM alpine:latest as app

EXPOSE 60000
WORKDIR /app

## If you are not using React you can comment out this section
COPY --from=node_builder /usr/src/web/dist /app/web/dist

COPY --from=go_builder /usr/src/mindscape /app/
CMD ["/app/mindscape", "-e", "development"]
