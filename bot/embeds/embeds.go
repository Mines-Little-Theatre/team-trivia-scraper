package embeds

import "github.com/bwmarrin/discordgo"

// EmbedProvider asynchronously provides a MessageEmbed for the output.
//
// The function should send a maximum of one embed on the provided channel and close the channel when it is finished.
type EmbedProvider func(chan<- *discordgo.MessageEmbed)

// CollectEmbeds calls the given provider functions and returns a slice of their results in the order that they are passed.
//
// If a provider closes its channel without sending an embed, it is skipped in the output.
//
// The provider functions are run concurrently.
func CollectEmbeds(providers ...EmbedProvider) []*discordgo.MessageEmbed {
	resultChannels := make([]<-chan *discordgo.MessageEmbed, len(providers))
	for i := range providers {
		channel := make(chan *discordgo.MessageEmbed, 1)
		go providers[i](channel)
		resultChannels[i] = channel
	}

	results := make([]*discordgo.MessageEmbed, 0, len(providers))
	for _, resultChannel := range resultChannels {
		result, ok := <-resultChannel
		if ok {
			results = append(results, result)
		}
	}

	return results
}
