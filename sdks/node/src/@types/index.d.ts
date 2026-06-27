export interface KubeRPCConfig {
  coreURL: string;
  serviceName: string;
  port?: number;
  host?: string;
}

export type Handler = (args: Record<string, any>) => any | Promise<any>;
