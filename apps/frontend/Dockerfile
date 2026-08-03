FROM node:24-alpine AS build

WORKDIR /workspace

RUN npm install --global pnpm@11.17.0

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
COPY apps/frontend/package.json apps/frontend/package.json
RUN pnpm install --frozen-lockfile --filter @avito-recap/frontend...

COPY contracts contracts
COPY apps/frontend apps/frontend
RUN pnpm --filter @avito-recap/frontend build

FROM caddy:2.11.4-alpine

COPY apps/frontend/Caddyfile /etc/caddy/Caddyfile
COPY --from=build /workspace/apps/frontend/dist /srv

EXPOSE 80
