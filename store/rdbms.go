// Package store provides ...
package store

import (
	"fmt"

	v2 "github.com/eiblog/migrate/v2"

	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type rdbms struct{}

// LoadEiBlog 读取数据
func (db rdbms) LoadEiBlog(from Store) (interface{}, error) {
	client, err := db.getClient(from)
	if err != nil {
		return nil, err
	}
	switch from.Version {
	case "v2":
		// blogger
		blogger := v2.Blogger{}
		err = client.First(&blogger).Error
		if err != nil {
			return nil, err
		}
		// account
		acct := v2.Account{}
		err = client.First(&acct).Error
		if err != nil {
			return nil, err
		}
		// articles
		var articles []v2.Article
		err = client.Find(&articles).Error
		if err != nil {
			return nil, err
		}
		// serie
		var series []v2.Serie
		err = client.Find(&series).Error
		if err != nil {
			return nil, err
		}
		blog := &v2.EiBlog{
			Blogger:  blogger,
			Account:  acct,
			Articles: articles,
			Series:   series,
		}
		return blog, nil
	}
	return nil, fmt.Errorf("unsupported version: %s", from.Version)
}

// StoreEiBlog 保存数据
func (db rdbms) StoreEiBlog(to Store, blog interface{}) error {
	client, err := db.getClient(to)
	if err != nil {
		return err
	}
	switch data := blog.(type) {
	case *v2.EiBlog:
		err = client.AutoMigrate(v2.Blogger{},
			v2.Account{},
			v2.Article{},
			v2.Serie{})
		if err != nil {
			return err
		}
		// blogger
		err = client.Create(&data.Blogger).Error
		if err != nil {
			return err
		}
		// account
		err = client.Create(&data.Account).Error
		if err != nil {
			return err
		}
		// articles
		err = client.CreateInBatches(data.Articles, len(data.Articles)).Error
		if err != nil {
			return err
		}
		// series
		err = client.CreateInBatches(data.Series, len(data.Series)).Error
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported data version: %T", blog)
	}
	return nil
}

func (db rdbms) getClient(info Store) (*gorm.DB, error) {
	switch info.Driver {
	case "mysql":
		// https://github.com/go-sql-driver/mysql
		return gorm.Open(mysql.Open(info.Source), &gorm.Config{})
	case "postgres":
		// https://github.com/go-gorm/postgres
		return gorm.Open(postgres.Open(info.Source), &gorm.Config{})
	case "sqlite":
		// github.com/mattn/go-sqlite3
		return gorm.Open(sqlite.Open(info.Source), &gorm.Config{})
	case "sqlserver":
		// github.com/denisenkom/go-mssqldb
		return gorm.Open(sqlserver.Open(info.Source), &gorm.Config{})
	case "clickhouse":
		return gorm.Open(clickhouse.Open(info.Source), &gorm.Config{})
	}
	return nil, fmt.Errorf("unsupported driver: %s", info.Driver)
}
