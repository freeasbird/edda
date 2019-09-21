package dao

import (
	"context"
	"strings"
	"time"
)

// 数据库
type DB interface {
	Init(ctx context.Context, option ...Option) (err error)
	Insert(table string, instance interface{}) (id string, err error)
	Find(table string, filter interface{}, queryF interface{}, skip int64, limit int64, sort int) error
	FindOne(coll string, filter interface{}, result interface{}) (err error)
	Delete(coll string, filter interface{}) (err error)
	Update(coll string, filter, update interface{}) (err error)
	Count(coll string, filter interface{}) (num int64, err error)
	Aggregation(coll string, pipe interface{}, cursorF interface{}) (err error) // 聚合查询
}

func NewDB(driver string) (db DB) {
	switch strings.ToLower(driver) {
	case "mongodb", "mongo":
		return new(MongoCli)
	}
	return nil
}

type Options struct {
	Host      string
	Port      string
	Username  string
	Password  string
	Timeout   time.Duration
	Database  string
	CollIndex map[string]string
}

type Option func(opts *Options)

func WithHost(host string) Option {
	return func(opts *Options) {
		opts.Host = host
	}
}

func WithPort(port string) Option {
	return func(opts *Options) {
		opts.Port = port
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(opts *Options) {
		opts.Timeout = timeout
	}
}

func WithUsername(user string) Option {
	return func(opts *Options) {
		opts.Username = user
	}
}

func WithPwd(pwd string) Option {
	return func(opts *Options) {
		opts.Password = pwd
	}
}

func WithDB(db string) Option {
	return func(opts *Options) {
		opts.Database = db
	}
}

func WithCollIndex(ci map[string]string) Option {
	return func(opts *Options) {
		opts.CollIndex = ci
	}
}
