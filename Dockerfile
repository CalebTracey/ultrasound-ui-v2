## build environment
#FROM node:16.13.0-alpine AS build
#WORKDIR /usr/src/app
## Copies package.json and package-lock.json to Docker environment
#COPY package*.json ./
## Installs all node packages
#RUN yarn install
## Copies everything over to Docker environment
#COPY . .
## Uses port which is used by the actual application
#EXPOSE $PORT
## Finally runs the application
#CMD [ "yarn", "start" ]

FROM node:16.13.0-alpine AS builder

WORKDIR /opt/web
COPY package.json yarn.lock ./
RUN npm install

ENV PATH="./node_modules/.bin:$PATH"

COPY . ./
RUN npm run build

FROM nginx:1.17-alpine
RUN apk --no-cache add curl
RUN curl -L https://github.com/a8m/envsubst/releases/download/v1.1.0/envsubst-`uname -s`-`uname -m` -o envsubst && \
    chmod +x envsubst && \
    mv envsubst /usr/local/bin
COPY ./nginx.conf /etc/nginx/nginx.template
CMD ["/bin/sh", "-c", "envsubst < /etc/nginx/nginx.template > /etc/nginx/conf.d/default.conf && nginx -g 'daemon off;'"]
COPY --from=builder /opt/web/build /usr/share/nginx/html