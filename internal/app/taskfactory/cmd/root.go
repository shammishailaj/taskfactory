/*
    Copyright © 2022  <>

    Licensed under the HLT License, Version 0.0.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

    

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/
package cmd

import (
    "taskfactory/pkg/utils"
    "fmt"
    homedir "github.com/mitchellh/go-homedir"
    log "github.com/sirupsen/logrus"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    "os"
    "taskfactory/pkg/schemas"
)

var cfgFile string

var (
    rootCmd = &cobra.Command{
        Use:   "taskfactory [sub-command] [flags] [args]",
        Short: "The  CLI",
        Long:  `A scheduler to be used as a cron replacement`,
        // Uncomment the following line if your bare application
        // has an action associated with it:
        // Run: func(cmd *cobra.Command, args []string) { },
    }
    logger     = log.New()
    u          *utils.Utils
    SEMVER     *schemas.SemanticVersion
    GitBranch  string
    GitState   string
    GitSummary string
    BuildDate  string
    Version    string
    GitCommit  string
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
    if err := rootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
    }
}

func init() {
    logger = log.New()
    u = utils.NewUtils(logger)
    cobra.OnInitialize(initConfig)

    // Here you will define your flags and configuration settings.
    // Cobra supports persistent flags, which, if defined here,
    // will be global for your application.

    rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.taskfactory.yaml)")

    // Cobra also supports local flags, which will only run
    // when this action is called directly.
    rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

    SEMVER = new(schemas.SemanticVersion)
    SEMVER.GitBranch = GitBranch
    SEMVER.GitState = GitState
    SEMVER.GitSummary = GitSummary
    SEMVER.BuildDate = BuildDate
    SEMVER.Version = Version
    SEMVER.GitCommit = GitCommit
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
    if cfgFile != "" {
        // Use config file from the flag.
        viper.SetConfigFile(cfgFile)
    } else {
        // Find home directory.
        home, err := homedir.Dir()
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }

        // Search config in home directory with name ".taskfactory (without extension).
        viper.AddConfigPath(home)
        viper.SetConfigName(".taskfactory")
    }

    viper.AutomaticEnv() // read in environment variables that match

    // If a config file is found, read it in.
    if err := viper.ReadInConfig(); err == nil {
        fmt.Println("Using config file:", viper.ConfigFileUsed())
    }
}
