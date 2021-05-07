// Package store provides ...
package store

import (
	v1 "github.com/eiblog/migrate/v1"
	v2 "github.com/eiblog/migrate/v2"

	"github.com/lib/pq"
)

type migrate interface {
	// LoadEiBlog 读取数据
	LoadEiBlog(from Store) (interface{}, error)
	// StoreEiBlog 保存数据
	StoreEiBlog(to Store, blog interface{}) error
}

// Store 存储信息
type Store struct {
	Version string `yaml:"version"`
	Driver  string `yaml:"driver"`
	Source  string `yaml:"source"`
}

// Migrate driver->store
var Migrate = make(map[string]migrate)

func init() {
	Migrate["mongodb"] = mongodb{}
	Migrate["mysql"] = rdbms{}
	Migrate["postgres"] = rdbms{}
	Migrate["sqlite"] = rdbms{}
	Migrate["sqlserver"] = rdbms{}
	Migrate["clickhouse"] = rdbms{}
}

// V1ToV2 升级数据v1->v2
func V1ToV2(v1 *v1.EiBlog) *v2.EiBlog {
	blogger := v2.Blogger{
		BlogName:    v1.Account.BlogName,
		SubTitle:    v1.Account.SubTitle,
		BeiAn:       v1.Account.BeiAn,
		BTitle:      v1.Account.BTitle,
		Copyright:   v1.Account.Copyright,
		SeriesSay:   v1.Account.SeriesSay,
		ArchivesSay: v1.Account.ArchivesSay,
	}
	acct := v2.Account{
		Username: v1.Account.Username,
		Password: v1.Account.Password,
		Email:    v1.Account.Email,
		PhoneN:   v1.Account.PhoneN,
		Address:  v1.Account.Address,

		LogoutAt:  v1.Account.LogoutTime,
		LoginIP:   v1.Account.LoginIP,
		LoginUA:   "", // no
		LoginAt:   v1.Account.LoginTime,
		CreatedAt: v1.Account.CreateTime,
	}
	var series []v2.Serie
	for _, v := range v1.Account.Series {
		serie := v2.Serie{
			ID:        int(v.ID),
			Slug:      v.Slug,
			Name:      v.Name,
			Desc:      v.Desc,
			CreatedAt: v.CreateTime,
		}
		series = append(series, serie)
	}
	var articles []v2.Article
	for _, v := range v1.Articles {
		article := v2.Article{
			ID:      int(v.ID),
			Author:  v.Author,
			Slug:    v.Slug,
			Title:   v.Title,
			Count:   v.Count,
			Content: v.Content,
			SerieID: int(v.SerieID),
			Tags:    pq.StringArray(v.Tags),
			IsDraft: v.IsDraft,

			DeletedAt: v.DeleteTime,
			UpdatedAt: v.UpdateTime,
			CreatedAt: v.CreateTime,
		}
		articles = append(articles, article)
	}
	return &v2.EiBlog{
		Blogger:  blogger,
		Account:  acct,
		Series:   series,
		Articles: articles,
	}
}
