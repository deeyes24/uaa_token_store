package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func init() {

	Token_Store_File = "/tmp/store.json"

}
func TestTokenStore(t *testing.T) {
	var tStore = TokenStore{}

	tStore.UaaURL = "https://uaaUrl.com"
	tStore.ClientID = "clientId"
	tStore.ClientSecret = "clientSecret"
	tStore.TokenInfo = TokenInfo{}

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
	var ts = TokenStore{}
	ts.ClientID = "id"
	ts.ClientSecret = "secret"
	ts.Name = "name"
	ts.UaaURL = "uaa url"
	ts.TokenInfo = TokenInfo{}
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

	var tss = make([]TokenStore, 0)
	json.Unmarshal([]byte(collectionString), &tss)

	if "id" != tss[0].ClientID {
		panic("Error" + tss[0].ClientID)
	}

	tss = LoadSavedTokenStores()

	if len(tss) != 0 {
		t.Error("Should have been 0")
	}

	tss = AddToken(ts, tss)
	if len(tss) != 1 {
		t.Error("Should have been 1")
	}

	if !NameTaken(tss, "name") {
		t.Error("Name should have been taken")
	}

	os.Remove(Token_Store_File)
	fmt.Println("Token Store FIle " + Token_Store_File)

	data, err := ioutil.ReadFile("test_input.json")
	if err != nil {
		panic("Error while reading from file " + err.Error())
	}
	var tsStores = make([]TokenStore, 0)
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
