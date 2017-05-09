package api

import (
	"context"
	"database/sql"
	"fmt"
)

//Device represents a device in the inventory database
type Device struct {
	InventoryNumber string `json:"inventory_number"`
	SerialNumber    string `json:"serial_number"`
	BagTag          string `json:"bag_tag"`
	Status          string `json:"status"`
	Model           string `json:"model"`
	User            string `json:"user"`
}

//readDeviceBySerialNumber returns the device with the given serial number
func readDeviceBySerialNumber(ctx context.Context, serialNumber string) (*Device, error) {
	tx := ctx.Value(TransactionKey).(*sql.Tx)

	device := &Device{SerialNumber: serialNumber}

	row := tx.QueryRow("SELECT inventory_number, bag_tag, status, model, user FROM devices WHERE serial_number=?;", serialNumber)
	err := row.Scan(
		&(device.InventoryNumber),
		&(device.BagTag),
		&(device.Status),
		&(device.Model),
		&(device.User),
	)

	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, &Error{Description: fmt.Sprintf("Could not query Device(%s)", serialNumber), Err: err}
	}

	return device, nil
}

//ReadDeviceByGoogleID returns the device with the given Google ID
func ReadDeviceByGoogleID(ctx context.Context, svc *ChromebookService, id string) (*Device, error) {
	serial, err := svc.Get(id)
	if err != nil {
		return nil, err
	}
	return readDeviceBySerialNumber(ctx, serial)
}
