import { RESTPostAPIChannelMessageJSONBody } from "discord-api-types/v10";
import { fetcher } from "itty-fetcher";

const discord = fetcher({
  base: "https://discord.com/api/v10",
});

export function webhook(id: string, token: string) {
  return async (
    message: RESTPostAPIChannelMessageJSONBody,
    ...files: readonly File[]
  ): Promise<void> => {
    if (files.length === 0) {
      await discord.post(`/webhooks/${id}/${token}?wait=true`, message);
    } else {
      const formData = new FormData();
      formData.append("payload_json", JSON.stringify(message));
      files.forEach((file, i) => {
        formData.append(`files[${String(i)}]`, file);
      });
      await discord.post(`/webhooks/${id}/${token}?wait=true`, formData);
    }
  };
}
