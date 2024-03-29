package cmd

import (
	"context"
	"fmt"
	"github.com/r3volut1oner/go-karbo/config"
	"github.com/r3volut1oner/go-karbo/cryptonote"
	"github.com/r3volut1oner/go-karbo/p2p"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
	"os/signal"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "krbd",
	Short:   "Karbo node daemon.",
	Long:    `Karbo node daemon.`,
	Version: "0.0.1",
	Run:     handleCommand,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.krbd.yml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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

		// Search config in home directory with name ".krbd" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".krbd")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func handleCommand(cmd *cobra.Command, args []string) {
	mainnet := config.MainNet()

	// Initialize blockchain storage
	storage := cryptonote.NewMemoryStorage()

	logrusLogger := logrus.New()
	logrusLogger.Out = os.Stdout
	logrusLogger.Level = logrus.TraceLevel

	bc := cryptonote.NewBlockChain(mainnet, storage, logrusLogger)

	if err := bc.Init(); err != nil {
		panic(fmt.Errorf("failed to init blockchain: %w", err))
	}

	ctx := interruptListener()
	cfg := p2p.HostConfig{
		BindAddr: "127.0.0.1:32447",
		Network:  mainnet,
	}

	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	host := p2p.NewNode(bc, cfg, zapLogger)

	fmt.Println("Server started.")

	if err := host.Run(ctx); err != nil {
		panic(err)
	}

	fmt.Println("Server stopped.")
}

func interruptListener() context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		interruptChannel := make(chan os.Signal, 1)
		signal.Notify(interruptChannel, os.Interrupt)

		select {
		case sig := <-interruptChannel:
			fmt.Printf("Received signal (%s). Shutting down...\n", sig)
		}

		cancel()
	}()

	return ctx
}
