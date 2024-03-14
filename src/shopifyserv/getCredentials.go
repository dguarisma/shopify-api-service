package shopifyserv

import (
	"fmt"
	"net/http"
	"strings"
)

func New(version string, credentials []ShopifyCredentials) ShopifyService {
	wareHouses := make(map[string]WareHouseShopify, len(credentials))

	for _, cred := range credentials {
		wareHouses[cred.Name] = WareHouseShopify{
			UrlBase:             cred.UrlBase,
			XShopifyAccessToken: cred.XShopifyAccessToken,
		}
	}

	return &ShopifyServ{
		warehouses: wareHouses,
		version:    version,
		client:     &http.Client{},
	}
}

func GetCredentials(credentialsStrings ...string) []ShopifyCredentials {
	credential := make([]ShopifyCredentials, len(credentialsStrings))
	for i, str := range credentialsStrings {
		cred, err := GetCredential(str)
		if err != nil {
			panic(err)
		}
		credential[i] = cred
	}

	return credential
}

func GetCredential(credentialString string) (credential ShopifyCredentials, err error) {
	if credentialString == "" {
		return credential, fmt.Errorf(ErrEmptyCredentials)
	}

	cred := strings.Split(credentialString, ",")
	if len(cred) != 3 {
		return credential, fmt.Errorf(ErrNotEnoughtCredentials)
	}

	credential = ShopifyCredentials{
		Name:                cred[0],
		UrlBase:             cred[1],
		XShopifyAccessToken: cred[2],
	}

	return
}
