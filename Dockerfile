FROM node:16.14.2-buster AS builder

WORKDIR /code

COPY package.json package.json
COPY yarn.lock yarn.lock

RUN yarn install --frozen-lockfile

COPY . .

RUN yarn build

#NGINX web server
FROM nginx:1.20-alpine AS prod

COPY --from=builder /code/build /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]