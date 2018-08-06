/**
 * @leizm/slimiter-client
 *
 * @author Zongmin Lei <leizongmin@gmail.com>
 */

import * as net from "net";
import * as events from "events";
import * as assert from "assert";
import * as protocol from "./protocol";

export const DEFAULT_ADDRESS = "127.0.0.1:16789";

export interface ClientOptions {
  server: string;
  db: number;
}

export class Client extends events.EventEmitter {
  protected readonly socket: net.Socket;
  protected readonly db: number;
  protected buffer: Buffer = Buffer.alloc(0);
  protected callback: Record<string, Array<[Function, Function]>> = {
    [protocol.OpType.OpGetResult]: [],
    [protocol.OpType.OpIncrResult]: [],
  };

  constructor(options: Partial<ClientOptions> = {}) {
    super();
    const { host, port } = parseServerAddress(options.server || DEFAULT_ADDRESS);
    this.db = options.db || 0;
    assert(this.db >= 0, `invalid db number: ${this.db}`);
    this.socket = net.createConnection(port, host, () => {
      this.emit("connect");
    });
    this.socket.on("data", data => {
      // console.log("on data", data);
      this.buffer = Buffer.concat([this.buffer, data]);
      this.processBuffer();
    });
    this.socket.on("error", err => {
      this.emit("error", err);
    });
    this.socket.on("close", () => {
      this.emit("close");
    });
  }

  protected processBuffer() {
    try {
      while (this.buffer.length > 0) {
        const { data, rest } = protocol.unpack(this.buffer);
        this.buffer = rest;
        switch (data.op) {
          case protocol.OpType.OpGetResult:
            {
              const v = data.data.readUInt32BE(0);
              const fn = this.callback[protocol.OpType.OpGetResult].shift();
              if (fn) {
                fn[0](v);
              }
            }
            break;
          case protocol.OpType.OpIncrResult:
            {
              const v = data.data.readUInt32BE(0);
              const fn = this.callback[protocol.OpType.OpIncrResult].shift();
              if (fn) {
                fn[0](v);
              }
            }
            break;
          default:
          // unknown
        }
      }
    } catch (err) {
      // RangeError: Index out of range
      // console.log(err);
    }
  }

  protected send(op: protocol.OpType, data: Buffer): Promise<any> {
    return new Promise((resolve, reject) => {
      this.socket.write(
        protocol.pack({
          version: 1,
          op,
          length: data.length,
          data,
        }),
      );
      switch (op) {
        case protocol.OpType.OpGet:
          this.callback[protocol.OpType.OpGetResult].push([resolve, reject]);
          break;
        case protocol.OpType.OpIncr:
          this.callback[protocol.OpType.OpIncrResult].push([resolve, reject]);
          break;
      }
    });
  }

  public incr(ns: string, milliseconds: number, key: string, count: number = 1): Promise<number> {
    return this.send(
      protocol.OpType.OpIncr,
      protocol.packCmdArg({
        db: this.db,
        ns,
        milliseconds,
        key,
        count,
      }),
    );
  }

  public get(ns: string, milliseconds: number, key: string): Promise<number> {
    return this.send(
      protocol.OpType.OpGet,
      protocol.packCmdArg({
        db: this.db,
        ns,
        milliseconds,
        key,
        count: 0,
      }),
    );
  }

  public close() {
    this.socket.end();
  }
}

function parseServerAddress(str: string): { host: string; port: number } {
  const b = str.split(":");
  assert(b.length === 2, `invalid server address format: ${str}`);
  const port = Number(b[1]);
  assert(port > 0, `invalid server address format: ${str}`);
  return { host: b[0], port };
}

async function test() {
  function sleep(ms: number) {
    return new Promise(resolve => {
      setTimeout(resolve, ms);
    });
  }
  const c = new Client();
  // console.log(c);
  for (let i = 0; i < Number.MAX_SAFE_INTEGER; i++) {
    c.incr("hello", 1000, "www").then(console.log);
    await sleep(0);
    // c.incr("hello", 199, "www2", 2).then(console.log);
    // c.incr("hello", 199, "www2", 5).then(console.log);
  }
}
test().catch(console.log);
