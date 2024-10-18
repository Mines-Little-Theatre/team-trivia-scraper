import slug from "slug";

export async function generateImage(
  { AI, IMAGE_GENERATION_MODEL }: Env,
  answer: string,
): Promise<File> {
  const imageStream = await AI.run(IMAGE_GENERATION_MODEL, {
    prompt: answer,
  });
  const imageBlob = await new Response(imageStream).blob();
  return new File([imageBlob], `${slug(answer)}.png`, { type: "image/png" });
}
