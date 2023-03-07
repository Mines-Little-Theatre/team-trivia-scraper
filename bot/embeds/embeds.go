package embeds

import (
	"log"

	"github.com/Quantaly/shunt"
	"github.com/bwmarrin/discordgo"
)

// EmbedProvider provides a MessageEmbed for the output.
type EmbedProvider func() (*discordgo.MessageEmbed, error)

// CollectEmbeds calls the given provider functions in parallel and returns a slice of their results in the order that they are passed.
func CollectEmbeds(providers ...EmbedProvider) []*discordgo.MessageEmbed {
	tasks := make([]shunt.Task[*discordgo.MessageEmbed], 0, len(providers))
	for _, f := range providers {
		tasks = append(tasks, shunt.Do(f))
	}

	results := make([]*discordgo.MessageEmbed, 0, len(providers))
	for _, task := range tasks {
		result, err := task.Join()
		if err != nil {
			log.Println(err)
		} else if result != nil {
			results = append(results, result)
		}
	}

	return results
}
