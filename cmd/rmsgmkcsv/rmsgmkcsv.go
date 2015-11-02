package main

import (
	"log"
	"os"

	"github.com/larsth/rmsgmkcsv/cmd/rmsgmkcsv/internal/command"
	"github.com/spf13/cobra"
)

func main() {
	var cmd *cobra.Command = command.RootCmd
	//log.Logger settings
	log.SetFlags(log.Ldate | log.Lshortfile | log.LUTC)
	log.SetOutput(os.Stderr)
	log.SetPrefix(command.CommandName)

	if cmd == nil {
		log.Fatalln("cmd er <nil>")
	}
	if err := cmd.Execute(); err != nil {
		log.Println(err.Error())
		os.Exit(-2)
	}
	os.Exit(0)
}
