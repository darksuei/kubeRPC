import express from "express";
import { KubeRPC } from "kuberpc-node";
import dotenv from "dotenv";

dotenv.config();

const SERVICE_NAME = process.env.SERVICE_NAME || "server";
const HTTP_PORT = Number(process.env.HTTP_PORT || 8081);

const rpc = new KubeRPC({
  coreURL: process.env.KUBERPC_CORE_URL,
  serviceName: SERVICE_NAME,
  host: process.env.RPC_HOST || "localhost",
});

function fib(n) {
  if (n <= 1) return n;
  let a = 0, b = 1;
  for (let i = 2; i <= n; i++) [a, b] = [b, a + b];
  return b;
}

express()
  .get("/", (_, res) => res.send("ok"))
  .listen(HTTP_PORT, async () => {
    try {
      await rpc.register("fib", async ({ n }) => fib(n));
      console.log(`[${SERVICE_NAME}] ready — HTTP :${HTTP_PORT}  RPC :7749`);
    } catch (err) {
      console.error(`[${SERVICE_NAME}] failed to start:`, err.message);
      process.exit(1);
    }
  });
