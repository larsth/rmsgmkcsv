package command

func init() {
	flags = new(TopLevelCommands)
	RootCmd.AddCommand(versionCmd)
	RootCmd.AddCommand(createCmd)
	initCreate()
}
