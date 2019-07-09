package postgres

import (
	"database/sql"
	"time"

	uuid "github.com/satori/go.uuid"
)

// Payment Method shape
type PaymentMethod struct {
	ID            uuid.UUID `db:"id" json:"id"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
	MongoID       string    `db:"mongo_id" json:"mongo_id"`
	Code          string    `db:"code" json:"code"`
	Name          string    `db:"name" json:"name"`
	Module        string    `db:"module" json:"module"`
	ModuleChannel string    `db:"module_channel" json:"module_channel"`
	OrderIndex    string    `db:"order_index" json:"order_index"`
}

type Vendor struct {
	ID          uuid.UUID `db:"id" json:"id,omitempty"`
	CreatedAt   time.Time `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at,omitempty"`
	MongoID     string    `db:"mongo_id" json:"mongo_id,omitempty"`
	Name        string    `db:"name" json:"name,omitempty"`
	Description string    `db:"description" json:"description,omitempty"`
	Products    []Product `json:"products,omitempty"`
}

// Product shape
type Product struct {
	ID               uuid.UUID      `db:"id" json:"id,omitempty"`
	CreatedAt        time.Time      `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt        time.Time      `db:"updated_at" json:"updated_at,omitempty"`
	MongoID          sql.NullString `db:"mongo_id" json:"mongo_id,omitempty"`
	Photo            sql.NullString `db:"photo" json:"photo,omitempty"`
	Code             sql.NullString `db:"code" json:"code,omitempty"`
	IsVirtualProduct bool           `db:"is_virtual_product" json:"is_virtual_product,omitempty"`
	Barcode          sql.NullString `db:"barcode" json:"barcode,omitempty"`
	Descriptions     LanguageJson   `db:"descriptions" json:"descriptions,omitempty"`
	BrandNames       LanguageJson   `db:"brand_names" json:"brand_names,omitempty"`
	Names            LanguageJson   `db:"names" json:"names,omitempty"`
	OptionalData     LanguageJson   `db:"optional_data" json:"optional_data,omitempty"`
	VendorID         uuid.NullUUID  `db:"vendor_id" json:"vendor_id,omitempty"`
	SupplierID       uuid.NullUUID  `db:"supplier_id" json:"supplier_id,omitempty"`
}

type LanguageJson struct {
	En string `db:"en" json:"en,omitempty"`
	Zh string `db:"zh" json:"zh,omitempty"`
}
