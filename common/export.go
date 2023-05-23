package common

import (
	"html/template"
	"time"
)

type IFlow interface {
	Id() int64
	Accessed()
	Config() FlowConfig
	LastAccessedAt() time.Time
	ObjectIds() []string
	IdentifierString() string
	Title() string
	IdentifierTemplate() *template.Template
	Load(ac AppContexter) error
	Unload(ac AppContexter) error
	Update(ac AppContexter, title, identifier string, objIds []string, config FlowConfig) error
	Save(ac AppContexter) error
}

type IObject interface {
	Id() int64
	Object() *Json
	Update(ac AppContexter, m *Json) error
	IsEnabled() bool
	Enable(ac AppContexter) error
	Disable(ac AppContexter) error
	Save(ac AppContexter) error
}

type IApi interface {
	Id() int64
	Method() string
	Route() string
	ResponseBody() string
	ResponseHeaders() JSONSimpleStrDict
	Status() int
	IsEnabled() bool
	Enable(ac AppContexter) error
	Disable(ac AppContexter) error
	Save(ac AppContexter) error
	Update(ac AppContexter, route, method string, status int, respHeaders JSONSimpleStrDict, respBody string) error
	ResponseTemplate() *template.Template
}
