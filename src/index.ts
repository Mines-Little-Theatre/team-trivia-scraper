import { AnswerData, fetchFreeAnswer } from "./trivia";
import { webhook } from "./webhook";

export default {
  async scheduled(_, env) {
    const post = webhook(env.DISCORD_WEBHOOK_ID, env.DISCORD_WEBHOOK_TOKEN);

    let answer: AnswerData;
    try {
      answer = await fetchFreeAnswer(env.TEAM_TRIVIA_REGION_ID);
    } catch (e) {
      await post({
        content: env.BOT_MESSAGE,
        embeds: [
          {
            description: `Failed to retrieve the free answer: ${String(e)}`,
            color: 0xffcc00,
          },
        ],
      });
      return;
    }

    await post({
      content: env.BOT_MESSAGE,

      embeds: [
        {
          description: `\`\`\`json
${JSON.stringify(answer, undefined, 2)}
\`\`\``,
          color: 0x0069b5,
        },
      ],
    });
  },
} satisfies ExportedHandler<Env>;
