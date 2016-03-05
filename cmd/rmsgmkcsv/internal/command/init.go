package command

func initCommandPkg() {
	flags = new(TopLevelCommands)
	RootCmd.AddCommand(versionCmd)
	RootCmd.AddCommand(createCmd)
	initCreate()
}

func init() {
	initCommandPkg()
}
