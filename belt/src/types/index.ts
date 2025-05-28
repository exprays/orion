// Type Interface Definitions for Belt


export interface ServerStats {
  uptime_in_seconds: number;
  uptime_in_days: number;
  used_memory: number;
  used_memory_human: string;
  keyspace_info: string;
  total_keys: number;
  connections: number;
  version: string;
  port: number;
}

export interface KeyInfo {
  key: string;
  type: 'string' | 'set' | 'hash';
  value?: string;
  members?: string[];
  fields?: Record<string, string>;
  ttl: number;
  size: number;
}

export interface ApiResponse<T> {
  success: boolean;
  data?: T;
  error?: string;
}

export interface CommandResult {
  command: string;
  result: string;
  timestamp: number;
  success: boolean;
}