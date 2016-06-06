// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/FactomProject/fctwallet/Wallet"
	"github.com/FactomProject/web"

	"github.com/FactomProject/factomd/common/primitives"
	"github.com/FactomProject/factomd/wsapi"
)

func FctBalance(adr string) (int64, error) {
	err := refresh()
	if err != nil {
		return 0, err
	}
	return Wallet.FactoidBalance(adr)
}

func ECBalance(adr string) (int64, error) {
	return Wallet.ECBalance(adr)
}

func HandleEntryCreditBalance(ctx *web.Context, adr string) {
	req := primitives.NewJSON2Request("entry-credit-balance", 1, AddressRequest{Address: adr})

	jsonResp, jsonError := HandleV2Request(req)
	if jsonError != nil {
		reportResults(ctx, jsonError.Message, false)
		return
	}

	str := fmt.Sprintf("%d", jsonResp.Result.(*EntryCreditBalanceResponse).Balance)
	reportResults(ctx, str, true)
}

func HandleV2EntryCreditBalance(params interface{}) (interface{}, *primitives.JSONError) {
	req := new(AddressRequest)
	err := wsapi.MapToObject(params, req)
	if err != nil {
		return nil, wsapi.NewInvalidParamsError()
	}
	adr := req.Address

	v, err := ECBalance(adr)
	if err != nil {
		return nil, wsapi.NewInvalidParamsError()
	}
	resp := new(EntryCreditBalanceResponse)
	resp.Balance = v

	return resp, nil
}

func HandleFactoidBalance(ctx *web.Context, adr string) {
	req := primitives.NewJSON2Request("factoid-balance", 1, AddressRequest{Address: adr})

	jsonResp, jsonError := HandleV2Request(req)
	if jsonError != nil {
		reportResults(ctx, jsonError.Message, false)
		return
	}

	str := fmt.Sprintf("%d", jsonResp.Result.(*FactoidBalanceResponse).Balance)
	reportResults(ctx, str, true)
}

func HandleV2FactoidBalance(params interface{}) (interface{}, *primitives.JSONError) {
	req := new(AddressRequest)
	err := wsapi.MapToObject(params, req)
	if err != nil {
		return nil, wsapi.NewInvalidParamsError()
	}
	adr := req.Address

	v, err := FctBalance(adr)
	if err != nil {
		return nil, wsapi.NewInvalidParamsError()
	}
	resp := new(FactoidBalanceResponse)
	resp.Balance = v

	return resp, nil
}

func HandleResolveAddress(ctx *web.Context, adr string) {
	req := primitives.NewJSON2Request("resolve-address", 1, AddressRequest{Address: adr})

	jsonResp, jsonError := HandleV2Request(req)
	if jsonError != nil {
		reportResults(ctx, jsonError.Message, false)
		return
	}

	type x struct {
		Fct, Ec string
	}

	t := new(x)
	t.Fct = jsonResp.Result.(*ResolveAddressResponse).FactoidAddress
	t.Ec = jsonResp.Result.(*ResolveAddressResponse).EntryCreditAddress
	p, err := json.Marshal(t)
	if err != nil {
		reportResults(ctx, err.Error(), false)
		return
	}

	reportResults(ctx, string(p), true)
}

func HandleV2ResolveAddress(params interface{}) (interface{}, *primitives.JSONError) {
	req := new(AddressRequest)
	err := wsapi.MapToObject(params, req)
	if err != nil {
		return nil, wsapi.NewInvalidParamsError()
	}
	adr := req.Address

	fAddress, ecAddress, err := Wallet.NetkiResolve(adr)
	if err != nil {
		return nil, wsapi.NewCustomInternalError(err.Error())
	}

	resp := new(ResolveAddressResponse)
	resp.FactoidAddress = fAddress
	resp.EntryCreditAddress = ecAddress

	return resp, nil
}
