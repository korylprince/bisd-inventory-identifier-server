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

//ReadDeviceBySerialNumber returns the device with the given serial number
func ReadDeviceBySerialNumber(ctx context.Context, tx *sql.Tx, serialNumber string) (*Device, error) {
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
		return nil, fmt.Errorf("Could not query Device(%s): %v", serialNumber, err)
	}

	return device, nil
}
