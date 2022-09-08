package helpers

import (
	"blob-svc/internal/data"
	"blob-svc/internal/data/horizon"
	"blob-svc/internal/service/accountcreator"
	"blob-svc/xdrbuild"
	"context"
	horizon2 "gitlab.com/tokend/horizon-connector"
	"gitlab.com/tokend/keypair"
	"gitlab.com/tokend/regources"
	"net/http"

	"gitlab.com/distributed_lab/logan/v3"
)

type ctxKey int

const (
	logCtxKey ctxKey = iota
	blobsQCtxKey
	accountQCtxKey
	coreInfoCtxKey
	txBuilderCtxKey
	systemSettingsCtxKey
	horizonConnectorCtxKey
)

func AccountCreator(r *http.Request) accountcreator.AccountCreator {
	txbuilderbuilder := r.Context().Value(txBuilderCtxKey).(data.Infobuilder)
	info := r.Context().Value(coreInfoCtxKey).(data.Info)
	master := keypair.MustParseAddress(CoreInfo(r).MasterAccountID)
	tx := txbuilderbuilder(info, master)

	return accountcreator.New(
		tx,
		Horizon(r),
		SystemSettings(r),
	)
}

func CtxLog(entry *logan.Entry) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, logCtxKey, entry)
	}
}

func Log(r *http.Request) *logan.Entry {
	return r.Context().Value(logCtxKey).(*logan.Entry)
}

func CtxBlobsQ(entry data.BlobsQ) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, blobsQCtxKey, entry)
	}
}

func BlobsQ(r *http.Request) data.BlobsQ {
	return r.Context().Value(blobsQCtxKey).(data.BlobsQ).New()
}

func CtxAccountQ(q *horizon.AccountQ) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, accountQCtxKey, q)
	}
}

func AccountQ(r *http.Request) *horizon.AccountQ {
	return r.Context().Value(accountQCtxKey).(*horizon.AccountQ)
}

func CtxCoreInfo(s data.Info) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, coreInfoCtxKey, s)
	}
}

func CoreInfo(r *http.Request) *regources.Info {
	info, err := r.Context().Value(coreInfoCtxKey).(data.Info).Info()
	if err != nil {
		//TODO handle error
		panic(err)
	}
	return info
}

func CtxTransaction(txbuilder data.Infobuilder) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		ctx = context.WithValue(ctx, txBuilderCtxKey, txbuilder)
		return ctx
	}
}

func Transaction(r *http.Request) *xdrbuild.Transaction {
	txbuilderbuilder := r.Context().Value(txBuilderCtxKey).(data.Infobuilder)
	info := r.Context().Value(coreInfoCtxKey).(data.Info)
	master := keypair.MustParseAddress(CoreInfo(r).MasterAccountID)
	return txbuilderbuilder(info, master)
}

func TransactionWithSource(r *http.Request, source keypair.Address) *xdrbuild.Transaction {
	txbuilderbuilder := r.Context().Value(txBuilderCtxKey).(data.Infobuilder)
	info := r.Context().Value(coreInfoCtxKey).(data.Info)
	return txbuilderbuilder(info, source)
}

func CtxHorizon(q *horizon2.Connector) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, horizonConnectorCtxKey, q)
	}
}

func Horizon(r *http.Request) *horizon2.Connector {
	return r.Context().Value(horizonConnectorCtxKey).(*horizon2.Connector).Clone()
}

func CtxSystemSettings(q data.SystemSettings) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, systemSettingsCtxKey, q)
	}
}

func SystemSettings(r *http.Request) data.SystemSettings {
	return r.Context().Value(systemSettingsCtxKey).(data.SystemSettings)
}
