# Build context must be sdks/node/ so the local SDK reference resolves correctly.
# docker build -f samples/server/Dockerfile -t docker.io/darksueii/kuberpc-sample-server:2.0.0 .

FROM node:20-alpine AS sdk-builder
WORKDIR /sdk
COPY package.json yarn.lock ./
RUN yarn install --frozen-lockfile
COPY tsconfig.json ./
COPY src/ ./src/
RUN yarn build

FROM node:20-alpine
# Replicate the layout that package.json's "file:../../" expects:
#   /sdk/           ← SDK root (dist/ lives here)
#   /sdk/samples/server/ ← sample working dir
WORKDIR /sdk
COPY --from=sdk-builder /sdk/dist ./dist
COPY --from=sdk-builder /sdk/node_modules ./node_modules
COPY package.json ./

WORKDIR /sdk/samples/server
COPY samples/server/package.json samples/server/yarn.lock ./
RUN yarn install --frozen-lockfile
COPY samples/server/index.js ./

CMD ["node", "index.js"]
