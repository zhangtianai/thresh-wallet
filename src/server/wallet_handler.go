// thresh-wallet
//
// Copyright 2019 by KeyFuse
//
// GPLv3 License

package server

import (
	"encoding/json"
	"net/http"

	"proto"
)

func (h *Handler) walletCheck(w http.ResponseWriter, r *http.Request) {
	var userExists bool
	var backupExists bool
	var backupTimestamp int64
	var backupCloudService string

	log := h.log
	wdb := h.wdb
	resp := newResponse(log, w)

	// UID.
	uid, err := h.userinfo("walletCheck", r)
	if err != nil {
		log.Error("api.wallet.check.uid.error:%+v", err)
		resp.writeError(err)
		return
	}

	// Request.
	req := &proto.WalletCheckRequest{}
	err = json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		log.Error("api.wallet[%v].check.decode.body.error:%+v", uid, err)
		resp.writeError(err)
		return
	}
	log.Info("api.wallet.check.req:%+v", req)

	wallet := wdb.Wallet(uid)
	if wallet != nil {
		userExists = true
		backupExists = (wallet.Backup.EncryptedPrvKey != "")
		backupTimestamp = wallet.Backup.Time
		backupCloudService = wallet.Backup.CloudService
	}
	// Response.
	rsp := proto.WalletCheckResponse{
		UserExists:         userExists,
		BackupExists:       backupExists,
		BackupTimestamp:    backupTimestamp,
		BackupCloudService: backupCloudService,
	}
	log.Info("api.wallet.check.rsp:%+v", rsp)
	resp.writeJSON(rsp)
}

func (h *Handler) walletCreate(w http.ResponseWriter, r *http.Request) {
	log := h.log
	wdb := h.wdb
	resp := newResponse(log, w)

	// UID.
	uid, err := h.userinfo("walletCreate", r)
	if err != nil {
		log.Error("api.wallet.create.uid.error:%+v", err)
		resp.writeError(err)
		return
	}

	// Request.
	req := &proto.WalletCreateRequest{}
	err = json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		log.Error("api.wallet[%v].create.decode.body.error:%+v", uid, err)
		resp.writeError(err)
		return
	}
	log.Info("api.wallet.create.req:%+v", req)

	// Verify the pub/prv key pairing.
	if err := verifyPubKey(req.MasterPubKey, req.Signature); err != nil {
		log.Error("api.wallet[%v].create.verify.pubkey.error:%+v", uid, err)
		resp.writeErrorWithStatus(400, err)
		return
	}

	// Create wallet.
	if err := wdb.CreateWallet(uid, req.MasterPubKey); err != nil {
		log.Error("api.wallet[%v].wdb.create.error:%+v", uid, err)
		resp.writeError(err)
		return
	}
	// Response.
	rsp := proto.WalletCreateResponse{}
	log.Info("api.wallet.create.rsp:%+v", rsp)
	resp.writeJSON(rsp)
}

func (h *Handler) walletBalance(w http.ResponseWriter, r *http.Request) {
	log := h.log
	wdb := h.wdb
	resp := newResponse(log, w)

	// UID.
	uid, err := h.userinfo("walletBalance", r)
	if err != nil {
		log.Error("api.wallet.balance.uid.error:%+v", err)
		resp.writeError(err)
		return
	}

	// Request.
	req := &proto.WalletBalanceRequest{}
	err = json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		log.Error("api.wallet[%v].balance.decode.body.error:%+v", uid, err)
		resp.writeError(err)
		return
	}
	log.Info("api.wallet.balance.req:%+v", req)

	// Balance.
	balance, err := wdb.Balance(uid)
	if err != nil {
		log.Error("api.wallet.balance.wdb.balance.error:%+v", err)
		resp.writeError(err)
		return
	}

	// Response.
	rsp := proto.WalletBalanceResponse{
		CoinValue: balance.TotalBalance,
	}
	resp.writeJSON(rsp)
}

func (h *Handler) walletUnspent(w http.ResponseWriter, r *http.Request) {
	log := h.log
	wdb := h.wdb
	resp := newResponse(log, w)

	// UID.
	uid, err := h.userinfo("walletUnspent", r)
	if err != nil {
		log.Error("api.wallet.unspent.uid.error:%+v", err)
		resp.writeError(err)
		return
	}

	// Request.
	req := &proto.WalletUnspentRequest{}
	err = json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		log.Error("api.wallet[%v].unspent.decode.body.error:%+v", uid, err)
		resp.writeError(err)
		return
	}
	log.Info("api.wallet.unspent.req:%+v", req)

	unspents, err := wdb.Unspents(uid, req.Amount)
	if err != nil {
		log.Error("api.wallet[%v].unspent.by.amount.error:%+v", uid, err)
		resp.writeError(err)
		return
	}

	var rsp []proto.WalletUnspentResponse
	for _, unspent := range unspents {
		rsp = append(rsp, proto.WalletUnspentResponse{
			Pos:          unspent.Pos,
			Txid:         unspent.Txid,
			Vout:         unspent.Vout,
			Value:        unspent.Value,
			Address:      unspent.Address,
			Confirmed:    unspent.Confirmed,
			SvrPubKey:    unspent.SvrPubKey,
			Scriptpubkey: unspent.Scriptpubkey,
		})
	}
	log.Info("api.wallet.unspent.rsp:%+v", rsp)
	resp.writeJSON(rsp)
}

