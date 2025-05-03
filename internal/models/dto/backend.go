package dto

import (
	"github.com/asaskevich/govalidator"
	"strings"
)

type AddBackendRequest struct {
	ServerURL string `json:"server_url" valid:"required,url"`
}

type DeleteBackendRequest struct {
	ServerURL string `json:"server_url" valid:"required,url"`
}

func (a *AddBackendRequest) Validate() error {
	a.ServerURL = strings.TrimRight(a.ServerURL, "/")
	_, err := govalidator.ValidateStruct(a)
	return err
}

func (d *DeleteBackendRequest) Validate() error {
	d.ServerURL = strings.TrimRight(d.ServerURL, "/")
	_, err := govalidator.ValidateStruct(d)
	return err
}
