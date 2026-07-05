# syntax=docker/dockerfile:1.7

# ---------- Build stage ----------
FROM node:20-alpine AS builder

WORKDIR /app

# Cache npm layer
COPY package.json package-lock.json ./
RUN --mount=type=cache,target=/root/.npm \
    npm ci

# Build the Vite production bundle
ARG VITE_API_BASE_URL=https://project-alpha-bt-224586-8-1326559147.sh.run.tcloudbase.com/api/v1
ENV VITE_API_BASE_URL=${VITE_API_BASE_URL}
COPY . .
RUN npm run build

# ---------- Runtime stage ----------
FROM nginx:1.27-alpine AS runtime

# Drop default config, ship ours
RUN rm /etc/nginx/conf.d/default.conf
COPY nginx.conf /etc/nginx/conf.d/app.conf

# Static assets built in the previous stage
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget -qO- http://127.0.0.1/healthz || exit 1

CMD ["nginx", "-g", "daemon off;"]
