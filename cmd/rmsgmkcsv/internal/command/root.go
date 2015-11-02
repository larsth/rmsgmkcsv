package command

import "github.com/spf13/cobra"

var (
	RootCmd = &cobra.Command{
		Use: CommandName,
		Short: "rmsgmkcsv is a CLI tool that create decimal " +
			"bearings/azimuths to microprocessor H-bridge actuator integers.",
		Long: "rmsgmkcsv either creates predictable linear values, or calls " +
			"an external program/script interpreter to make one translation " +
			"from one decimal bearing/azimuth to one whole integer value, " +
			"and does that for a specified range of decimal degrees, with a " +
			"step degree.",
	}
)
