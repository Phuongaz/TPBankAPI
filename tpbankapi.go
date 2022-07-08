package tpbankapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
)

const (
	TOKEN_URL   = "https://ebank.tpb.vn/gateway/api/auth/login"
	HISTORY_URL = "https://ebank.tpb.vn/gateway/api/smart-search-presentation-service/v1/account-transactions/find"
)

type TPBank struct {
	Username string
	Password string
	Token    string
}

func NewTPBank(username string, password string) *TPBank {
	return &TPBank{
		Username: username,
		Password: password,
	}
}

func (tpbank *TPBank) Login() error {
	postBody := map[string]string{
		"username": tpbank.Username,
		"password": tpbank.Password,
	}
	body, err := json.Marshal(postBody)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", TOKEN_URL, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("DEVICE_ID", "LYjkjqGZ3HhGP5520GxPP2j94RDMC7Xje77MI7"+strconv.Itoa(rand.Intn(100000000)))
	req.Header.Set("PLATFORM_VERSION", "91")
	req.Header.Set("DEVICE_NAME", "Chrome")
	req.Header.Set("SOURCE-APP", "HYDRO")
	req.Header.Set("Authorization", "Bearer")
	req.Header.Set("Accept", "aapplication/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,vi;q=0.8")
	req.Header.Set("Referer", " https://ebank.tpb.vn/retail/vX/login?returnUrl=%2Fmain")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.132 Safari/537.36")
	req.Header.Set("PLATFORM_NAME", "WEB")
	req.Header.Set("APP_VERSION", "1.3")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err2 := ioutil.ReadAll(resp.Body)

	if err2 != nil {
		return err2
	}
	var json_resp map[string]string
	json.Unmarshal(body, &json_resp)
	if len(json_resp) == 0 {
		return errors.New("login failed")
	}
	tpbank.Token = json_resp["access_token"]
	return nil
}

// GetHistory returns the history of the account
func (tpbank *TPBank) GetHistory(stk_bank string, fromDate string, toDate string) ([]byte, error) {
	postBody := map[string]string{
		"fromDate":  fromDate,
		"toDate":    toDate,
		"currency":  "VND",
		"accountNo": stk_bank,
		"keyword":   "",
	}
	body, err := json.Marshal(postBody)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", HISTORY_URL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	if tpbank.Token == "" {
		err := tpbank.Login()
		if err != nil {
			return nil, err
		}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("connection", "keep-alive")
	req.Header.Set("DEVICE_ID", "LYjkjqGZ3HhGP5520GxPP2j94RDMC7Xje77MI75RYBVR")
	req.Header.Set("PLATFORM_VERSION", "91")
	req.Header.Set("DEVICE_NAME", "Chrome")
	req.Header.Set("SOURCE-APP", "HYDRO")
	req.Header.Set("Authorization", "Bearer "+tpbank.Token)
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,vi;q=0.8")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.132 Safari/537.36")
	req.Header.Set("PLATFORM_NAME", "WEB")
	req.Header.Set("APP_VERSION", "1.3")
	req.Header.Set("Referer", " https://ebank.tpb.vn/retail/vX/main/inquiry/account/transaction?id="+stk_bank)
	req.Header.Set("Cookie", "_ga=GA1.2.1679888794.1623516; _gid=GA1.2.580582711.16277; _gcl_au=1.1.756417552.1626666; Authorization="+tpbank.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		return nil, err2
	}
	return body, nil
	//creditDebitIndicator:
	//	+ DBIT: Tru tien
	//	+ CRDT: Cong tien
}
