/*
Package rpc implements bridge to Lachesis full node API interface.

We recommend using local IPC for fast and the most efficient inter-process communication between the API server
and an Opera/Lachesis node. Any remote RPC connection will work, but the performance may be significantly degraded
by extra networking overhead of remote RPC calls.

You should also consider security implications of opening Lachesis RPC interface for a remote access.
If you considering it as your deployment strategy, you should establish encrypted channel between the API server
and Lachesis RPC interface with connection limited to specified endpoints.

We strongly discourage opening Lachesis RPC interface for unrestricted Internet access.
*/
package rpc

import (
	"encoding/base64"
	"encoding/json"
	"fantom-api-graphql/internal/types"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// requestTimeout is number of seconds we wait for the information request to finish.
const requestTimeout = 5 * time.Second

// nameCheckRegex is the expression used to check for validtor name validity
var stiNameCheckRegex = regexp.MustCompile(`^[\w\d\s.\-_'$()]+$`)

// ValidatorInfo extracts extended information for a validator.
func (ftm *FtmBridge) ValidatorInfo(id *hexutil.Big) (*types.ValidatorInfo, error) {
	if id == nil {
		return nil, fmt.Errorf("validator ID not given")
	}

	// keep track of the operation
	ftm.log.Debugf("loading staker information for staker #%d", id.ToInt().Uint64())

	// instantiate the contract and display its name
	stUrl, err := ftm.SfcContract().GetValidatorInfo(nil, (*big.Int)(id))
	if err != nil {
		ftm.log.Errorf("failed to get the validator information: %v", err)
		return nil, err
	}

	// var url string
	if len(stUrl) == 0 {
		ftm.log.Debugf("no information for validator #%d", id.ToInt().Uint64())
		return nil, nil
	}

	// try to download JSON for the info
	return ftm.downloadValidatorInfo(stUrl)
}

// downloadValidatorInfo tries to download validator information from the given URL address.
func (ftm *FtmBridge) downloadValidatorInfo(stUrl string) (*types.ValidatorInfo, error) {
	var data []byte
	// check for data url
	if strings.HasPrefix(stUrl, "data:application/json;base64,") {
		ftm.log.Debugf("using validator info address from dataUrl [%s]", stUrl)

		var err error
		// extract base64 encoded json
		data, err = base64.StdEncoding.DecodeString(stUrl[29:])
		if err != nil {
			return nil, err
		}
	} else {
		// log what we are about to do
		ftm.log.Debugf("downloading validator info address [%s]", stUrl)

		// make a http client
		cl := http.Client{Timeout: requestTimeout}

		// prep request
		req, err := http.NewRequest(http.MethodGet, stUrl, nil)
		if err != nil {
			ftm.log.Errorf("can not request given validator info url; %s", err.Error())
			return nil, err
		}

		// be honest, set agent
		req.Header.Set("User-Agent", "Camino GraphQL API Server")

		// process the request
		res, err := cl.Do(req)
		if err != nil {
			ftm.log.Errorf("can not download validator info; %s", err.Error())
			return nil, err
		}

		// read the response
		data, err = ioutil.ReadAll(res.Body)
		if err != nil {
			ftm.log.Errorf("can not read validator info response; %s", err.Error())
			return nil, err
		}
	}

	// try to parse
	var info types.ValidatorInfo
	err := json.Unmarshal(data, &info)
	if err != nil {
		ftm.log.Errorf("invalid response for validator info; %s", err.Error())
		return nil, err
	}

	// do we have anything?
	if !ftm.isValidValidatorInfo(&info) {
		ftm.log.Errorf("invalid response for validator info [%s]", stUrl)
		return nil, err
	}

	ftm.log.Debugf("found validator [%s]", *info.Name)
	return &info, nil
}

// isValidValidatorInfo check if the validator information is valid and can be used.
func (ftm *FtmBridge) isValidValidatorInfo(info *types.ValidatorInfo) bool {
	// name must be available
	if nil == info.Name || 0 == len(*info.Name) || !stiNameCheckRegex.Match([]byte(*info.Name)) {
		ftm.log.Error("validator name not valid")
		return false
	}

	// check the logo URL
	if !isValidValidatorInfoUrl(info.LogoUrl, true) {
		ftm.log.Error("validator logo URL not valid")
		return false
	}

	// check the website
	if !isValidValidatorInfoUrl(info.Website, false) {
		ftm.log.Error("validator website URL not valid")
		return false
	}

	// check the contact URL
	if !isValidValidatorInfoUrl(info.Contact, false) {
		ftm.log.Error("validator contact URL not valid")
		return false
	}
	return true
}

// isValidValidatorInfoUrl validates the given URL address from the validator info.
func isValidValidatorInfoUrl(addr *string, reqHttps bool) bool {
	// do we even have an URL; it's ok if not
	if nil == addr || 0 == len(*addr) {
		return true
	}

	// try to decode the address
	u, err := url.ParseRequestURI(*addr)
	if err != nil || u.Scheme == "" || (reqHttps && u.Scheme != "https") {
		return false
	}
	return true
}
