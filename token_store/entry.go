package token_store

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type TokenStore struct {
	UaaURL       string    `json:"uaaURL"`
	ClientID     string    `json:"clientID"`
	ClientSecret string    `json:"clientSecret"`
	TokenInfo    TokenInfo `json:"tokenInfo,omitempty"`
	Name         string    `json:"name"`
}
type TokenInfo struct {
	AccessToken string    `json:"accessToken"`
	CreatedTime time.Time `json:"createdTime"`
	ExpiresAt   time.Time `json:"expiresAt"`
}
type EnvInfo struct {
	httpProxy string
}

var Token_Store_File = os.Getenv("HOME") + "/.token_stores.json"

func EntryPoint() {

	var savedTokenStores = LoadSavedTokenStores()

	addPtr := flag.Bool("add", false, "Add token to the TokenStore")
	addFromFilePtr := flag.String("add-from-file", "", `Add to the TokenStore from file path. Ex :
		[
			{
				"name" : "dummyName",
				"uaaURL": "dummy Url",
				"clientId": "clientid ",
				"clientSecret":"ClientSecret"
			}
		]
	`)
	flag.Parse()
	reader := bufio.NewReader(os.Stdin)
	if *addPtr {

		fmt.Printf("Enter UAA URL:  ")
		uaaUrl, uaaReadError := reader.ReadString('\n')
		uaaUrl = strings.TrimSuffix(uaaUrl, "\n")
		uaaUrl = strings.TrimSuffix(uaaUrl, "/")
		if uaaReadError != nil || len(uaaUrl) == 0 {
			panic("What did you enter ? I was expecting a non-empty target UAA url")

		}
		fmt.Printf("Enter the clientId:  ")
		cid, clientIdReadError := reader.ReadString('\n')
		cid = strings.TrimSuffix(cid, "\n")
		if clientIdReadError != nil || len(cid) == 0 {
			panic("What did you enter ? I  was expecting a non-empty clientID for the above URL")
		}
		fmt.Printf("Enter the clientSecret: ")
		csecret, clientSecretReadError := reader.ReadString('\n')
		csecret = strings.TrimSuffix(csecret, "\n")
		if clientSecretReadError != nil || len(csecret) == 0 {
			panic("What did you enter? I was expecting a non-empty clientSecret for the above clientId")
		}
		fmt.Printf("Enter a unique name for the token store:  ")
		name, nameReadError := reader.ReadString('\n')
		name = strings.TrimSuffix(name, "\n")
		if nameReadError != nil || len(name) == 0 {
			panic("What did you enter ? I was expecting a non-empty name for the above TokenStore")
		}
		var isNameTaken = NameTaken(savedTokenStores, name)
		if isNameTaken {
			fmt.Printf("Name " + name + " is already taken. Try another unique name\n")
			fmt.Printf("Enter a unique name for the token store ")
			name, nameReadError = reader.ReadString('\n')
			if nameReadError != nil || len(name) == 0 {
				panic("What did you enter ? I was expecting a non-empty name for the above TokenStore")
			}

		}
		var ts = TokenStore{
			UaaURL:       uaaUrl,
			ClientID:     cid,
			ClientSecret: csecret,
			TokenInfo:    getTokenInfo(getTokenResponse(uaaUrl, cid, csecret)),
			Name:         name,
		}

		savedTokenStores = AddToken(ts, savedTokenStores)

		fmt.Println("Access Token is " + ts.TokenInfo.AccessToken)
		return

	}

	if len(*addFromFilePtr) != 0 {
		data, err := ioutil.ReadFile(*addFromFilePtr)
		if err != nil {
			panic("Error while reading from file " + err.Error())
		}
		var tsStores = make([]TokenStore, 0)
		json.Unmarshal(data, &tsStores)

		for index := 0; index < len(tsStores); index++ {
			var curInput = tsStores[index]
			if len(curInput.ClientID) == 0 {
				panic("clientID is missing ")
			}
			if len(curInput.ClientSecret) == 0 {
				panic("clientSecret is missing")
			}
			if len(curInput.Name) == 0 {
				panic("name is missing")
			}
			if len(curInput.UaaURL) == 0 {
				panic("uaaURL is missing")
			}
			var ts = TokenStore{
				UaaURL:       curInput.UaaURL,
				ClientID:     curInput.ClientID,
				ClientSecret: curInput.ClientSecret,
				TokenInfo:    getTokenInfo(getTokenResponse(curInput.UaaURL, curInput.ClientID, curInput.ClientSecret)),
				Name:         curInput.Name,
			}

			savedTokenStores = AddToken(ts, savedTokenStores)

			fmt.Println("Token for  " + curInput.Name + " saved in the TokenStore")

		}

	}

	if len(savedTokenStores) > 0 {
		fmt.Printf("Choose your tokenStore. Enter number b/w 0  -  %d \n", len(savedTokenStores)-1)
		for index := 0; index < len(savedTokenStores); index++ {
			fmt.Printf("%d. %s \n ", index, savedTokenStores[index].Name)
		}
		var choice int
		_, err := fmt.Scanf("%d", &choice)
		if err != nil {
			panic("I was expecting a number ")
		}
		if !(choice >= 0 && choice < len(savedTokenStores)) {
			fmt.Printf("Enter a number b/w 0 - %d \n", len(savedTokenStores)-1)
			fmt.Println("Retry")
			return
		}

		var askedToken = savedTokenStores[choice]
		if askedToken.TokenInfo.ShouldRenew() {
			askedToken.TokenInfo = getTokenInfo(getTokenResponse(askedToken.UaaURL, askedToken.ClientID, askedToken.ClientSecret))
			savedTokenStores[choice] = askedToken
			UpdateTokenStore(savedTokenStores)

		}
		fmt.Println(askedToken.TokenInfo.AccessToken)

	} else {
		fmt.Println("No Tokens are currently Managed. Use the -add option to add single token or -add-from-file=<file_path> to load from file")
	}

}

