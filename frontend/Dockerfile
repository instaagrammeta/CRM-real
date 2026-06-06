# =====================================================================
# Multi-stage build for the Nuxt 3 frontend.
# =====================================================================

# ============== build stage ==============
FROM node:22-alpine AS builder

WORKDIR /app

# Cache deps
COPY package.json package-lock.json* ./
# `npm ci` is preferred when package-lock.json exists; fall back to install otherwise.
RUN if [ -f package-lock.json ]; then npm ci; else npm install; fi

# Copy source and build
COPY . .

ENV NODE_ENV=production
RUN npm run build

# ============== runtime stage ==============
FROM node:22-alpine

# wget for healthchecks (busybox-included in alpine).
RUN addgroup -S app && adduser -S app -G app

WORKDIR /app

COPY --from=builder --chown=app:app /app/.output ./.output
COPY --from=builder --chown=app:app /app/package.json ./

USER app

ENV HOST=0.0.0.0 \
    PORT=3000 \
    NODE_ENV=production \
    NITRO_HOST=0.0.0.0 \
    NITRO_PORT=3000

EXPOSE 3000

HEALTHCHECK --interval=20s --timeout=5s --retries=3 --start-period=15s \
  CMD wget -qO- http://127.0.0.1:3000/ >/dev/null 2>&1 || exit 1

CMD ["node", ".output/server/index.mjs"]
