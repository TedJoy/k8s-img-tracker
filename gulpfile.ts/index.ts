import { xjs } from "@/utils";

export async function hello() {
  await xjs("echo Hello World");
}
