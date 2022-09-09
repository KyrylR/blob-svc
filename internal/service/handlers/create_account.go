package handlers

import (
	"blob-svc/internal/data"
	"blob-svc/internal/service/helpers"
	"blob-svc/internal/service/requests"
	"bytes"
	"encoding/json"
	"fmt"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/go/keypair"
	"gitlab.com/tokend/go/network"
	"gitlab.com/tokend/go/xdrbuild"
	"io/ioutil"
	"net/http"
)

func CreateAccount(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewCreateAccountRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	resp, err := http.Get("http://127.0.0.1:80")
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get horizon data")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	var horizonInfo data.HorizonInfo
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &horizonInfo); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}

	txBuilder := xdrbuild.NewBuilder(horizonInfo.NetworkPassphrase, int64(horizonInfo.TxExpirationPeriod))

	//sourceKP, err := keypair.FromRawSeed(network.ID(request.Data.ID))
	sourceKP, err := keypair.Random()
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	tx := txBuilder.Transaction(sourceKP)
	var signersData []xdrbuild.SignerData
	signersData = append(signersData, xdrbuild.SignerData{
		PublicKey: sourceKP.Address(),
		RoleID:    uint64(request.SignerData[0].Attributes.RoleID),
		Weight:    uint32(request.SignerData[0].Attributes.Weight),
		Identity:  uint32(request.SignerData[0].Attributes.Identity),
		Details:   json.RawMessage{},
	})
	var createAccount = xdrbuild.CreateAccount{
		Destination: sourceKP.Address(),
		Referrer:    nil,
		RoleID:      1,
		Signers:     signersData,
	}
	tx.Op(&createAccount)

	signerKP, err := keypair.FromRawSeed(network.ID(horizonInfo.NetworkPassphrase))
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	envelope, err := tx.Sign(signerKP).Marshal()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to marshal transaction")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	values := map[string]string{"tx": envelope}

	jsonValue, _ := json.Marshal(values)

	resp, err = http.Post(
		"http://localhost:8000/_/api/v3/transactions",
		"application/json",
		bytes.NewBuffer(jsonValue))
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to publish transaction")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, resp)
}
