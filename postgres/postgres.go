package postgres

import (
	"database/sql"
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
		err = rows.Scan(
			&r.ID,
			&r.CreatedAt,
			&r.UpdatedAt,
			&r.MongoID,
			&r.Photo,
			&r.Code,
			&r.IsVirtualProduct,
			&r.Barcode,
			&r.Descriptions,
			&r.BrandNames,
			&r.Names,
			&r.OptionalData,
			&r.VendorID,
			&r.SupplierID,
		)
		if err != nil {
			return products, fmt.Errorf("Error scanning rows: %+v", err)
		}
		products = append(products, r)
	}
	return products, nil
}

func (d *Db) GetVendorStores(vendorIDs []uuid.UUID) ([]Store, error) {
	// Create Store struct for holding each row's data
	var r Store
	// Create slice of Stores for our response
	stores := []Store{}
	// Make query with our stmt, passing in phoneDeviceID argument
	//rows, err := d.Query("SELECT vendor.*, array_to_json(array_agg(row_to_json(product.*))) AS products FROM vendor JOIN product ON product.vendor_id = vendor.id GROUP BY vendor.id WHERE vendor.id IN $1", vendorIDs)
	rows, err := d.Query(`SELECT store.* FROM store JOIN vendor ON store.vendor_id = vendor.id WHERE vendor.id = ANY($1)`, pq.Array(vendorIDs))

	if err != nil {
		return stores, fmt.Errorf("GetVendorStores Query Err: %+v", err)
	}

	// Copy the columns from row into the values pointed at by r (Store)
	for rows.Next() {
		err = rows.Scan(
			&r.ID,
			&r.CreatedAt,
			&r.UpdatedAt,
			&r.MongoID,
			&r.Code,
			&r.Name,
			&r.Model,
			&r.Address,
			&r.LastOnlineAt,
			&r.LastGet,
			&r.LastSync,
			&r.LastRefill,
			&r.LastReset,
			&r.UnsubmittedOrderCount,
			&r.VendorID,
		)
		if err != nil {
			return stores, fmt.Errorf("Error scanning rows: %+v", err)
		}
		stores = append(stores, r)
	}
	return stores, nil
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

func (d *Db) EditVendors(vendorID uuid.UUID) error {
	// Make query with our stmt, passing in phoneDeviceID argument
	//rows, err := d.Query("SELECT vendor.*, array_to_json(array_agg(row_to_json(product.*))) AS products FROM vendor JOIN product ON product.vendor_id = vendor.id GROUP BY vendor.id WHERE vendor.id IN $1", vendorIDs)
	rows, err := d.Exec(`UPDATE vendor SET description = $1 WHERE id = $2 `, "testssssss", vendorID)
	fmt.Println(rows.LastInsertId)

	if err != nil {
		return fmt.Errorf("GetVendors Query Err: %+v", err)
	}

	return nil
}
