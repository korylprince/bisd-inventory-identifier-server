package api

import (
	"context"
	"fmt"
	"io/ioutil"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/admin/directory/v1"
)

//ChromebookService is a service to request information about Chromebooks from Google
type ChromebookService struct {
	service *admin.ChromeosdevicesService
}

//Get returns the serial number for the chromebook with the given id
func (s *ChromebookService) Get(id string) (serialNumber string, err error) {
	dev, err := s.service.Get("my_customer", id).Fields("serialNumber").Do()
	if err != nil {
		return "", &Error{Description: fmt.Sprintf("Could not GET Chromebook(%s)", id), Err: err}
	}
	return dev.SerialNumber, nil
}

//NewChromebookService returns a ChromebookService
func NewChromebookService(configPath, impersonateUser string) (*ChromebookService, error) {
	buf, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, &Error{Description: fmt.Sprintf("Could not read OAuth JSON(%s)", configPath), Err: err}
	}

	config, err := google.JWTConfigFromJSON(buf, admin.AdminDirectoryDeviceChromeosReadonlyScope)
	if err != nil {
		return nil, &Error{Description: fmt.Sprintf("Could not read OAuth JSON(%s)", configPath), Err: err}
	}
	config.Subject = impersonateUser

	adminSvc, err := admin.New(config.Client(context.Background()))
	if err != nil {
		return nil, &Error{Description: fmt.Sprintf("Could not read create admin service"), Err: err}
	}

	return &ChromebookService{service: admin.NewChromeosdevicesService(adminSvc)}, nil
}
