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

  constructor(options: Partial<ClientOptions> = {}) {
    super();
    const { host, port } = parseServerAddress(options.server || DEFAULT_ADDRESS);
    this.db = options.db || 0;
    assert(this.db >= 0, `invalid db number: ${this.db}`);
    this.socket = net.createConnection(port, host, () => {
      this.emit("connect");
      this.send(protocol.OpType.OpPing, Buffer.alloc(0));
      this.send(
        protocol.OpType.OpIncr,
        protocol.packCmdArg({
          db: this.db,
          ns: "hello",
          milliseconds: 299,
          key: "woaa",
          count: 1,
        }),
      );
    });
    this.socket.on("data", data => {
      console.log(data);
    });
    this.socket.on("error", err => {
      this.emit("error", err);
    });
    this.socket.on("close", () => {
      this.emit("close");
    });
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
    });
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
  const c = new Client();
  console.log(c);
}
test().catch(console.log);
