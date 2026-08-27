# build stage
FROM node:22-alpine AS build

WORKDIR /app

COPY package.json pnpm-lock.yaml ./

# Install pnpm
RUN npm install -g pnpm

RUN pnpm install

COPY . .
RUN pnpm run generate

# production stage
FROM nginx:stable-alpine AS production
COPY --from=build /app/.output/public /usr/share/nginx/html

# Copy the nginx configuration file
COPY ./nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 3000
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD wget -qO- http://127.0.0.1:3000/ >/dev/null || exit 1
CMD ["nginx", "-g", "daemon off;"]
