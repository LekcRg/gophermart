package config

import (
	"fmt"
	"log"
	"os"

	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	flag "github.com/spf13/pflag"
)

// don't forget to change the getFlags func
type Config struct {
	IsDev          bool   `koanf:"IS_DEV"`
	DBPass         string `koanf:"POSTGRES_PASSWORD"`
	DBUser         string `koanf:"POSTGRES_USER"`
	DBHost         string `koanf:"POSTGRES_HOST"`
	DBName         string `koanf:"POSTGRES_DB"`
	DBPort         int    `koanf:"POSTGRES_PORT"`
	DBURI          string `koanf:"DATABASE_URI"`
	Address        string `koanf:"RUN_ADDRESS"`
	AccrualAddress string `koanf:"ACCRUAL_SYSTEM_ADDRESS"`
	JWTSecret      string `koanf:"JWT_SECRET_KEY"`
}

var k = koanf.New(".")

func getFlags() {
	fl := flag.NewFlagSet("config", flag.ContinueOnError)
	fl.Usage = func() {
		fmt.Println(fl.FlagUsages())
		os.Exit(0)
	}

	fl.BoolP("IS_DEV", "x", false, "Development mode")
	fl.StringP("RUN_ADDRESS", "a", "", "Address to run gophermart")
	fl.StringP("DATABASE_URI", "d", "", "Database URI")
	fl.StringP("ACCRUAL_SYSTEM_ADDRESS", "r", "", "Accrual system address")
	fl.StringP("JWT_SECRET_KEY", "j", "", "JWT secret key")

	fl.Parse(os.Args[1:])
	if err := k.Load(posflag.Provider(fl, ".", k), nil); err != nil {
		log.Fatalf("error loading config: %v", err)
	}
}

func getDefaults() {
	defaultCfg := Config{
		DBPass:         "postgres",
		DBUser:         "postgres",
		DBHost:         "localhost",
		DBPort:         5432,
		DBName:         "gophermart",
		DBURI:          "",
		Address:        "localhost:8080",
		AccrualAddress: "localhost:8000",
	}
	k.Load(structs.Provider(defaultCfg, "koanf"), nil)
}

func getDotEnv() {
	envfile := file.Provider("../../.env")
	if err := k.Load(envfile, dotenv.Parser()); err != nil {
		fmt.Println(err)
	}
}

func getEnvVars() {
	if err := k.Load(env.Provider("", ".", nil), nil); err != nil {
		fmt.Println(err)
	}
}

func Get() Config {
	getDefaults()
	getDotEnv()
	getFlags()
	getEnvVars()

	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		fmt.Println(err)
	}

	if cfg.DBURI == "" && cfg.DBPass != "" &&
		cfg.DBUser != "" && cfg.DBHost != "" &&
		cfg.DBName != "" && cfg.DBPort != 0 {
		cfg.DBURI = fmt.Sprintf("postgresql://%s:%s@%s:%d/%s",
			cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName)
	}

	return cfg
}
