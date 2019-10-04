package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	token_store "github.com/uaa_token_store/token_store"
)

func init() {

	token_store.Token_Store_File = "/tmp/store.json"

}
func TestTokenStore(t *testing.T) {
	var tStore = token_store.TokenStore{}

	tStore.UaaURL = "https://uaaUrl.com"
	tStore.ClientID = "clientId"
	tStore.ClientSecret = "clientSecret"
	tStore.TokenInfo = token_store.TokenInfo{}

	if tStore.UaaURL != "https://uaaUrl.com" {
		t.Error("UAA url mismatch")
	}

	if tStore.ClientID != "clientId" {
		t.Error("Client Id mismatch")
	}
	if tStore.ClientSecret != "clientSecret" {
		t.Error("Client Secret mismatch")
	}

}

func TestReadWriteToFile(t *testing.T) {
	var ts = token_store.TokenStore{}
	ts.ClientID = "id"
	ts.ClientSecret = "secret"
	ts.Name = "name"
	ts.UaaURL = "uaa url"
	ts.TokenInfo = token_store.TokenInfo{}
	//persistToFile(ts)
	byteArray, err := json.Marshal(ts)
	if err != nil {
		fmt.Println("Unable to Marshal to JSON ")
	}
	var jsonString = "{\"uaaURL\":\"uaa url\",\"clientID\":\"id\",\"clientSecret\":\"secret\",\"tokenInfo\":{\"accessToken\":\"\",\"createdTime\":\"0001-01-01T00:00:00Z\",\"expiresAt\":\"0001-01-01T00:00:00Z\"},\"name\":\"name\"}"
	if jsonString != string(byteArray) {
		t.Error("Found this message", string(byteArray))
	}
	var collectionString = "[" + jsonString + "]"

	var tss = make([]token_store.TokenStore, 0)
	json.Unmarshal([]byte(collectionString), &tss)

	if "id" != tss[0].ClientID {
		panic("Error" + tss[0].ClientID)
	}

	tss = token_store.LoadSavedTokenStores()

	if len(tss) != 0 {
		t.Error("Should have been 0")
	}

	tss = token_store.AddToken(ts, tss)
	if len(tss) != 1 {
		t.Error("Should have been 1")
	}

	if !token_store.NameTaken(tss, "name") {
		t.Error("Name should have been taken")
	}

	os.Remove(token_store.Token_Store_File)
	fmt.Println("Token Store FIle " + token_store.Token_Store_File)

	data, err := ioutil.ReadFile("test_input.json")
	if err != nil {
		panic("Error while reading from file " + err.Error())
	}
	var tsStores = make([]token_store.TokenStore, 0)
	err = json.Unmarshal(data, &tsStores)
	if err != nil {
		t.Error(err.Error())
	}

	for index := 0; index < len(tsStores); index++ {
		fmt.Println(tsStores[index])
	}
	if len(tsStores) != 1 {
		t.Error(tsStores)
	}

}
