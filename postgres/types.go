package postgres

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/lib/pq"

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

// Store shape
type Store struct {
	ID                    uuid.UUID      `db:"id" json:"id,omitempty"`
	CreatedAt             time.Time      `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt             time.Time      `db:"updated_at" json:"updated_at,omitempty"`
	MongoID               sql.NullString `db:"mongo_id" json:"mongo_id,omitempty"`
	Code                  sql.NullString `db:"code" json:"code,omitempty"`
	Name                  sql.NullString `db:"name" json:"name,omitempty"`
	Model                 sql.NullString `db:"model" json:"model,omitempty"`
	Address               sql.NullString `db:"address" json:"address,omitempty"`
	LastOnlineAt          pq.NullTime    `db:"last_online_at" json:"last_online_at,omitempty"`
	LastGet               pq.NullTime    `db:"last_get" json:"last_get,omitempty"`
	LastSync              pq.NullTime    `db:"last_sync" json:"last_sync,omitempty"`
	LastRefill            pq.NullTime    `db:"last_refill" json:"last_refill,omitempty"`
	LastReset             pq.NullTime    `db:"last_reset" json:"last_reset,omitempty"`
	UnsubmittedOrderCount sql.NullInt64  `db:"unsubmitted_order_count" json:"unsubmitted_order_count,omitempty"`
	VendorID              uuid.NullUUID  `db:"vendor_id" json:"vendor_id,omitempty"`
}

type LanguageJson struct {
	En string `db:"en" json:"en,omitempty"`
	Zh string `db:"zh" json:"zh,omitempty"`
}

// Make the Attrs struct implement the driver.Valuer interface. This method
// simply returns the JSON-encoded representation of the struct.
func (lj LanguageJson) Value() (driver.Value, error) {
	return json.Marshal(lj)
}

// Make the Attrs struct implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (lj *LanguageJson) Scan(value interface{}) error {
	newLJ := LanguageJson{}
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	err := json.Unmarshal(b, &newLJ)
	if err != nil {
		return err
	}
	// need to assign since if the field does not exist, it will inherit into the next scan with the same type
	*lj = newLJ
	return nil
}
