package autorization

import (
	"time"
)

type Team struct {
	ID       uint      `json:"id,omitempty"`
	Name     string    `json:"name,omitempty"`
	CreateAt time.Time `json:"createAt,omitempty"`
	Users    []*User   `json:"users,omitempty"`
}

type User struct {
	ID       uint      `json:"id,omitempty"`
	Name     string    `json:"name,omitempty"`
	Password string    `json:"paasWord,omitempty"`
	Mail     string    `json:"mail,omitempty"`
	Phone    string    `json:"phone,omitempty"`
	Count    int       `json:"count,omitempty"`
	CreateAt time.Time `json:"createAt,omitempty"`
	GroupId  uint      `json:"groudId,omitempty"`
	RoleId   uint      `json:"roleId,omitempty"`
}

type Role struct {
	ID          uint          `json:"id,omitempty"`
	Name        string        `json:"name,omitempty"`
	CreateAt    time.Time     `json:"createAt,omitempty"`
	Permissions []*Permission `json:"permissions,omitempty"`
}

type Permission struct {
	ID       uint      `json:"id,omitempty"`
	Name     string    `json:"name,omitempty"`
	Url      string    `json:"url,omitempty"`
	CreateAt time.Time `json:"createAt,omitempty"`
	RoleId   uint      `json:"roleId,omitempty"`
}
