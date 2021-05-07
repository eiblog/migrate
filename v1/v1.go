// Package v1 provides ...
package v1

import "time"

// db info
const (
	DB                 = "eiblog"
	COLLECTION_ACCOUNT = "account"
	COLLECTION_ARTICLE = "article"
	COUNTER_SERIE      = "serie"
	COUNTER_ARTICLE    = "article"
)

// Account 账户
type Account struct {
	// 账户信息
	Username   string    // 账户名
	Password   string    // 账户密码
	Token      string    // 二次验证token
	Email      string    // 账户
	PhoneN     string    // 手机号
	Address    string    // 住址
	CreateTime time.Time // 创建时间
	LoginTime  time.Time // 最后登录时间
	LogoutTime time.Time // 登出时间
	LoginIP    string    // 最后登录ip
	// 博客信息
	BlogName  string // 博客名
	SubTitle  string // SubTitle
	BeiAn     string // 备案号
	BTitle    string // 底部title
	Copyright string // 版权声明
	SeriesSay string // 专题，倒序
	Series    []struct {
		ID         int32     // 自增id
		Name       string    // 名称unique
		Slug       string    // 缩略名
		Desc       string    // 专题描述
		CreateTime time.Time // 创建时间
	} // 专题
	ArchivesSay string // 归档描述
	Archives    []struct {
		Time time.Time
	} // 归档
}

// Article 文章
type Article struct {
	ID         int32     // 自增id
	Author     string    // 作者名
	Title      string    // 标题
	Slug       string    // 文章名: how-to-get-girlfriend
	Count      int       // 评论数量
	Content    string    // markdown文档
	SerieID    int32     // 归属专题
	Tags       []string  // tagname
	IsDraft    bool      // 是否是草稿
	CreateTime time.Time // 创建时间
	UpdateTime time.Time // 更新时间
	DeleteTime time.Time // 开始删除时间
}

// EiBlog v1数据
type EiBlog struct {
	Account  Account
	Articles []Article
}
