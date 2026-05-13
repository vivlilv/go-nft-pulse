package reducer

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/vivlilv/go_nft_trader/internal/domain"
)

func TestMain(m *testing.M) {
	os.Chdir("../..") // now CWD = project root
	if err := InitConfig(); err != nil {
		log.Fatalf("Config load for ./internal/reducer/ %v", err)
	}
	fmt.Printf("%s", string(domain.MyWallet.Address))
	os.Exit(m.Run())
}

func InitConfig() error {
	godotenv.Load()
	return envconfig.Process("", domain.MyWallet)
}
