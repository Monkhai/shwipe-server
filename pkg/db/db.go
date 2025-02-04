package db

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	pool *pgxpool.Pool
	ctx  context.Context
}

const DB_PASSWORD string = "qxk*NQG5gau9key-edu"
const PATH_TO_SUPABASE_CERT string = "./prod-ca-2021.crt"

func NewDB(ctx context.Context) (*DB, error) {
	caCert, err := os.ReadFile(PATH_TO_SUPABASE_CERT)
	if err != nil {
		return nil, err
	}
	caPool := x509.NewCertPool()
	caPool.AppendCertsFromPEM(caCert)

	dsn := fmt.Sprintf("postgresql://postgres.xbmpgyybutdgjmvexqsr:%s@aws-0-us-west-1.pooler.supabase.com:5432/postgres", DB_PASSWORD)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	config.ConnConfig.TLSConfig = &tls.Config{
		RootCAs:            caPool,
		InsecureSkipVerify: false,
	}

	config.MaxConns = 10

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return &DB{pool: pool, ctx: ctx}, nil
}

func (db *DB) Close() {
	if db.pool != nil {
		db.pool.Close()
	}
	log.Println("DB closed")
}

func (db *DB) CreateQuery(queryFile, templateName string, data interface{}) (string, error) {
	sqlBytes, err := SqlFiles.ReadFile(queryFile)
	if err != nil {
		return "", err
	}

	tmpl, err := template.New(templateName).Parse(string(sqlBytes))
	if err != nil {
		return "", err
	}

	var queryBuf bytes.Buffer
	err = tmpl.Execute(&queryBuf, data)
	if err != nil {
		return "", err
	}

	return queryBuf.String(), nil
}

func (db *DB) RunQuery(sql string, args ...interface{}) (pgx.Rows, error) {
	return db.pool.Query(db.ctx, sql, args...)
}

func (db *DB) ExecuteQuery(sql string, args ...interface{}) error {
	_, err := db.pool.Exec(db.ctx, sql, args...)
	return err
}
