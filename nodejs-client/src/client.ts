/**
 * @leizm/slimiter-client
 *
 * @author Zongmin Lei <leizongmin@gmail.com>
 */

import * as net from "net";
import * as events from "events";
import * as assert from "assert";

export const DEFAULT_ADDRESS = "127.0.0.1:16789";

export interface ClientOptions {
  server: string;
  db: number;
}

export class Client extends events.EventEmitter {
  protected readonly socket: net.Socket;
  protected readonly db: number;

  constructor(options: Partial<ClientOptions> = {}) {
    super();
    const { host, port } = parseServerAddress(options.server || DEFAULT_ADDRESS);
    this.db = options.db || 0;
    assert(this.db > 0, `invalid db number: ${this.db}`);
    this.socket = net.createConnection(port, host, () => {
      this.emit("connect");
    });
    this.socket.on("data", data => {});
    this.socket.on("error", err => {
      this.emit("error", err);
    });
    this.socket.on("close", () => {
      this.emit("close");
    });
  }
}

function parseServerAddress(str: string): { host: string; port: number } {
  const b = str.split(":");
  assert(b.length === 2, `invalid server address format: ${str}`);
  const port = Number(b[1]);
  assert(port > 0, `invalid server address format: ${str}`);
  return { host: b[0], port };
}

interface IPackage {
  version: number;
  op: number;
  length: number;
  data: Buffer;
}

function pack(data: IPackage): Buffer {
  const buf = Buffer.alloc(5 + data.length);
  buf.writeUInt16BE(data.version, 0);
  buf.writeUInt8(data.op, 2);
  buf.writeUInt16BE(data.length, 3);
  data.data.copy(buf, 5);
  return buf;
}

function unpack(buf: Buffer): { data: IPackage; rest: Buffer } {
  const version = buf.readUInt16BE(0);
  const op = buf.readUInt8(2);
  const length = buf.readUInt16BE(3);
  const data = buf.slice(5, 5 + length);
  const rest = buf.slice(5 + length);
  return { data: { version, op, length, data }, rest };
}

{
  const d: IPackage = { version: 1, op: 2, length: 3, data: Buffer.from("abc") };
  const b = pack(d);
  const p = unpack(Buffer.concat([b, Buffer.from("xyz")]));
  assert.deepEqual(p.data, d);
  assert.deepEqual(p.rest, Buffer.from("xyz"));
}

interface ICmdArg {
  db: number;
  ns: string;
  milliseconds: number;
  key: string;
  count: number;
}

function packCmdArg(a: ICmdArg): Buffer {
  const ns = Buffer.from(a.ns);
  const key = Buffer.from(a.key);
  const buf = Buffer.alloc(4 + 1 + ns.length + 4 + 1 + key.length + 4);
  buf.writeUInt32BE(a.db, 0);
  buf.writeUInt8(ns.length, 4);
  ns.copy(buf, 5);
  buf.writeUInt32BE(a.milliseconds, 5 + ns.length);
  buf.writeUInt8(key.length, 5 + ns.length + 4);
  key.copy(buf, 5 + ns.length + 4 + 1);
  buf.writeUInt32BE(a.count, 5 + ns.length + 4 + 1 + key.length);
  return buf;
}

function unpackCmdArg(buf: Buffer): ICmdArg {
  const db = buf.readUInt32BE(0);
  const nsLen = buf.readUInt8(4);
  const ns = buf.slice(5, 5 + nsLen);
  const milliseconds = buf.readUInt32BE(5 + nsLen);
  const keyLen = buf.readUInt8(5 + nsLen + 4);
  const key = buf.slice(5 + nsLen + 5, 5 + nsLen + 5 + keyLen);
  const count = buf.readUInt32BE(5 + nsLen + 5 + keyLen);
  return { db, ns: ns.toString(), milliseconds, key: key.toString(), count };
}

{
  const d: ICmdArg = { db: 1, ns: "hello", milliseconds: 100, key: "world", count: 2 };
  const b = packCmdArg(d);
  const p = unpackCmdArg(b);
  assert.deepEqual(d, p);
}
