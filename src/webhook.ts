import { RESTPostAPIChannelMessageJSONBody } from "discord-api-types/v10";
import { fetcher } from "itty-fetcher";

const discord = fetcher({
  base: "https://discord.com/api/v10",
});

export function webhook(id: string, token: string) {
  return (message: RESTPostAPIChannelMessageJSONBody) =>
    discord.post(`/webhooks/${id}/${token}`, message);
}
