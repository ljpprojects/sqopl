const { X509Certificate } = require("node:crypto");
const fs = require("node:fs/promises");

(async () => {
  const lines = (
    await fs.readFile("/Volumes/monster/gibberish/hello.gibberish")
  )
    .toString()
    .split("\n")
    .map((l) => l.trimStart())
    .filter((l) => l[0] !== "#" && l.length > 0);

  let depth = 0;

  const Types = {
    Nothing: -1,
    Numeric: 0,
    Reference: 1,
  };

  let li = 0;

  while (li < lines.length) {
    let line = lines[li];

    let ci = 0;

    while (ci < line.length) {
      if (line.slice(ci, ci + 6) === "repeat") {
        console.log("REPEAT");

        ci += 7;

        console.log(line.slice(ci));

        const n = (() => {
          let n = "";

          while (
            ci < line.length &&
            [...line[ci]].every((c) => new Set([..."0123456789"]).has(c))
          ) {
            n += line[ci++];
          }

          return Number(n);
        })();

        depth++;

        line = lines[++li];
        ci = 0;

        while (
          li < lines.length &&
          line.slice(0, depth) === ">".repeat(depth)
        ) {
          console.log(line);
          line = lines[++li];
          ci = 0;
        }

        ci += 6;
      }

      ci++;
    }

    li++;
  }
})();
