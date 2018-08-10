/**
 * @leizm/yedda
 *
 * @author Zongmin Lei <leizongmin@gmail.com>
 */

export enum OpType {
  OpPing = 0x1,
  OpPong = 0x2,
  OpGet = 0x3,
  OpGetResult = 0x4,
  OpIncr = 0x5,
  OpIncrResult = 0x6,
}

export interface IPackage {
  // uint16
  version: number;
  // uint32
  id: number;
  // uint16
  op: OpType;
  // uint16
  length: number;
  // buffer
  data: Buffer;
}

export function pack(data: IPackage): Buffer {
  const buf = Buffer.alloc(10 + data.length);
  buf.writeUInt16BE(data.version, 0);
  buf.writeUInt32BE(data.id, 2);
  buf.writeUInt16BE(data.op, 6);
  buf.writeUInt16BE(data.length, 8);
  data.data.copy(buf, 10);
  return buf;
}

export function unpack(buf: Buffer): { ok: boolean; data: IPackage; rest: Buffer } {
  const version = buf.readUInt16BE(0);
  const id = buf.readUInt32BE(2);
  const op = buf.readUInt16BE(6);
  const length = buf.readUInt16BE(8);
  const data = buf.slice(10, 10 + length);
  const rest = buf.slice(10 + length);
  const ok = data.length === length;
  return { ok, data: { version, id, op, length, data }, rest };
}

export interface ICmdArg {
  // uint32
  db: number;
  // uint8 + buffer
  ns: string;
  // uint32
  milliseconds: number;
  // uint8 + buffer
  key: string;
  // uint32
  count: number;
}

export function packCmdArg(a: ICmdArg): Buffer {
  const ns = Buffer.from(a.ns);
  const key = Buffer.from(a.key);
  const buf = Buffer.alloc(14 + ns.length + key.length);
  buf.writeUInt32BE(a.db, 0);
  buf.writeUInt8(ns.length, 4);
  ns.copy(buf, 5);
  buf.writeUInt32BE(a.milliseconds, 5 + ns.length);
  buf.writeUInt8(key.length, 9 + ns.length);
  key.copy(buf, 10 + ns.length);
  buf.writeUInt32BE(a.count, 10 + ns.length + key.length);
  return buf;
}

export function unpackCmdArg(buf: Buffer): ICmdArg {
  const db = buf.readUInt32BE(0);
  const nsLen = buf.readUInt8(4);
  const ns = buf.slice(5, 5 + nsLen);
  const milliseconds = buf.readUInt32BE(5 + nsLen);
  const keyLen = buf.readUInt8(9 + nsLen);
  const key = buf.slice(10 + nsLen, 10 + nsLen + keyLen);
  const count = buf.readUInt32BE(10 + nsLen + keyLen);
  return { db, ns: ns.toString(), milliseconds, key: key.toString(), count };
}

export function uint64ToBuffer(n: number): Buffer {
  const b = Buffer.alloc(8);
  b.writeUIntBE(n, 0, 8);
  return b;
}
