package pokecache

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestAdd(t *testing.T) {
	fmt.Println("=== Add() ===")
	tests := []struct {
		key  string
		val  []byte
	}{
		{key: "https://example.com", val: []byte("testdata")},
		{key: "https://example.com/path", val: []byte("moretestdata")},
	}

	for i, tc := range tests {
		cache := NewCache(5 * time.Second)
		cache.Add(tc.key, tc.val)

		// Using your 'got' system
		got, ok := cache.Cache[tc.key]
		if !ok {
			t.Fatalf("test %d: expected key %s to exist, but it was not found", i+1, tc.key)
		}

		if !reflect.DeepEqual(tc.val, got.val) {
			t.Fatalf("test %d: expected: %v, got: %v", i+1, tc.val, got.val)
		} else {
			fmt.Printf("Passed test: %d\n", i+1)
		}
	}
	fmt.Print("=============\n\n")
}

func TestGet(t *testing.T) {
	fmt.Println("=== Get() ===")
	tests := []struct {
		keyToSet string
		valToSet []byte
		keyToGet string
		wantVal  []byte
		wantOk   bool
	}{
		{
			keyToSet: "key1",
			valToSet: []byte("val1"),
			keyToGet: "key1",
			wantVal:  []byte("val1"),
			wantOk:   true,
		},
		{
			keyToSet: "key1",
			valToSet: []byte("val1"),
			keyToGet: "nonexistent",
			wantVal:  []byte(nil),
			wantOk:   false,
		},
	}

	for i, tc := range tests {
		cache := NewCache(5 * time.Second)
		cache.Add(tc.keyToSet, tc.valToSet)

		gotVal, gotOk := cache.Get(tc.keyToGet)

		if gotOk != tc.wantOk {
			t.Fatalf("test %d: expected ok: %v, got ok: %v\n", i+1, tc.wantOk, gotOk)
		}

		if !reflect.DeepEqual(tc.wantVal, gotVal) {
			t.Fatalf("test %d: expected: %v, got: %v\n", i+1, tc.wantVal, gotVal)
		} else {
			fmt.Printf("Passed test: %d\n", i+1)
		}
	}

	fmt.Print("=============\n\n")
}