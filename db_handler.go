package main

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type DBHandler struct {
	dbPath         string
	contextTimeout time.Duration
}

func (d *DBHandler) Init(ctx context.Context) error {
	db, err := sql.Open("sqlite3", d.dbPath)
	defer db.Close()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(ctx, d.contextTimeout)
	defer cancel()
	_, err = db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS "proxy" (
		"id" INTEGER NOT NULL,
		"name" TEXT NOT NULL,
		"create_time" INTEGER NOT NULL,
		"s"	TEXT NOT NULL,
		"p"	INTEGER NOT NULL,
		"b"	TEXT NOT NULL,
		"l"	INTEGER NOT NULL,
		"k"	TEXT NOT NULL,
		"m"	TEXT NOT NULL,
		"o"	TEXT NOT NULL,
		"op" TEXT NOT NULL,
		"oo" TEXT NOT NULL,
		"oop" TEXT NOT NULL,
		"t" INTEGER NOT NULL,
		"f" TEXT NOT NULL,
		"status" INTEGER NOT NULL,
		PRIMARY KEY("id" AUTOINCREMENT)
	)`)
	return err
}

func (d *DBHandler) InsertProxy(ctx context.Context, p *Proxy) (int64, error) {
	db, err := sql.Open("sqlite3", d.dbPath)
	defer db.Close()
	if err != nil {
		return 0, err
	}
	ctx, cancel := context.WithTimeout(ctx, d.contextTimeout)
	defer cancel()
	var id int64
	db.QueryRowContext(
		ctx,
		"insert into proxy(name,create_time,s,p,b,l,k,m,o,op,oo,oop,t,f,status) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?) RETURNING id",
		p.Name,
		p.CreateTime,
		p.S,
		p.P,
		p.B,
		p.L,
		p.K,
		p.M,
		p.O,
		p.Op,
		p.Oo,
		p.Oop,
		p.T,
		p.F,
		p.Status,
	).Scan(&id)
	return id, err
}

func (d *DBHandler) UpdateProxy(ctx context.Context, p *Proxy) error {
	db, err := sql.Open("sqlite3", d.dbPath)
	defer db.Close()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(ctx, d.contextTimeout)
	defer cancel()
	_, err = db.ExecContext(
		ctx,
		"update proxy set name = ?, s = ?, p = ?, b = ?, l = ?, k = ?, m = ?, o = ?, op = ?, oo = ?, oop = ?, t = ?, f = ?, status = ? where id = ?",
		p.Name,
		p.S,
		p.P,
		p.B,
		p.L,
		p.K,
		p.M,
		p.O,
		p.Op,
		p.Oo,
		p.Oop,
		p.T,
		p.F,
		p.Status,
		p.Id,
	)
	return err
}

func (d *DBHandler) DeleteProxy(ctx context.Context, id int64) error {
	db, err := sql.Open("sqlite3", d.dbPath)
	defer db.Close()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(ctx, d.contextTimeout)
	defer cancel()
	_, err = db.ExecContext(
		ctx,
		"delete from proxy where id = ?",
		id,
	)
	return err
}

func (d *DBHandler) GetProxies(ctx context.Context, onlyStatusOne bool) ([]Proxy, error) {
	db, err := sql.Open("sqlite3", d.dbPath)
	defer db.Close()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, d.contextTimeout)
	defer cancel()
	query := "select id,name,create_time,s,p,b,l,k,m,o,op,oo,oop,t,f,status from proxy"
	if onlyStatusOne {
		query = query + " where status = 1"
	}
	rows, err := db.QueryContext(
		ctx,
		query,
	)
	if err != nil {
		return nil, err
	}
	ret := make([]Proxy, 0)
	for rows.Next() {
		p := Proxy{}
		err = rows.Scan(&p.Id, &p.Name, &p.CreateTime, &p.S, &p.P, &p.B, &p.L, &p.K, &p.M, &p.O, &p.Op, &p.Oo, &p.Oop, &p.T, &p.F, &p.Status)
		if err != nil {
			return nil, err
		}
		ret = append(ret, p)
	}
	return ret, nil
}

func (d *DBHandler) GetProxy(ctx context.Context, id int64) (*Proxy, error) {
	db, err := sql.Open("sqlite3", d.dbPath)
	defer db.Close()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, d.contextTimeout)
	defer cancel()
	row := db.QueryRowContext(
		ctx,
		"select id,name,create_time,s,p,b,l,k,m,o,op,oo,oop,t,f,status from proxy where id = ?",
		id,
	)
	p := Proxy{}
	err = row.Scan(&p.Id, &p.Name, &p.CreateTime, &p.S, &p.P, &p.B, &p.L, &p.K, &p.M, &p.O, &p.Op, &p.Oo, &p.Oop, &p.T, &p.F, &p.Status)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func NewDBHandler(dbName string, contextTimeout time.Duration) *DBHandler {
	return &DBHandler{
		dbPath:         dbName,
		contextTimeout: contextTimeout,
	}
}
