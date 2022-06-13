import * as util from "util";
import { exec } from "shelljs";

export async function x(...args) {
  let res;
  try {
    res = await util.promisify(require("child_process").exec)(...args);
  } catch (e) {
    res = e;
  }
  if (res.stdout) console.log(res.stdout);
  if (res.stderr) console.error(res.stderr);
  return res;
}

export async function xjs(cmd, opts = {}) {
  return new Promise(function (resolve, reject) {
    // Execute the command, reject if we exit non-zero (i.e. error)
    exec(cmd, opts, function (code, stdout, stderr) {
      if (code != 0) return reject(new Error(stderr));
      return resolve(stdout);
    });
  });
}
export async function runMod(mod: string) {
  delete require.cache[require.resolve(mod)];
  try {
    require(mod);
  } catch (e) {
    console.log(e);
  }
}
