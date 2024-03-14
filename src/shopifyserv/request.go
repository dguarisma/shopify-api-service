package shopifyserv

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (ss *ShopifyServ) getRequest(url, token string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Shopify-Access-Token", token)

	data, status, err := ss.request(req)
	if err != nil {
		return nil, err
	}
	if status != 200 {
		return nil, ErrStatusUnexpected(200, status)
	}
	return data, nil
}

func (ss *ShopifyServ) request(req *http.Request) ([]byte, int, error) {
	res, err := ss.client.Do(req)
	if err != nil {
		return nil, 0, ErrResponse(err)
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, 0, ErrResponseFormat(err)
	}
	return data, res.StatusCode, nil
}

func (ss *ShopifyServ) graphqlQuery(query Query, warehouse WareHouseShopify) (data []byte, status int, err error) {
	b2, err := json.Marshal(query)
	if err != nil {
		return
	}

	url := fmt.Sprintf(
		"%s/admin/api/%s/graphql.json",
		warehouse.UrlBase,
		ss.version,
	)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b2))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("X-Shopify-Access-Token", warehouse.XShopifyAccessToken)

	res, err := ss.client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	data, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, 0, ErrResponseFormat(err)
	}
	return data, res.StatusCode, nil
}
