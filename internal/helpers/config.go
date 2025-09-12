package helpers

import (
	"sync"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

var (
	Addr             = "89.207.255.214:8080"
	DownloadFilesDir = "./"
)

type Config struct {
	Addr               string `env:"AK_GRPC_ADDR,required"`
	DBDSN              string `env:"DB_DSN,required"`
	JWTSecret          string `env:"JWT_SECRET,required" json:"-"`
	FilesDir           string `env:"FILES_DIR"`
	AccessTTL          string `env:"ACCESS_TTL"`
	DownloadFilesDir   string `env:"DOWNLOAD_FILES_DIR"`
	MigrationsFilesDir string `env:"MIGRATIONS_FILES_DIR"`
}

var (
	cfg      Config
	loadOnce sync.Once
	loadErr  error
)

func LoadConfig() (*Config, error) {
	loadOnce.Do(func() {
		_ = godotenv.Load()
		loadErr = env.Parse(&cfg)
	})
	return &cfg, loadErr
}

// (опционально, только для тестов)
// func resetConfigForTests() { loadOnce = sync.Once{}; cfg = Config{} }
