package handlers

import (
	"blob-svc/internal/service/requests"
	"net/http"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func CreateAccount(w http.ResponseWriter, r *http.Request) {
	_, err := requests.NewCreateAccountRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	//account, err := helpers.AccountQ(r).Account(request.Data.ID)
	//if err != nil {
	//	helpers.Log(r).WithError(err).Error("failed to get account")
	//	ape.RenderErr(w, problems.InternalError())
	//	return
	//}
	//
	//if account != nil {
	//	ape.RenderErr(w, problems.Conflict())
	//	return
	//}
	//
	//signers, err := getSigners(request)
	//if err != nil {
	//	ape.RenderErr(w, problems.BadRequest(err)...)
	//	return
	//}
	//
	//err = helpers.AccountCreator(r).CreateAccount(r.Context(), request.Data.ID, signers)
	//if err != nil {
	//	helpers.Log(r).WithError(err).Error("failed to create account")
	//	ape.RenderErr(w, problems.InternalError())
	//	return
	//}
	//
	//w.WriteHeader(http.StatusCreated)
}

//func getSigners(request resources.CreateAccountResponse) ([]data.AccountSigner, error) {
//	var signers []data.AccountSigner
//	for _, signerKey := range request.Data.Relationships.Signers.Data {
//		signer := request.Included.MustSigner(signerKey)
//		if signer == nil {
//			return nil, validation.Errors{"/included": errors.New("missed signer include")}
//		}
//		signers = append(signers, data.AccountSigner{
//			SignerID: signerKey.ID,
//			RoleID:   uint64(signer.Attributes.RoleId),
//			Weight:   uint32(signer.Attributes.Weight),
//			Identity: uint32(signer.Attributes.Identity),
//		})
//	}
//	return signers, nil
//}
