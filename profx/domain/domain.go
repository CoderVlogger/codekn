package domain

import (
	"crypto/sha256"
	"fmt"
	"time"
)

type RuleType int
type ResourceKind string

const (
	IncludeRuleType RuleType = iota + 1
	ExcludeRuleType

	ArticleResourceKind ResourceKind = "article"
	NewsResourceKind                 = "news"
)

// Model is a base domain model.
type Model struct {
	Created time.Time `bson:"created"   json:"created"   db:"created"`
}

// SysLog is exported.
type SysLog struct {
	Model
	Message string `bson:"message"   json:"message"   db:"message"`
}

// Source is exported.
type Source struct {
	Model
	Name string       `bson:"name"   json:"name"   db:"name"`
	URL  string       `bson:"url"    json:"url"    db:"url"`
	Kind ResourceKind `bson:"kind"   json:"kind"   db:"kind"`
}

// Rule is exported.
type Rule struct {
	Model
	Type   RuleType `bson:"type"     json:"type"     db:"type"`
	Source string   `bson:"source"   json:"source"   db:"source"`
	Rule   string   `bson:"rule"     json:"rule"     db:"rule"`
}

// Link is exported.
type Link struct {
	Model
	Hash string `bson:"hash"   json:"hash"   db:"hash"`
	URL  string `bson:"url"    json:"url"    db:"url"`

	Source  string `bson:"source"     json:"source"     db:"source"`
	FromURL string `bson:"from_url"   json:"from_url"   db:"from_url"`
}

// Resource is exported.
type Resource struct {
	Hash string `bson:"hash"    json:"hash"    db:"hash"`
	URL  string `bson:"url"     json:"url"     db:"url"`

	Kind    ResourceKind `bson:"kind"       json:"kind"       db:"kind"`
	Source  string       `bson:"source"     json:"source"     db:"source"`
	FromURL string       `bson:"from_url"   json:"from_url"   db:"from_url"`

	Title       string `bson:"title"   json:"title"   db:"title"`
	Description string `bson:"desc"    json:"desc"    db:"desc"`
}

// Hash returns the hash value of the given URL.
func Hash(url string) string {
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(url)))
	return hash
}
