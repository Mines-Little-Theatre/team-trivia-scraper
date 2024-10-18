import { generateImage } from "./image-generation";
import { AnswerData, fetchFreeAnswer, freeAnswerURL } from "./trivia";
import { webhook } from "./webhook";

const embedColor = 0x0069b5;
const errorColor = 0xffcc00;

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
            color: errorColor,
          },
        ],
      });
      return;
    }

    let image: File;
    try {
      image = await generateImage(env, answer.answer);
    } catch (e) {
      await post({
        content: env.BOT_MESSAGE,
        embeds: [
          {
            title: answer.title,
            url: freeAnswerURL,
            fields: [
              {
                name: answer.date,
                value: answer.answer,
              },
            ],
            footer: { text: `Failed to generate image: ${String(e)}` },
            color: embedColor,
          },
        ],
      });
      return;
    }

    await post(
      {
        content: env.BOT_MESSAGE,
        embeds: [
          {
            title: answer.title,
            url: freeAnswerURL,
            fields: [
              {
                name: answer.date,
                value: answer.answer,
              },
            ],
            image: {
              url: `attachment://${image.name}`,
              width: 1024,
              height: 1024,
            },
            footer: {
              text: `Image is AI-generated (${env.IMAGE_GENERATION_MODEL})`,
            },
            color: embedColor,
          },
        ],
      },
      image,
    );
  },
} satisfies ExportedHandler<Env>;
