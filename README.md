# uaats
A Tool to Create and Refresh OAuthTokens for CF-UAA. UAA Token Store

## Overview:

  When working with multiple [UAA](https://docs.cloudfoundry.org/uaa/uaa-overview.html) the user would need to request Access Tokens from the UAA [Token endpoint](http://docs.cloudfoundry.org/api/uaa/version/4.21.0/#authorization-code-grant). These tokens have an expiration time. When working with mulitple UAAs it's tiresome to Fetch and Referesh (when expired) the tokens by Posting to the UAA Token Endpoint and constructing Authorization Headers. 
  This would frustrate when in the middle of a debugging session. Suppose you were testing the RESTful Apis either through command line and through a tool like POSTMAN. It would helpful if the tokens were easier to fetch and also automatically refresh upon expiry.
  
## Usage Modes:
 Currently there are two ways you can use this tool.
 1. Standalone mode.
 2. CloudFoundry Plugin Mode
  
## How to build:

### Standalone mode
If you are running on Mac Machine, then you may choose to run the uaats command directly.
Otherwise, download the code and build using the below command.
   `$ cd stand_alone`
   `$ go build -o uaats`
### CF Plugin Mode
   `$ cd cf_plugin`
   `$ go build -o cf_uaats_plugin`
   `$ cf install-plugin cf_uaats_plugin`
 
## How to use the tool:
  
  For standalone mode directly run the [uaats binary for mac](https://github.com/deeyes24/uaa_token_store/tree/master/stand_alone).
  For CloudFoundry you can run as "cf uaats" followed by below mentioned arguments. You can use this [prebuilt plugin for mac](https://github.com/deeyes24/uaa_token_store/tree/master/cf_plugin)
  
### Scenario1: Inital State

`$ uaats`

```
Output:
 No Tokens are currently Managed. Use the -add option to add single token or -add-from-file=<file_path> to load from file
 ```
 
### Scenario2: Adding a single TokenInfo

`$ uaats -add`

```
Output:
Enter UAA URL <replace_with_uaa_url>
Enter the clientId  <relace_with_client_id>
Enter the clientSecret <replace_with_client_secret>
Enter a unique name for the token store <replace_with_a_unique_name>
Access Token is <OUTPUT_FOR_ACCESS_TOKEN_IS_DISPLAYED_HERE>
```


### Scenario3: After adding tokeninfo

`$  uaats`

```Output:
Choose your tokenStore. Enter number b/w 0  -  0 
0. <unique_name_that_was_saved_in_scenario1_is_display>
```


### Scenario4: Adding with same name
```
Output: 
Name <token_store_name> is already taken. Try another unique name. Enter a unique name for the token store
```

### Scenario5: Adding from a json file.

`$ uaats -add-from-file`


```
Output:
flag needs an argument: -add-from-file
Usage of uaats:
  -add
        Add token to the TokenStore
  -add-from-file string
        Add to the TokenStore from file path. Ex :
                        [
                                {
                                        "name" : "uniqueName",
                                        "uaaURL": "actual UAA Url",
                                        "clientId": "clientid ",
                                        "clientSecret":"ClientSecret"
                                }
                        ]
 ```
                        

Using the sample output create file with multiple token store details.
 `$ uaats -add-from-file=input.json`
 
`$ uaats -add-from-file="<replace_with_token_store_input.json_as_created_from_above>"`

```
Output:
Token for  cf3-ad-app-client saved in the TokenStore
Token for  em-htc-app-client saved in the TokenStore
Token for  em-htc-admin saved in the TokenStore
Token for  default-app-client-int saved in the TokenStore
Choose your tokenStore. Enter number b/w 0  -  4 
 0. em-htc-admin 
 1. cf3-ad-app-client 
 2. em-htc-app-client 
 3. em-htc-admin 
 4. default-app-client-int 
 ```
 
 After adding you may immediately choose to select the TokenStore from available options or exit.
 

Thank You.
 
