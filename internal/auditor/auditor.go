package auditor

import (
	"context"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type databaseDriver interface {
	Connect(context.Context, string) (*pgx.Conn, error)
	QueryRow(context.Context, string, ...interface{}) *pgx.Row
}

type Auditor struct {
	sugar *zap.SugaredLogger
	db    *pgx.Conn
}

func CreateAuditor(sugar *zap.SugaredLogger) *Auditor {
	sugar.Debugln("Creating Auditor")
	config
	connectionUrl := fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v?sslmode=disable",
		user, password, host, port, name,
	)
	db, err := pgx.Connect(context.Background(), connectionUrl)
	if err != nil {
		sugar.Errorw("Unable to connect to database",
			"error",
			err,
		)
	}
	return &Auditor{db: db}
}

func (a *Auditor) StartAuditor() {
	a.sugar.Debugln("Starting auditor")
}

func (a *Auditor) StopAuditor() {
	a.sugar.Debugln("Stopping auditor")
	a.db.Close(context.Background())
}

func (a *Auditor) AddLog(string) {
	a.sugar.Debugln("Adding audit log")
	var name string
	var weight int64
	err := a.db.QueryRow(context.Background(), "select name, weight from widgets where id=$1", 42).Scan(&name, &weight)
	if err != nil {
		a.sugar.Errorw("Unable to query database",
			"error",
			err,
		)
	}
}
