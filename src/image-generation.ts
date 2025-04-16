import { Base64 } from "js-base64";
import slug from "slug";

export async function generateImage(env: Env, answer: string): Promise<File> {
  const { image } = await env.AI.run(env.IMAGE_GENERATION_MODEL, {
    prompt: answer,
  });
  if (image) {
    return new File([Base64.toUint8Array(image)], `${slug(answer)}.jpg`, {
      type: "image/jpg",
    });
  } else {
    throw new Error("no image returned by API");
  }
}
