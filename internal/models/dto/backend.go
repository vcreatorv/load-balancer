package dto

import "github.com/asaskevich/govalidator"

type AddBackendRequest struct {
	ServerURL string `json:"server_url" valid:"required,url"`
}

type DeleteBackendRequest struct {
	ServerURL string `json:"server_url" valid:"required,url"`
}

func (a *AddBackendRequest) Validate() error {
	_, err := govalidator.ValidateStruct(a)
	return err
}

func (d *DeleteBackendRequest) Validate() error {
	_, err := govalidator.ValidateStruct(d)
	return err
}
