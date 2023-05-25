package core

import (
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var Cfg *Config

type Config struct {
	Season        int
	Debug         bool
	Api configApi
	Database configDatabase
	AWS configAWS
}

type configApi struct {
	Host, Key string
}
type configDatabase struct {
	User, Password, Location, Name string
	Port                           int
}
type configAWS struct {
	Region string
	AccessKeyId string `mapstructure:"access-key-id"`
	SecretAccessKey string `mapstructure:"secret-access-key"`
	BucketName string `mapstructure:"bucket-name"`
}

func SetupViper() {
	viper.SetDefault("use-config-file", true)
	viper.SetEnvPrefix("nf")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()
}

// this is intentionally kept separate from SetupViper() as configFile will eventually be passed in from a Cobra command
func SetupConfigFile(configFile string) {
	if viper.GetBool("use-config-file") {
		viper.SetConfigFile(configFile)
		if err := viper.ReadInConfig(); err == nil {
			Log.Infof("Loaded config file: %s", viper.ConfigFileUsed())
		} else {
			Log.Fatalf("Could not load config file: %s!", configFile)
		}
	} else {
		Log.Info("Config file NOT being used...requiring NF_ENVIRONMENT_VARIABLES")
		bindViperEnvVars()
	}

	if viper.GetBool("debug") {
		Log.SetLevel(logrus.DebugLevel)
		Log.Debug("Debug logging enabled!")
	}

	Cfg = &Config{}
	if err := viper.Unmarshal(Cfg); err != nil {
		Log.Fatalf("Error decoding Config struct: %v", err)
	} else {
		Log.Infof("Config successfully initialized")
	}
}

// viper needs a little help with these nested variables...
func bindViperEnvVars() {
	viper.BindEnv("api.host")
	viper.BindEnv("api.key")
	viper.BindEnv("database.user")
	viper.BindEnv("database.password")
	viper.BindEnv("database.location")
	viper.BindEnv("database.port")
	viper.BindEnv("database.name")
	viper.BindEnv("aws.region")
	viper.BindEnv("aws.access-key-id")
	viper.BindEnv("aws.secret-access-key")
	viper.BindEnv("aws.bucket-name")
}