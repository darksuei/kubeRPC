export class KubeRpcError extends Error {
  constructor(
    public readonly code: string,
    message: string,
  ) {
    super(message);
    this.name = "KubeRpcError";
    Object.setPrototypeOf(this, KubeRpcError.prototype);
  }

  static methodNotFound(method: string, service: string) {
    return new KubeRpcError(
      "METHOD_NOT_FOUND",
      `Method "${method}" not found in service "${service}"`,
    );
  }

  static connectionFailed(host: string, port: number) {
    return new KubeRpcError(
      "CONNECTION_FAILED",
      `Failed to connect to ${host}:${port}`,
    );
  }

  static coreUnreachable(url: string) {
    return new KubeRpcError(
      "CORE_UNREACHABLE",
      `kubeRPC core unreachable at ${url}`,
    );
  }
}
