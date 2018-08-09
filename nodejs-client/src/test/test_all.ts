/**
 * @leizm/slimiter-client
 *
 * @author Zongmin Lei <leizongmin@gmail.com>
 */

import { expect } from "chai";
import { IPackage, ICmdArg, pack, unpack, packCmdArg, unpackCmdArg } from "../lib/protocol";
import { Client } from "../lib/client";

function sleep(ms: number) {
  return new Promise(resolve => {
    setTimeout(resolve, ms);
  });
}

describe("protocol", function() {
  it("package", function() {
    const d: IPackage = { version: 1, id: 999, op: 2, length: 3, data: Buffer.from("abc") };
    const b = pack(d);
    const p = unpack(Buffer.concat([b, Buffer.from("xyz")]));
    expect(p.data).to.deep.equal(d);
    expect(p.rest).to.deep.equal(Buffer.from("xyz"));
  });

  it("arg", function() {
    const d: ICmdArg = { db: 1, ns: "hello", milliseconds: 100, key: "world", count: 2 };
    const b = packCmdArg(d);
    const p = unpackCmdArg(b);
    expect(d).to.deep.equal(p);
  });
});

describe("client", function() {
  it("ping", async function() {
    this.timeout(10000);
    const c = new Client();
    expect(await c.ping()).to.gte(0);
    await sleep(1);
    c.close();
  });

  it("get & incr", async function() {
    this.timeout(10000);
    const c = new Client();
    expect(await c.incr("hello", "www2", 199, 2)).to.equal(2);
    expect(await c.incr("hello", "www2", 200, 1)).to.equal(1);
    expect(await c.incr("hello", "www2", 199, 2)).to.equal(4);
    expect(await c.incr("hello", "www2", 199, 1)).to.equal(5);
    c.close();
  });

  it("many", async function() {
    this.timeout(10000);
    const c = new Client();
    const list = [];
    const N = 10000;
    for (let i = 0; i < N; i++) {
      list.push(c.incr("hello", "www", 1000));
    }
    const ret = await Promise.all(list);
    expect(ret.length).to.equal(N);
  });
});
