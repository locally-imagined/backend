package login

import (
	"backend/gen/auth"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/mysql"
)

type Service struct{}

func (s *Service) Login(ctx context.Context, p *auth.LoginPayload) (string, error) {
	cfg := mysql.Cfg(os.Getenv("INSTANCE")+"-m", os.Getenv("DBUSER"), os.Getenv("DBPASS"))
	cfg.DBName = os.Getenv("DBNAME")
	db, err := mysql.DialCfg(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(fmt.Sprintf("USE %s", cfg.DBName))
	if err != nil {
		log.Fatal(err)
	}
	return "hi", err

}
