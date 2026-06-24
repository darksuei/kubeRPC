import net from "net";
import axios, { AxiosInstance } from "axios";
import { encode, decode } from "@msgpack/msgpack";
import { Handler, KubeRPCConfig } from "./@types";
import { KubeRpcError } from "./errors";

// Wire format (positional arrays to avoid encoding key strings on every call):
//   request  → [method: string, args: object]
//   response → [null, result]  on success
//   response → [errorMsg: string]  on error

export class KubeRPC {
  private http: AxiosInstance;
  private config: Required<KubeRPCConfig>;
  private handlers = new Map<string, Handler>();
  private server: net.Server | null = null;
  private pool = new Map<string, net.Socket>();
  private endpointCache = new Map<string, { host: string; port: number }>();
  private locks = new Map<string, Promise<unknown>>();

  constructor({ coreURL, serviceName, port, host = "0.0.0.0" }: KubeRPCConfig) {
    this.config = { coreURL, serviceName, port, host };
    this.http = axios.create({ baseURL: coreURL });
  }

  async register(name: string, handler: Handler): Promise<void> {
    this.handlers.set(name, handler);

    if (!this.server) {
      await this.startServer();
    }

    await this.http.post("/register-methods", {
      service_name: this.config.serviceName,
      methods: [{ name, params: [], description: "" }],
    });
  }

  async call(
    service: string,
    method: string,
    args: Record<string, any> = {},
  ): Promise<any> {
    const endpoint = await this.resolve(service, method);
    const socket = await this.getConnection(service, endpoint);
    return this.enqueue(service, socket, method, args);
  }

  close(): void {
    for (const socket of this.pool.values()) socket.destroy();
    this.pool.clear();
    this.locks.clear();
    this.endpointCache.clear();
    if (this.server) this.server.close();
  }

  private async startServer(): Promise<void> {
    await this.http.put(`/update-service?name=${this.config.serviceName}`, {
      host: this.config.host,
      port: this.config.port,
    });

    return new Promise((resolve, reject) => {
      this.server = net.createServer((socket) => {
        socket.setNoDelay(true);
        socket.setKeepAlive(true, 0);
        let buf = Buffer.alloc(0);

        socket.on("data", (chunk) => {
          buf = Buffer.concat([buf, chunk]);
          this.drainInbound(socket, buf).then((remaining) => {
            buf = remaining;
          });
        });

        socket.on("error", () => socket.destroy());
      });

      this.server.listen(this.config.port, () => resolve());
      this.server.once("error", reject);
    });
  }

  private async drainInbound(socket: net.Socket, buf: Buffer): Promise<Buffer> {
    while (buf.length >= 4) {
      const len = buf.readUInt32BE(0);
      if (buf.length < 4 + len) break;
      const frame = buf.slice(4, 4 + len);
      buf = buf.slice(4 + len);
      await this.handleInboundFrame(socket, frame);
    }
    return buf;
  }

  private async handleInboundFrame(socket: net.Socket, frame: Buffer): Promise<void> {
    try {
      const [method, args] = decode(frame) as [string, any];
      const handler = this.handlers.get(method);

      if (!handler) {
        this.writeFrame(socket, [`Method "${method}" not found`]);
        return;
      }

      const result = await handler(args);
      this.writeFrame(socket, [null, result]);
    } catch (err: any) {
      this.writeFrame(socket, [err?.message ?? "Internal error"]);
    }
  }

  private writeFrame(socket: net.Socket, payload: unknown): void {
    const body = Buffer.from(encode(payload));
    const frame = Buffer.alloc(4 + body.length);
    frame.writeUInt32BE(body.length, 0);
    body.copy(frame, 4);
    socket.write(frame);
  }

  private async resolve(
    service: string,
    method: string,
  ): Promise<{ host: string; port: number }> {
    const cached = this.endpointCache.get(service);
    if (cached) return cached;

    const { data } = await this.http
      .get(`/get-method?name=${service}&method=${method}`)
      .catch(() => {
        throw KubeRpcError.methodNotFound(method, service);
      });

    if (!data.host || !data.port) {
      throw KubeRpcError.methodNotFound(method, service);
    }

    const endpoint = { host: data.host, port: Number(data.port) };
    this.endpointCache.set(service, endpoint);
    return endpoint;
  }

  private async getConnection(
    service: string,
    endpoint: { host: string; port: number },
  ): Promise<net.Socket> {
    const existing = this.pool.get(service);
    if (existing && !existing.destroyed) return existing;

    const socket = await new Promise<net.Socket>((resolve, reject) => {
      const s = new net.Socket();
      s.setNoDelay(true);
      s.setKeepAlive(true, 0);
      s.connect(endpoint.port, endpoint.host, () => resolve(s));
      s.once("error", () =>
        reject(KubeRpcError.connectionFailed(endpoint.host, endpoint.port)),
      );
    });

    this.pool.set(service, socket);
    socket.once("close", () => {
      this.pool.delete(service);
      this.locks.delete(service);
      this.endpointCache.delete(service);
    });
    socket.once("error", () => {
      this.pool.delete(service);
      this.locks.delete(service);
    });

    return socket;
  }

  private enqueue(
    service: string,
    socket: net.Socket,
    method: string,
    args: Record<string, any>,
  ): Promise<any> {
    const tail = (this.locks.get(service) ?? Promise.resolve()).then(() =>
      this.sendFrame(socket, method, args),
    );
    this.locks.set(service, tail.catch(() => {}));
    return tail;
  }

  private sendFrame(
    socket: net.Socket,
    method: string,
    args: Record<string, any>,
  ): Promise<any> {
    return new Promise((resolve, reject) => {
      const body = Buffer.from(encode([method, args]));
      const frame = Buffer.alloc(4 + body.length);
      frame.writeUInt32BE(body.length, 0);
      body.copy(frame, 4);

      let buf = Buffer.alloc(0);

      const onData = (chunk: Buffer) => {
        buf = Buffer.concat([buf, chunk]);
        if (buf.length < 4) return;
        const len = buf.readUInt32BE(0);
        if (buf.length < 4 + len) return;
        socket.removeListener("data", onData);
        const r = decodeResponse(buf.slice(4, 4 + len));
        if (r.ok) resolve(r.value);
        else reject(new Error(r.error));
      };

      socket.on("data", onData);
      socket.once("error", reject);
      socket.write(frame);
    });
  }
}

function decodeResponse(buf: Buffer): { ok: true; value: any } | { ok: false; error: string } {
  const raw = decode(buf);
  if (!Array.isArray(raw)) {
    return { ok: false, error: `Protocol mismatch: expected array, got ${typeof raw}. Restart the server.` };
  }
  if (raw.length === 1) return { ok: false, error: raw[0] as string };
  return { ok: true, value: raw[1] };
}
