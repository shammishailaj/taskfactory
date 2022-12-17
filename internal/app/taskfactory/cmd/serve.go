package cmd

import (
    "taskfactory/internal/app/taskfactory/server/handlers"
    "taskfactory/internal/app/taskfactory/server/middleware"
    "taskfactory/internal/app/taskfactory/server/routes"
    "taskfactory/pkg/utils"
    "flag"
    "fmt"
    "github.com/fsnotify/fsnotify"
    router "github.com/go-chi/chi"
    log "github.com/sirupsen/logrus"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    "net/http"
    "os"
    "os/signal"
    "syscall"
)

func readConfig(env string) {
    if len(env) > 0 {
        env = fmt.Sprintf(".%s", env)
    }

    configFile := flag.String("c", "", "Manually pass the path to the configuration file")
    flag.Parse()

    if *configFile == "" {
        *configFile = fmt.Sprintf("./configs/app%s.yml", env)
    }

    cwd, cwdErr := os.Getwd()
    if cwdErr != nil {
        log.Printf("Error getting current directory. %s", cwdErr.Error())
    }
    log.Infof("current working directory = %s", cwd)

    log.Infof("Checking for existence of configFile %t", u.FileExists(*configFile))

    log.Infof("Loading configuration from file: %s", *configFile)
    viper.SetConfigFile(*configFile)
    viper.SetConfigType("yaml")

    viper.WatchConfig()
    viper.OnConfigChange(func(e fsnotify.Event) {
        log.Infof("Config file changed: %s. Reloading...", e.Name)
    })

    if err := viper.ReadInConfig(); err != nil {
        panic(err)
    }
}

var serveCmd = &cobra.Command{
    Use:   "serve [flags]",
    Short: "Used to start the web services",
    Long: `Can be used to serve the web services`,
    Run: func(cmd *cobra.Command, args []string) {
        env := os.Getenv(utils.KEY_ENV)

        logger.Info(fmt.Sprintf("Starting %s on %s env..", utils.APP_NAME, env))

        readConfig(env)

        host, port := viper.GetString("host"), viper.GetString("port")

        stop := make(chan os.Signal, 1)
        signal.Notify(stop, os.Interrupt, syscall.SIGINT)

        r := routes.GetRoutesList(logger, router.NewRouter(), middleware.NewLog(logger, true))
        r.Add("/ping", "GET", handlers.NewPing(logger, SEMVER).Handler)
        r.Parse()

        httpErr := make(chan error, 1)
        go func() {
            logger.Info(fmt.Sprintf("Started server on %s:%s..", host, port))
            httpErr <- http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), r.Router)
        }()

        select {
            case err := <-httpErr:
            logger.Error(err.Error())
            case <-stop:
            logger.Info("Stopped via signal")
        }

        logger.Info(fmt.Sprintf("Stopping %s..", utils.APP_NAME))
    },
}

func init() {
    rootCmd.AddCommand(serveCmd)

    // Here you will define your flags and configuration settings.

    // Cobra supports Persistent Flags which will work for this command
    // and all subcommands, e.g.:
    // serveCmd.PersistentFlags().String("foo", "", "A help for foo")

    // Cobra supports local flags which will only run when this command
    // is called directly, e.g.:
    //serveCmd.Flags().BoolP("md", "m", false, "Generate Markdwon Documentation")
    //serveCmd.Flags().BoolP("man", "n", false, "Generate Man Page Documentation")
}
