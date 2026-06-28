# Build context must be sdks/node/ so the local SDK reference resolves correctly.
# docker build -f samples/client/Dockerfile -t docker.io/darksueii/kuberpc-sample-client:2.0.0 .

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

WORKDIR /sdk/samples/client
COPY samples/client/package.json samples/client/yarn.lock ./
RUN yarn install --frozen-lockfile
COPY samples/client/index.js ./

CMD ["node", "index.js"]
