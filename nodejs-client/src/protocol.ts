/**
 * @leizm/slimiter-client
 *
 * @author Zongmin Lei <leizongmin@gmail.com>
 */

import * as assert from "assert";

export enum OpType {
  OpPing = 0x1,
  OpPong = 0x2,
  OpGet = 0x3,
  OpGetResult = 0x4,
  OpIncr = 0x5,
  OpIncrResult = 0x6,
}

export interface IPackage {
  version: number;
  op: OpType;
  length: number;
  data: Buffer;
}

export function pack(data: IPackage): Buffer {
  const buf = Buffer.alloc(5 + data.length);
  buf.writeUInt16BE(data.version, 0);
  buf.writeUInt8(data.op, 2);
  buf.writeUInt16BE(data.length, 3);
  data.data.copy(buf, 5);
  return buf;
}

export function unpack(buf: Buffer): { data: IPackage; rest: Buffer } {
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

export interface ICmdArg {
  db: number;
  ns: string;
  milliseconds: number;
  key: string;
  count: number;
}

export function packCmdArg(a: ICmdArg): Buffer {
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

export function unpackCmdArg(buf: Buffer): ICmdArg {
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
