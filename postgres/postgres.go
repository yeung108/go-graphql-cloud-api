package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"

	// postgres driver

	"github.com/lib/pq"
	_ "github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
)

// Db is our database struct used for interacting with the database
type Db struct {
	*sql.DB
}

// New makes a new database using the connection string and
// returns it, otherwise returns the error
func New(connString string) (*Db, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	// Check that our connection is good
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &Db{db}, nil
}

// ConnString returns a connection string based on the parameters it's given
// This would normally also contain the password, however we're not using one
func ConnString(host string, port int, user string, password string, dbName string, sslMode string) string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbName, sslMode,
	)
}

func (d *Db) GetVendorProducts(vendorIDs []uuid.UUID) ([]Product, error) {
	// Create Vendor struct for holding each row's data
	var r Product
	// Create slice of Users for our response
	products := []Product{}
	// Make query with our stmt, passing in phoneDeviceID argument
	//rows, err := d.Query("SELECT vendor.*, array_to_json(array_agg(row_to_json(product.*))) AS products FROM vendor JOIN product ON product.vendor_id = vendor.id GROUP BY vendor.id WHERE vendor.id IN $1", vendorIDs)
	rows, err := d.Query(`SELECT product.* FROM product JOIN vendor ON product.vendor_id = vendor.id WHERE vendor.id = ANY($1)`, pq.Array(vendorIDs))

	if err != nil {
		return products, fmt.Errorf("GetVendorProducts Query Err: %+v", err)
	}

	// Copy the columns from row into the values pointed at by r (Product)
	for rows.Next() {
		var descriptions string
		var brandNames string
		var names string
		var optionalData string
		err = rows.Scan(
			&r.ID,
			&r.CreatedAt,
			&r.UpdatedAt,
			&r.MongoID,
			&r.Photo,
			&r.Code,
			&r.IsVirtualProduct,
			&r.Barcode,
			&descriptions,
			&brandNames,
			&names,
			&optionalData,
			&r.VendorID,
			&r.SupplierID,
		)
		descriptionErr := json.Unmarshal([]byte(descriptions), &r.Descriptions)
		brandNamesErr := json.Unmarshal([]byte(brandNames), &r.BrandNames)
		namesErr := json.Unmarshal([]byte(names), &r.Names)
		optionalDataErr := json.Unmarshal([]byte(optionalData), &r.OptionalData)
		// productErr := json.Unmarshal([]byte(products), &r.Products)
		if err != nil {
			return products, fmt.Errorf("Error scanning rows: %+v", err)
		} else if descriptionErr != nil {
			return products, fmt.Errorf("Error scanning descriptions: %+v", descriptionErr)
		} else if brandNamesErr != nil {
			return products, fmt.Errorf("Error scanning brandNames: %+v", brandNamesErr)
		} else if namesErr != nil {
			return products, fmt.Errorf("Error scanning names: %+v", namesErr)
		} else if optionalDataErr != nil {
			return products, fmt.Errorf("Error scanning optionalData: %+v", optionalDataErr)
		}
		products = append(products, r)
	}
	return products, nil
}
func (d *Db) GetVendors(vendorIDs []uuid.UUID) ([]Vendor, error) {
	// Create Vendor struct for holding each row's data
	var r Vendor
	// Create slice of Users for our response
	vendors := []Vendor{}
	// Make query with our stmt, passing in phoneDeviceID argument
	//rows, err := d.Query("SELECT vendor.*, array_to_json(array_agg(row_to_json(product.*))) AS products FROM vendor JOIN product ON product.vendor_id = vendor.id GROUP BY vendor.id WHERE vendor.id IN $1", vendorIDs)
	rows, err := d.Query(`SELECT * FROM vendor WHERE vendor.id = ANY($1)`, pq.Array(vendorIDs))

	if err != nil {
		return vendors, fmt.Errorf("GetVendors Query Err: %+v", err)
	}

	// Copy the columns from row into the values pointed at by r (Vendors)
	for rows.Next() {
		err = rows.Scan(
			&r.ID,
			&r.CreatedAt,
			&r.UpdatedAt,
			&r.MongoID,
			&r.Name,
			&r.Description,
		)
		if err != nil {
			return vendors, fmt.Errorf("Error scanning rows: %+v", err)
		}
		vendors = append(vendors, r)
	}
	return vendors, nil
}
