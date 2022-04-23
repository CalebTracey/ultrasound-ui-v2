FROM golang:1.18.1-alpine3.15 as builder
ADD . /app
WORKDIR /app/go-server/cmd/svr
RUN apk add git
RUN apk add build-base
RUN go mod download
RUN go build -o /main .

FROM node:16.14-buster AS node_builder
WORKDIR /code
COPY --from=builder /app .
RUN rm -rf /go-server/
RUN yarn install --frozen--lockfile
COPY . .
RUN yarn build

FROM alpine:3.15
RUN apk --no-cache add ca-certificates
COPY --from=builder /main .
COPY --from=builder /app/go-server/cmd/svr/config.json .
COPY --from=node_builder /code/build ./web
RUN chmod +x ./main
EXPOSE 80
CMD ./main

#NGINX web server
#FROM nginx:1.20-alpine AS prod
#COPY --from=builder /code/build /usr/share/nginx/html
#EXPOSE 80
#CMD ["nginx", "-g", "daemon off;"]



# Stage 2
#FROM alpine
#RUN adduser -S -D -H -h /app appuser
#USER appuser
#COPY --from=builder /code/build /app/
#WORKDIR /app
#CMD ["./server.go"]