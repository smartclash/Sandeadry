package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/smartclash/Sandeadry/indexer"
	"github.com/smartclash/Sandeadry/parser"
	"github.com/smartclash/Sandeadry/storage"
	"github.com/thatisuday/commando"
	"net/url"
	"strings"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Couldn't load .env file")
	}

	commando.
		SetExecutableName("Sandeadry").
		SetVersion("v0.2.0").
		SetDescription("Scrape sanfoundry MCQs, topics and subjects. Save it in a JSON file or a sqlite database")

	commando.
		Register("scrape").
		SetDescription("Scrape sanfoundry website").
		SetShortDescription("Scrape sanfoundry website").
		AddArgument("link", "Link to the degree you want to scrape subjects, topics and MCQs", "").
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
			link := args["link"].Value
			invokeParser(&link)
		})

	commando.
		Register("save").
		SetDescription("Save all subjects, topics and MCQs scrapped into an sqlite database").
		SetShortDescription("Save scrapped data into sqlite DB").
		AddFlag("database,d", "Custom name for the database", commando.String, "sandeadry").
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
			database := flags["database"].Value
			invokeStorage(database.(string))
		})

	commando.
		Register("index").
		SetDescription("Index all subjects, topics and MCQs scrapped into meilisearch").
		SetShortDescription("Index data into meilisearch").
		AddFlag("database,d", "Custom name for the database", commando.String, "sandeadry").
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
			database := flags["database"].Value
			invokeIndexer(database.(string))
		})

	commando.Parse(nil)
}

func invokeParser(link *string) {
	parse, err := url.Parse(*link)
	if err != nil {
		fmt.Println("Please enter a proper link")
		return
	}

	if !strings.EqualFold(parse.Hostname(), "www.sanfoundry.com") {
		fmt.Println("Enter only sanfoundry links")
		return
	}

	parser.Parser(*link)
}

func invokeStorage(database string) {
	err := storage.Init(database)
	if err != nil {
		fmt.Println("Couldn't save the files into database", err)
		return
	}
}

func invokeIndexer(database string) {
	err := indexer.Init(database)
	if err != nil {
		fmt.Println("Couldn't index files into meilisearch", err)
		return
	}
}
