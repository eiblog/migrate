// Package v2 provides ...
package v2

import (
	"time"

	"github.com/lib/pq"
)

// mongodb info
const (
	MongoDBName       = "eiblog"
	CollectionAccount = "account"
	CollectionArticle = "article"
	CollectionBlogger = "blogger"
	CollectionCounter = "counter"
	CollectionSerie   = "serie"

	CounterNameSerie   = "serie"
	CounterNameArticle = "article"
)

// Account 博客账户
type Account struct {
	Username string `gorm:"column:username;primaryKey" bson:"username"` // 用户名
	Password string `gorm:"column:password;not null" bson:"password"`   // 密码
	Email    string `gorm:"column:email;not null" bson:"email"`         // 邮件地址
	PhoneN   string `gorm:"column:phone_n;not null" bson:"phone_n"`     // 手机号
	Address  string `gorm:"column:address;not null" bson:"address"`     // 地址信息

	LogoutAt  time.Time `gorm:"column:logout_at;not null" bson:"logout_at"`        // 登出时间
	LoginIP   string    `gorm:"column:login_ip;not null" bson:"login_ip"`          // 最近登录IP
	LoginUA   string    `gorm:"column:login_ua;not null" bson:"login_ua"`          // 最近登录IP
	LoginAt   time.Time `gorm:"column:login_at;default:now()" bson:"login_at"`     // 最近登录时间
	CreatedAt time.Time `gorm:"column:created_at;default:now()" bson:"created_at"` // 创建时间
}

// Blogger 博客信息
type Blogger struct {
	BlogName  string `gorm:"column:blog_name;not null" bson:"blog_name"` // 博客名
	SubTitle  string `gorm:"column:sub_title;not null" bson:"sub_title"` // 子标题
	BeiAn     string `gorm:"column:bei_an;not null" bson:"bei_an"`       // 备案号
	BTitle    string `gorm:"column:b_title;not null" bson:"b_title"`     // 底部title
	Copyright string `gorm:"column:copyright;not null" bson:"copyright"` // 版权声明

	SeriesSay   string `gorm:"column:series_say;not null" bson:"series_say"`     // 专题说明
	ArchivesSay string `gorm:"column:archives_say;not null" bson:"archives_say"` // 归档说明
}

// Article 文章
type Article struct {
	ID      int            `gorm:"column:id;primaryKey" bson:"id"`               // ID, store自行控制
	Author  string         `gorm:"column:author;not null" bson:"author"`         // 作者名
	Slug    string         `gorm:"column:slug;not null;uniqueIndex" bson:"slug"` // 文章缩略名
	Title   string         `gorm:"column:title;not null" bson:"title"`           // 标题
	Count   int            `gorm:"column:count;not null" bson:"count"`           // 评论数量
	Content string         `gorm:"column:content;not null" bson:"content"`       // markdown内容
	SerieID int            `gorm:"column:serie_id;not null" bson:"serie_id"`     // 专题ID
	Tags    pq.StringArray `gorm:"column:tags;type:text[]" bson:"tags"`          // tags
	IsDraft bool           `gorm:"column:is_draft;not null" bson:"is_draft"`     // 是否是草稿

	DeletedAt time.Time `gorm:"column:deleted_at;not null" bson:"deleted_at"`      // 删除时间
	UpdatedAt time.Time `gorm:"column:updated_at;default:now()" bson:"updated_at"` // 更新时间
	CreatedAt time.Time `gorm:"column:created_at;default:now()" bson:"created_at"` // 创建时间
}

// Serie 专题
type Serie struct {
	ID        int       `gorm:"column:id;primaryKey" bson:"id"`                    // 自增ID
	Slug      string    `gorm:"column:slug;not null;uniqueIndex" bson:"slug"`      // 缩略名
	Name      string    `gorm:"column:name;not null" bson:"name"`                  // 专题名
	Desc      string    `gorm:"column:desc;not null" bson:"desc"`                  // 专题描述
	CreatedAt time.Time `gorm:"column:created_at;default:now()" bson:"created_at"` // 创建时间
}

// EiBlog 数据
type EiBlog struct {
	Blogger  Blogger
	Account  Account
	Articles []Article
	Series   []Serie
}
