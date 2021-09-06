package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
)

const serverAddress = "http://localhost:" + ServerPort

func TestValidAddress(t *testing.T) {
	addresses := map[string]bool{
		"0x07b1c12be0d62fe548a2b4b025ab7a5ca8def21e":  true,
		"0x8e26e2fc8140280fba3e34bfdca7fc1102c1ae04":  true,
		"0x0000000000000000000000000000000000000000":  true,
		"0x00000000000000000000000000000000000000000": false,
		"00x0000000000000000000000000000000000000000": false,
		"x00000000000000000000000000000000000000000":  false,
		"0x620cd19eae24fb8a02df908bb71b81b6e3aa1cc":   false,
		"0x645c3a387b8633df1d4d71ca4b50d27233bcb8":    false,
		"1x919fa96e88d67499339577fa202345436bcdaf79":  false,
		"009445bd19767f73dcae6f2de90e6cd31192f62589":  false,
		"07b1c12be0d62fe548a2b4b025ab7a5ca8def21e":    false,
		"": false,
	}

	for k, v := range addresses {
		ret := validAddress(k)
		if ret != v {
			t.Fatalf("%s, expected: %t, got: %t\n", k, v, ret)
		}
	}
}

func TestBadUrls(t *testing.T) {
	testBadUrl := func(url string, statusCode int) {
		resp, err := http.Get(url)
		if err != nil {
			t.Fatal(err)
		}

		expectedStatusCode := statusCode
		if resp.StatusCode != expectedStatusCode {
			t.Fatalf("Request failed for %s, got status code: %d, expected %d\n", url, resp.StatusCode, expectedStatusCode)
		}
	}

	testBadUrl(serverAddress, 404)
	testBadUrl(serverAddress+"/", 404)
	testBadUrl(serverAddress+"/2f2jo3i", 404)
	testBadUrl(serverAddress+"/assets/", 404)
	testBadUrl(serverAddress+"/assets/vxcv", 404)
	testBadUrl(serverAddress+"/assets/pools", 404)
	testBadUrl(serverAddress+"/assets/0x00/pools", 400)
	testBadUrl(serverAddress+"/assets/sdf/pools/abc", 404)
}

// Note that the test result may change as user add/remove pools on Uniswap
func TestAssetPool(t *testing.T) {
	assetId := "0xd533a949740bb3306d119cc777fa900ba034cd52" // CRV
	expected := map[string]bool{
		"0x07b1c12be0d62fe548a2b4b025ab7a5ca8def21e": false,
		"0x8e26e2fc8140280fba3e34bfdca7fc1102c1ae04": false,
		"0x4c83a7f819a5c37d64b4c5a2f8238ea082fa1f4e": false,
		"0x620cd19eae24fb8a02df908bb71b81b6e3aa1ccc": false,
		"0x645c3a387b8633df1d4d71ca4b50d27233bcb887": false,
		"0x919fa96e88d67499339577fa202345436bcdaf79": false,
		"0x9445bd19767f73dcae6f2de90e6cd31192f62589": false,
		"0xcbeb7da1ec121fc37dde2bc9010f3a4001e1ebcb": false,
	}

	endpoint := serverAddress + "/assets/" + assetId + "/pools"
	resp, err := http.Get(endpoint)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Request failed with status code: %d\n", resp.StatusCode)
	}

	body, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		t.Fatal(err2)
	}

	var pools []Pool
	err3 := json.Unmarshal(body, &pools)
	if err3 != nil {
		t.Fatal(err3)
	}

	if len(pools) != len(expected) {
		t.Fatalf("Incorrect number of pools. Expected %d, got %d\n", len(expected), len(pools))
	}

	for i := 0; i < len(pools); i++ {
		expected[pools[i].Id] = true
	}

	for k, v := range expected {
		if !v {
			t.Fatal("Missing pool " + k)
		}
	}
}
