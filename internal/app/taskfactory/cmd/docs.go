package cmd

import (
    log "github.com/sirupsen/logrus"
    "github.com/spf13/cobra"
    "github.com/spf13/cobra/doc"
)

var docsCmd = &cobra.Command{
    Use:   "docs [flags] /path/to/directory/to/output/docs/",
    Short: "Generates markdown documentation",
    Long: `Can be used to generate documentation in markdown`,
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) > 0 {
            err := doc.GenMarkdownTree(rootCmd, args[0])
            if err != nil {
                log.Fatal(err)
            }
        } else {
            log.Fatalln("Must provide path to output documentation")
        }
    },
}

func init() {
    rootCmd.AddCommand(docsCmd)

    // Here you will define your flags and configuration settings.

    // Cobra supports Persistent Flags which will work for this command
    // and all subcommands, e.g.:
    // docsCmd.PersistentFlags().String("foo", "", "A help for foo")

    // Cobra supports local flags which will only run when this command
    // is called directly, e.g.:
    docsCmd.Flags().BoolP("md", "m", false, "Generate Markdwon Documentation")
    docsCmd.Flags().BoolP("man", "n", false, "Generate Man Page Documentation")
}