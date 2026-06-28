# Build context must be sdks/node/ so the local SDK reference resolves correctly.
# docker build -f benchmark/server/Dockerfile -t docker.io/darksueii/kuberpc-bench-server:2.0.0 .

FROM node:20-alpine AS sdk-builder
WORKDIR /sdk
COPY package.json yarn.lock ./
RUN yarn install --frozen-lockfile
COPY tsconfig.json ./
COPY src/ ./src/
RUN yarn build

FROM node:20-alpine
WORKDIR /sdk
COPY --from=sdk-builder /sdk/dist ./dist
COPY --from=sdk-builder /sdk/node_modules ./node_modules
COPY package.json ./

WORKDIR /sdk/benchmark/server
COPY benchmark/server/package.json benchmark/server/yarn.lock ./
RUN yarn install --frozen-lockfile
COPY benchmark/server/index.js ./

CMD ["node", "index.js"]
