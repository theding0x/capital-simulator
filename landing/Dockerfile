# syntax=docker/dockerfile:1.7
# Build context: repo root (needs LETTER.md). Build with:
#   docker build -f landing/Dockerfile -t daskap-landing .
# Selective COPY (matching the Go service Dockerfiles) keeps node_modules/dist
# out of the image even though the context is the whole repo.
FROM node:20-alpine AS build
WORKDIR /src
COPY landing/package.json landing/package-lock.json* ./
RUN if [ -f package-lock.json ]; then npm ci; else npm install; fi
COPY landing/gen.mjs landing/template.html landing/styles.css ./
COPY landing/public ./public
COPY LETTER.md /LETTER.md
RUN node gen.mjs

FROM nginx:1.27-alpine
COPY --from=build /src/dist /usr/share/nginx/html
COPY landing/nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
