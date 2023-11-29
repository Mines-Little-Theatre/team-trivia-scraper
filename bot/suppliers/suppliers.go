package suppliers

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Quantaly/shunt"
	"github.com/bwmarrin/discordgo"
)

// Supplier provides data that the output configuration can reference
type Supplier interface {
	SupplyData(*SupplierContext) error
}

var suppliers map[string]Supplier = make(map[string]Supplier)

func RegisterSupplier(name string, supplier Supplier) {
	_, alreadyExists := suppliers[name]
	if alreadyExists {
		log.Fatalln("multiple suppliers registered with name:", name)
	}
	suppliers[name] = supplier
}

// SupplierContext provides input to and collects output from a Supplier
type SupplierContext struct {
	configPrefix string
	embeds       map[string]*discordgo.MessageEmbed
}

// Config reads a supplier-specific configuration variable
func (c *SupplierContext) Config(key string) string {
	return os.Getenv(c.configPrefix + key)
}

// AddEmbed supplies a named embed that the message can include
func (c *SupplierContext) AddEmbed(name string, embed *discordgo.MessageEmbed) {
	if c.embeds == nil {
		c.embeds = make(map[string]*discordgo.MessageEmbed)
	}
	c.embeds[name] = embed
}

// SupplierResult is the collected data from a number of Suppliers
type SupplierResults struct {
	Embeds map[string]*discordgo.MessageEmbed
}

// RunSuppliers runs the named Suppliers in parallel and collects the results
func RunSuppliers(supplierNames []string) SupplierResults {
	errorCount := 0
	tasks := make(map[string]shunt.Task[SupplierContext], len(supplierNames))
	for _, name := range supplierNames {
		supplier, ok := suppliers[name]
		if !ok {
			log.Println("unknown supplier:", name)
			errorCount++
		} else if _, alreadyUsed := tasks[name]; alreadyUsed {
			log.Println("duplicate supplier:", name)
			errorCount++
		} else {
			tasks[name] = shunt.Do(func() (SupplierContext, error) {
				context := SupplierContext{
					configPrefix: "TRIVIA_CONFIG_" + strings.ToUpper(name) + "_",
					embeds:       make(map[string]*discordgo.MessageEmbed),
				}
				err := supplier.SupplyData(&context)
				return context, err
			})
		}
	}

	results := SupplierResults{
		Embeds: make(map[string]*discordgo.MessageEmbed),
	}
	for name, task := range tasks {
		context, err := task.JoinWithoutPanicking()
		if err != nil {
			log.Println("supplier", name, ":", err)
			errorCount++
		} else {
			for embedName, embed := range context.embeds {
				results.Embeds[name+":"+embedName] = embed
			}
		}
	}

	if errorCount > 0 {
		results.Embeds["errors"] = &discordgo.MessageEmbed{
			Description: fmt.Sprintf("%v supplier(s) returned errors. Check the logs for details.", errorCount),
			Color:       0xffcc00,
		}
	}

	return results
}
