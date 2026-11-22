# The demo frontend as static files behind nginx, for docker compose (the
# frontend service, on port 3000) and the cluster (deploy/k8s/50-frontend.yaml).
# While working on the page, `npm run dev` reloads on every save.
#
# nginx serves the built page and passes /api/... on to the ticketing
# service, as Vite's dev server does, so the browser sees one origin and
# CORS never applies. TICKETING_URL picks the service (default
# http://ticketing, its name in the cluster).
#
# Build from this directory:
#   docker buildx build --platform linux/amd64 -t ... .

# --- build ---
# The output is plain HTML/JS/CSS, the same for every platform: build on the
# builder's own platform, no emulation.
FROM --platform=$BUILDPLATFORM node:24.12.0-alpine3.23 AS build
WORKDIR /src
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
RUN npm run build

# --- runtime ---
# nginx-unprivileged runs as uid 101 and listens on 8080: no root needed.
FROM nginxinc/nginx-unprivileged:1.30.0-alpine
ENV TICKETING_URL=http://ticketing
# compose sets TICKETING_URL=http://ticketing:8080.
# The image's entrypoint fills ${TICKETING_URL} into the template at start.
COPY nginx/default.conf.template /etc/nginx/templates/default.conf.template
COPY --from=build /src/dist /usr/share/nginx/html
EXPOSE 8080
