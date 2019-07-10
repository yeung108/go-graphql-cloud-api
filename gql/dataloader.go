package gql

import (
	"context"
	"fmt"

	"log"

	"github.com/graph-gophers/dataloader"
	uuid "github.com/satori/go.uuid"
)

func GetVendorProductsBatchFn(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	handleError := func(err error) []*dataloader.Result {
		var results []*dataloader.Result
		var result dataloader.Result
		result.Error = err
		results = append(results, &result)
		return results
	}
	var vendorIDs []uuid.UUID
	for _, key := range keys {
		k, err := uuid.FromString(key.String())
		if err != nil {
			fmt.Println("vendorIDs key error: %v ", err)
		}
		vendorIDs = append(vendorIDs, k)
	}
	products, err := keys[0].(*ResolverKey).client().resolver().db.GetVendorProducts(vendorIDs)
	if err != nil {
		return handleError(err)
	}

	var results []*dataloader.Result
	result := dataloader.Result{
		Data:  products,
		Error: nil,
	}
	results = append(results, &result)

	log.Printf("[GetVendorProductsBatchFn] batch size: %d", len(results))
	return results
}

func GetVendorStoresBatchFn(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	handleError := func(err error) []*dataloader.Result {
		var results []*dataloader.Result
		var result dataloader.Result
		result.Error = err
		results = append(results, &result)
		return results
	}
	var vendorIDs []uuid.UUID
	for _, key := range keys {
		k, err := uuid.FromString(key.String())
		if err != nil {
			fmt.Println("vendorIDs key error: %v ", err)
		}
		vendorIDs = append(vendorIDs, k)
	}
	stores, err := keys[0].(*ResolverKey).client().resolver().db.GetVendorStores(vendorIDs)
	if err != nil {
		return handleError(err)
	}

	var results []*dataloader.Result
	result := dataloader.Result{
		Data:  stores,
		Error: nil,
	}
	results = append(results, &result)

	log.Printf("[GetVendorStoresBatchFn] batch size: %d", len(results))
	return results
}

func GetVendorsBatchFn(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	handleError := func(err error) []*dataloader.Result {
		var results []*dataloader.Result
		var result dataloader.Result
		result.Error = err
		results = append(results, &result)
		return results
	}
	var vendorIDs []uuid.UUID
	for _, key := range keys {
		k, err := uuid.FromString(key.String())
		if err != nil {
			fmt.Println("vendorIDs key error: %v ", err)
		}
		vendorIDs = append(vendorIDs, k)
	}
	vendors, err := keys[0].(*ResolverKey).client().resolver().db.GetVendors(vendorIDs)
	if err != nil {
		return handleError(err)
	}

	var results []*dataloader.Result
	result := dataloader.Result{
		Data:  vendors,
		Error: nil,
	}
	results = append(results, &result)
	log.Printf("[GetVendorsBatchFn] batch size: %d", len(results))
	return results
}
