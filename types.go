package main

type ErrorResponse struct {
	Error error `json:"error"`
}

type AppConfig struct {
	DatabaseUrl string `json:"database_url"`
}

type Meta struct {
	Package string `json:"package"`
	ID      uint64 `json:"id"`
	Version uint64 `json:"version"`
	Extras  JSON   `db:"type:json" json:"extras,omitempty"`
}

type ModuleRequest struct {
	Data []byte `json:"data"`
	Meta Meta   `json:"meta"`
}

type Module struct {
	Data     []byte `json:"data"`
	Version  uint64 `json:"version"`
	Package  string `json:"package"`
	Extras   JSON   `sql:"type:json" json:"extras,omitempty"`
	IsActive bool   `json:"is_active" db:"is_active"`
	ID       uint64 `json:"id"`
}

type DeviceModuleRequest struct {
	ID      uint64 `json:"id"`
	Version uint64 `json:"version"`
}

type DeviceRequest struct {
	DeviceID         string                `json:"deviceId"`
	InstalledModules []DeviceModuleRequest `json:"installedModules"`
}
