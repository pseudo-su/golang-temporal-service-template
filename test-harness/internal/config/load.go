package config

import (
	"embed"
	"fmt"
	"maps"
	"os"
	"strings"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

//go:embed *
var envFilesFS embed.FS

func dotenvFiles(envStr string) []string {
	result := []string{}

	baseFiles := []string{
		".env",
		".env.local",
	}
	result = append(result, baseFiles...)

	// When an environment suffixed with `.ext` is selected we want to include both
	// the unsuffixed and suffixed dotenv files (eg .env, .env.sit, .env.sit.ext)
	if es, ok := strings.CutSuffix(envStr, ".ext"); ok {
		envFiles := []string{
			fmt.Sprintf(".env.%s", es),
			fmt.Sprintf(".env.%s.local", es),
		}
		result = append(result, envFiles...)
	}

	envFiles := []string{
		fmt.Sprintf(".env.%s", envStr),
		fmt.Sprintf(".env.%s.local", envStr),
	}
	result = append(result, envFiles...)

	return result
}

func loadEnvFromFiles(envStr string) (map[string]string, error) {
	envMap := map[string]string{}
	fileNames := dotenvFiles(envStr)

	for _, fileName := range fileNames {
		file, err := envFilesFS.Open(fileName)
		if err != nil {
			// Ignore any files if they don't exist
			continue
		}
		envVars, err := godotenv.Parse(file)
		if err != nil {
			return nil, err
		}

		maps.Copy(envMap, envVars)
	}

	return envMap, nil
}

func loadEnvFromShell() map[string]string {
	envMap := map[string]string{}
	for _, e := range os.Environ() {
		if i := strings.Index(e, "="); i >= 0 {
			envMap[e[:i]] = e[i+1:]
		}
	}
	return envMap
}

func loadEnv(envStr string) (map[string]string, error) {
	envMap := map[string]string{}
	fromFiles, err := loadEnvFromFiles(envStr)
	if err != nil {
		return nil, err
	}

	maps.Copy(envMap, fromFiles)
	maps.Copy(envMap, loadEnvFromShell())

	return envMap, nil
}

func IsValidEnv(envStr string) error {
	fileName := fmt.Sprintf(".env.%s", envStr)
	_, err := envFilesFS.Open(fileName)
	if err != nil {
		return fmt.Errorf("env %s not valid, config file not found: %w", envStr, err)
	}
	return nil
}

func LoadEnvConfig() (*TestsuiteEnvConfig, error) {
	envConfig := TestsuiteEnvConfig{}

	envMap, err := loadEnv(os.Getenv("ENV"))
	if err != nil {
		return nil, err
	}

	if err = env.ParseWithOptions(&envConfig, env.Options{
		Environment: envMap,
	}); err != nil {
		return nil, err
	}

	return &envConfig, nil
}
