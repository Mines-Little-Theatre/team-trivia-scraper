import { fetcher } from "itty-fetcher";
import { fetchFreeAnswer } from "./trivia";

const discord = fetcher({
  base: "https://discord.com/api/v10",
});

export default {
  async scheduled(_, env) {
    // const imageStream = await env.AI.run(env.IMAGE_GENERATION_MODEL, {
    //   prompt: "cyberpunk cat",
    // });
    // const imageBlob = await new Response(imageStream).blob();
    // const formData = new FormData();
    // formData.append("files[0]", imageBlob, "cybercat.png");
    // await discord.post(
    //   `/webhooks/${env.DISCORD_WEBHOOK_ID}/${env.DISCORD_WEBHOOK_TOKEN}`,
    //   formData,
    // );
    console.log("yeet");
    await discord.post(
      `/webhooks/${env.DISCORD_WEBHOOK_ID}/${env.DISCORD_WEBHOOK_TOKEN}`,
      {
        content: `\`\`\`json
${JSON.stringify(await fetchFreeAnswer(env.TEAM_TRIVIA_REGION_ID), undefined, 2)}
\`\`\``,
      },
    );
  },
} satisfies ExportedHandler<Env>;
