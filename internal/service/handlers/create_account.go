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
	"gitlab.com/tokend/go/strkey"
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
	body, err := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &horizonInfo); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}

	txBuilder := xdrbuild.NewBuilder(horizonInfo.NetworkPassphrase, int64(horizonInfo.TxExpirationPeriod))

	seed, err := strkey.Decode(strkey.VersionByteSeed, "SAMJKTZVW5UOHCDK5INYJNORF2HRKYI72M5XSZCBYAHQHR34FFR4Z6G4")
	var rawSeed [32]byte
	copy(rawSeed[:], seed)
	sourceKP, err := keypair.FromRawSeed(rawSeed)
	destKP, err := keypair.Random()
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	tx := txBuilder.Transaction(sourceKP)
	var signersData []xdrbuild.SignerData
	details := json.RawMessage("{}")
	signersData = append(signersData, xdrbuild.SignerData{
		PublicKey: destKP.Address(),
		RoleID:    uint64(request.SignerData[0].Attributes.RoleID),
		Weight:    uint32(request.SignerData[0].Attributes.Weight),
		Identity:  uint32(request.SignerData[0].Attributes.Identity),
		Details:   details,
	})
	var createAccount = xdrbuild.CreateAccount{
		Destination: destKP.Address(),
		Referrer:    nil,
		RoleID:      1,
		Signers:     signersData,
	}
	tx.Op(&createAccount)

	envelope, err := tx.Sign(sourceKP).Marshal()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to marshal transaction")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	values := map[string]string{"tx": envelope}

	jsonValue, _ := json.Marshal(values)

	response, err := http.Post(
		"http://localhost:8000/_/api/v3/transactions",
		"application/json",
		bytes.NewBuffer(jsonValue))
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to publish transaction")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(response.StatusCode)
}
