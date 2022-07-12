package cmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	ccmd "ntsc.ac.cn/ta-registry/pkg/cmd"
	"ntsc.ac.cn/ta/ntp-server/internal/ntp"
)

var ntpEnvs struct {
	ntpBindAddr string
}

var ntpCmd = &cobra.Command{
	Use:   "ntp",
	Short: "NTP server",
	Run: func(cmd *cobra.Command, args []string) {
		ccmd.InitGlobalVars()
		if err := ccmd.ValidateStringVar(
			&ntpEnvs.ntpBindAddr, "ntp_bind", true); err != nil {
			logrus.WithField("prefix", "ntp").
				Fatalf("validate var failed: %s", err.Error())
		}
		initNTP()
		ccmd.RunWithSysSignal(nil)
	},
}

func initNTP() {
	var err error
	var s *ntp.Server
	if s, err = ntp.NewServer(&ntp.Config{
		Listener: ntpEnvs.ntpBindAddr,
	}); err != nil {
		logrus.WithField("prefix", "ntp").
			Fatalf("create ntp server failed: %s", err.Error())
	}
	errChan := s.Start()
	go func() {
		select {
		case err := <-errChan:
			logrus.WithField("prefix", "ntp").
				Fatalf("start ntp server failed: %s", err.Error())
		default:
		}
	}()
}
func init() {
	cobra.OnInitialize(func() {})
	viper.AutomaticEnv()
	viper.SetEnvPrefix("NTP")
	ntpCmd.Flags().StringVar(&ccmd.GlobalEnvs.LoggerLevel,
		"logger-level", "DEBUG", "logger level")
	ntpCmd.Flags().StringVar(&ntpEnvs.ntpBindAddr,
		"ntp-bind", "0.0.0.0:123", "NTP server bind address")
}

// Execute TSA KMC main
func Execute() {
	if err := ntpCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
