package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Db     Db
	Server Server
	Token  Token
}

type Db struct {
	ConnString string
}

type Server struct {
	Port int
}

type Token struct {
	Secret string
	TTL    time.Duration
}

func LoadConfig() (Config, error) {
	err := loadEnvFromFile(".env")
	if err != nil {
		return Config{}, fmt.Errorf("Ошибка при загрузке файла .env:", err)
	}

	servPortEnv := os.Getenv("server_port")
	servPort, err := strconv.Atoi(servPortEnv)
	if err != nil {
		return Config{}, err
	}
	if servPort < 1 || servPort > 65535 {
		return Config{}, fmt.Errorf("Порт сервера должен находиться в диапазоне от 1 до 65535")
	}

	dbConnStringEnv := os.Getenv("conn_string")

	tokenSecretEnv := os.Getenv("token_secret")
	tokenTTLEnv := os.Getenv("token_ttl")
	tokenTTL, err := time.ParseDuration(tokenTTLEnv)
	if err != nil {
		return Config{}, err
	}

	return Config{
		Server: Server{
			Port: servPort,
		},
		Db: Db{
			ConnString: dbConnStringEnv,
		},
		Token: Token{
			TTL:    tokenTTL,
			Secret: tokenSecretEnv,
		},
	}, nil
}

func loadEnvFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "=")
		if len(parts) == 2 {
			key := parts[0]
			value := parts[1]
			os.Setenv(key, value)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
