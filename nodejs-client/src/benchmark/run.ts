import { Client } from "../lib/client";
import Benchmark from "@leizm/benchmark";

async function main() {
  const bench = new Benchmark({ title: "test slimiter", seconds: 10, delay: 1 });
  const c = new Client();

  const namespaces = "+"
    .repeat(10000)
    .split("")
    .map(_ => Math.random().toString(36));
  const keys = "+"
    .repeat(1000000)
    .split("")
    .map(_ => Math.random().toString(36) + Date.now());
  function getRandomNs() {
    const i = Math.floor(Math.random() * namespaces.length);
    return namespaces[i];
  }
  function getRandomKey() {
    const i = Math.floor(Math.random() * keys.length);
    return keys[i];
  }
  function getRandomMilliseconds() {
    return Math.floor(Math.floor(Math.random() * 100) * 10);
  }

  bench.addAsync("PING", () => c.ping());
  bench.addAsync("INCR (single ns & key)", () => c.incr("abc", "xxxx", 100));
  bench.addAsync("GET (single ns & key)", () => c.get("abc", "xxxx", 100));
  bench.addAsync(`INCR (${namespaces.length} ns & ${keys.length} key)`, () =>
    c.incr(getRandomNs(), getRandomKey(), getRandomMilliseconds()),
  );
  bench.addAsync(`GET (${namespaces.length} ns & ${keys.length} key)`, () =>
    c.get(getRandomNs(), getRandomKey(), getRandomMilliseconds()),
  );
  await bench.run();
  bench.print();
  c.close();
}

main().catch(err => console.error(err));
