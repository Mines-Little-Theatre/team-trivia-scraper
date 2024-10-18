import JPEG_ENC_WASM from "@jsquash/jpeg/codec/enc/mozjpeg_enc.wasm";
import encodeJPEG, { init as initJPEGEncode } from "@jsquash/jpeg/encode";
import PNG_DEC_WASM from "@jsquash/png/codec/pkg/squoosh_png_bg.wasm";
import decodePNG, { init as initPNGDecode } from "@jsquash/png/decode";
import slug from "slug";

export async function generateImage(env: Env, answer: string): Promise<File> {
  const [, , pngBuf] = await Promise.all([
    initJPEGEncode(JPEG_ENC_WASM),
    initPNGDecode(PNG_DEC_WASM),
    generatePNGBuf(env, answer),
  ]);
  const imageData = await decodePNG(pngBuf);
  const jpegBuf = await encodeJPEG(imageData);
  console.log(
    `JPEG encoding reduced image size from ${String(pngBuf.byteLength)} to ${String(jpegBuf.byteLength)} (${percentChange(
      pngBuf.byteLength,
      jpegBuf.byteLength,
    )})`,
  );
  return new File([jpegBuf], `${slug(answer)}.jpg`, { type: "image/jpeg" });
}

async function generatePNGBuf(env: Env, answer: string): Promise<ArrayBuffer> {
  const imageStream = await env.AI.run(env.IMAGE_GENERATION_MODEL, {
    prompt: answer,
  });
  return await new Response(imageStream).arrayBuffer();
}

const percentFormat = new Intl.NumberFormat("en-US", {
  style: "percent",
  signDisplay: "always",
  maximumSignificantDigits: 3,
});

function percentChange(from: number, to: number): string {
  if (from === to) {
    return "no change";
  } else {
    return percentFormat.format((to - from) / from);
  }
}
