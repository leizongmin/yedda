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
    [protocol.OpType.OpPong]: [],
    [protocol.OpType.OpGetResult]: [],
    [protocol.OpType.OpIncrResult]: [],
  };
  protected id: number = 0;

  constructor(options: Partial<ClientOptions> = {}) {
    super();
    const { host, port } = parseServerAddress(options.server || DEFAULT_ADDRESS);
    this.db = options.db || 0;
    assert(this.db >= 0, `invalid db number: ${this.db}`);
    this.socket = net.createConnection(port, host, () => {
      this.emit("connect");
    });
    this.socket.on("data", data => {
      // console.log("data", data);
      this.buffer = Buffer.concat([this.buffer, data]);
      // console.log("    ", this.buffer);
      this.processBuffer();
    });
    this.socket.on("error", err => {
      this.emit("error", err);
    });
    this.socket.on("close", () => {
      this.emit("close");
    });
    this.socket.setNoDelay(true);
  }

  protected processBuffer() {
    try {
      while (this.buffer.length > 0) {
        const { ok, data, rest } = protocol.unpack(this.buffer);
        if (!ok) break;
        this.buffer = rest;
        switch (data.op) {
          case protocol.OpType.OpPing:
            this.send(protocol.OpType.OpPong, data.data);
            break;
          case protocol.OpType.OpPong:
            this.runCallback(protocol.OpType.OpPong, Date.now() - data.data.readUIntBE(0, 8));
            break;
          case protocol.OpType.OpGetResult:
            this.runCallback(protocol.OpType.OpGetResult, data.data.readUInt32BE(0));
            break;
          case protocol.OpType.OpIncrResult:
            this.runCallback(protocol.OpType.OpIncrResult, data.data.readUInt32BE(0));
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
      const id = this.id++;
      const buf = protocol.pack({
        version: 1,
        id,
        op,
        length: data.length,
        data,
      });
      // console.log("send", buf);
      this.socket.write(buf);
      switch (op) {
        case protocol.OpType.OpPing:
          this.callback[protocol.OpType.OpPong].push([resolve, reject]);
          break;
        case protocol.OpType.OpGet:
          this.callback[protocol.OpType.OpGetResult].push([resolve, reject]);
          break;
        case protocol.OpType.OpIncr:
          this.callback[protocol.OpType.OpIncrResult].push([resolve, reject]);
          break;
      }
    });
  }

  protected runCallback(op: protocol.OpType, v: any): boolean {
    const fn = this.callback[op].shift();
    if (fn) {
      fn[0](v);
      return true;
    }
    return false;
  }

  public ping(): Promise<number> {
    return this.send(protocol.OpType.OpPing, protocol.uint64ToBuffer(Date.now()));
  }

  public incr(ns: string, key: string, milliseconds: number, count: number = 1): Promise<number> {
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

  public get(ns: string, key: string, milliseconds: number): Promise<number> {
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
