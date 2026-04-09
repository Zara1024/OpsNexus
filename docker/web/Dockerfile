FROM node:20-alpine AS builder

WORKDIR /src

COPY web/package.json web/package-lock.json ./web/
WORKDIR /src/web
RUN npm ci --no-audit --no-fund

COPY web/ ./
RUN npm run build

FROM nginx:1.27-alpine

RUN rm -rf /usr/share/nginx/html/*

COPY --from=builder /src/web/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