func (h *Handler) walletTxs(w http.ResponseWriter, r *http.Request) {
	log := h.log
	wdb := h.wdb
	resp := newResponse(log, w)

	// UID.
	uid, err := h.userinfo("walletTxs", r)
	if err != nil {
		log.Error("api.wallet.txs.uid.error:%+v", err)
		resp.writeError(err)
		return
	}

	// Request.
	req := &proto.WalletTxsRequest{}
	err = json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		log.Error("api.wallet[%v].txs.decode.body.error:%+v", uid, err)
		resp.writeError(err)
		return
	}
	log.Info("api.wallet.txs.req:%+v", req)

	if req.Limit > 16 {
		req.Limit = 16
	}
	txs, err := wdb.Txs(uid, req.Offset, req.Limit)
	if err != nil {
		log.Error("api.wallet[%v].txs.error:%+v", uid, err)
		resp.writeError(err)
		return
	}

	var rsp []proto.WalletTxsResponse
	for _, tx := range txs {
		rsp = append(rsp, proto.WalletTxsResponse{
			Txid:        tx.Txid,
			Fee:         tx.Fee,
			Link:        tx.Link,
			Value:       tx.Value,
			Confirmed:   tx.Confirmed,
			BlockTime:   tx.BlockTime,
			BlockHeight: tx.BlockHeight,
		})
	}
	log.Info("api.wallet.txs.rsp:%+v", rsp)
	resp.writeJSON(rsp)
}

func (h *Handler) walletSendFees(w http.ResponseWriter, r *http.Request) {
	log := h.log
	wdb := h.wdb
	resp := newResponse(log, w)

	// UID.
	uid, err := h.userinfo("walletSendFees", r)
	if err != nil {
		log.Error("api.wallet.send.fees.uid.error:%+v", err)
		resp.writeError(err)
		return
	}

	// Request.
	req := &proto.WalletSendFeesRequest{}
	err = json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		log.Error("api.wallet[%v].send.fees.decode.body.error:%+v", uid, err)
		resp.writeError(err)
		return
	}
	log.Info("api.wallet.send.fees.req:%+v", req)

	fees, err := wdb.SendFees(uid, req.Priority, req.SendValue)
	if err != nil {
		log.Error("api.wallet[%v].send.fees.wdb.send.fees.error:%+v", uid, err)
		resp.writeError(err)
		return
	}

	rsp := &proto.WalletSendFeesResponse{
		Fees:          fees.Fees,
		TotalValue:    fees.TotalValue,
		SendableValue: fees.SendableValue,
	}
	log.Info("api.wallet.fees.rsp:%+v", rsp)
	resp.writeJSON(rsp)
}

func (h *Handler) walletPortfolio(w http.ResponseWriter, r *http.Request) {
	log := h.log
	wdb := h.wdb
	resp := newResponse(log, w)
	code := "CNY"

	// UID.
	uid, err := h.userinfo("walletPortfolio", r)
	if err != nil {
		log.Error("api.wallet.portfolio.uid.error:%+v", err)
		resp.writeError(err)
		return
	}

	// Request.
	req := &proto.WalletPortfolioRequest{}
	err = json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		log.Error("api.wallet[%v].portfolio.decode.body.error:%+v", uid, err)
		resp.writeError(err)
		return
	}
	log.Info("api.wallet.portfolio.req:%+v", req)

	if req.Code != "" {
		code = req.Code
	}
	ticker, err := wdb.store.getTicker(code)
	if err != nil {
		log.Error("api.wallet[%v].get.ticker.error:%+v", uid, err)
		resp.writeError(err)
		return
	}

	rsp := &proto.WalletPortfolioResponse{
		CoinSymbol:   "BTC",
		FiatSymbol:   ticker.Symbol,
		CurrentPrice: ticker.Last,
	}
	log.Info("api.wallet.portfolio.rsp:%+v", rsp)
	resp.writeJSON(rsp)
}

func (h *Handler) walletPushTx(w http.ResponseWriter, r *http.Request) {
	log := h.log
	wdb := h.wdb
	resp := newResponse(log, w)

	// UID.
	uid, err := h.userinfo("walletPushTx", r)
	if err != nil {
		log.Error("api.wallet.push.tx.uid.error:%+v", err)
		resp.writeError(err)
		return
	}

	// Request.
	req := &proto.TxPushRequest{}
	err = json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		log.Error("api.wallet[%v].push.tx.decode.body.error:%+v", uid, err)
		resp.writeError(err)
		return
	}
	log.Info("api.wallet.push.tx.req:%+v", req)

	txid, err := wdb.chain.PushTx(req.TxHex)
	if err != nil {
		log.Error("api.wallet[%v].push.tx.error:%+v", uid, err)
		resp.writeError(err)
		return
	}
	rsp := &proto.TxPushResponse{
		TxID: txid,
	}
	log.Info("api.wallet.push.tx.rsp:%+v", rsp)
	resp.writeJSON(rsp)
}