func getTokenInfo(tokenResponse map[string]interface{}) TokenInfo {

	tokenInfo := TokenInfo{}

	if tkn, ok := tokenResponse["access_token"].(string); ok {
		tokenInfo.AccessToken = tkn
	} else {
		panic("Unable to parse access_token")
	}

	tokenInfo.CreatedTime = time.Now()
	if expDur, ok := tokenResponse["expires_in"].(float64); ok {
		sec := fmt.Sprintf("%ds", int(expDur))
		dur, _ := time.ParseDuration(sec)
		tokenInfo.ExpiresAt = tokenInfo.CreatedTime.Add(dur)
	} else {
		panic("Unable  to parse expires_in token  ")
	}
	return tokenInfo
}

func (t TokenInfo) ShouldRenew() bool {

	if time.Now().After(t.ExpiresAt) {
		return true
	} else {
		return false
	}
}

func getTokenResponse(UaaURL, ClientID, ClientSecret string) map[string]interface{} {
	req, err := http.NewRequest("POST", UaaURL+"/oauth/token?grant_type=client_credentials", nil)
	if err != nil {
		fmt.Println("Unable to create new Post request")
		fmt.Println(err.Error())
	}

	req.SetBasicAuth(ClientID, ClientSecret)

	client := &http.Client{}

	res, error := client.Do(req)
	if error != nil {
		fmt.Println("Caught Error while trying to fetch token ")
	}
	defer res.Body.Close()

	var jsonResponse map[string]interface{}

	json.NewDecoder(res.Body).Decode(&jsonResponse)
	return jsonResponse
}

func LoadSavedTokenStores() []TokenStore {
	data, err := ioutil.ReadFile(Token_Store_File)
	var tsStores = make([]TokenStore, 0)
	if err != nil {
		file, _ := json.Marshal(tsStores)
		ioutil.WriteFile(Token_Store_File, file, 0644)
	} else {
		json.Unmarshal(data, &tsStores)

	}
	return tsStores

}

func AddToken(ts TokenStore, tss []TokenStore) []TokenStore {

	tss = append(tss, ts)
	file, _ := json.Marshal(tss)
	error := ioutil.WriteFile(Token_Store_File, file, 0644)
	if error != nil {
		fmt.Println(error.Error())
	}
	return tss

}

func UpdateTokenStore(tss []TokenStore) {
	file, _ := json.Marshal(tss)
	error := ioutil.WriteFile(Token_Store_File, file, 0644)
	if error != nil {
		fmt.Println(error.Error())
	}
}
func NameTaken(tss []TokenStore, name string) bool {
	var isNameTaken = false
	for index := 0; index < len(tss); index++ {
		if tss[index].Name == name {
			isNameTaken = true
			break
		}
	}
	return isNameTaken
}
