import express from "express";
import { KubeRPC } from "kuberpc-node";
import dotenv from "dotenv";

dotenv.config();

const HTTP_PORT = Number(process.env.HTTP_PORT || 8081);

// In Kubernetes, KUBERPC_CORE_URL, KUBERPC_SERVICE_NAME, KUBERPC_HOST and KUBERPC_PORT
// are injected automatically via the admission webhook. Outside Kubernetes, set them
// in your .env file or pass them explicitly here.
const rpc = new KubeRPC();

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
      console.log(`[${process.env.KUBERPC_SERVICE_NAME ?? "server"}] ready - HTTP :${HTTP_PORT}  RPC :${process.env.KUBERPC_PORT ?? 7749}`);
    } catch (err) {
      console.error("failed to start:", err.message);
      process.exit(1);
    }
  });
