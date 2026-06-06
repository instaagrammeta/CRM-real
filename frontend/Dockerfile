# ============== build stage ==============
FROM node:22-alpine AS builder

WORKDIR /app

COPY package.json package-lock.json* ./
RUN npm ci || npm install

COPY . .

ENV NODE_ENV=production
RUN npm run build

# ============== runtime stage ==============
FROM node:22-alpine

RUN addgroup -S app && adduser -S app -G app

WORKDIR /app

COPY --from=builder --chown=app:app /app/.output ./.output
COPY --from=builder --chown=app:app /app/package.json ./

USER app

EXPOSE 3000

ENV HOST=0.0.0.0 PORT=3000 NODE_ENV=production

CMD ["node", ".output/server/index.mjs"]
