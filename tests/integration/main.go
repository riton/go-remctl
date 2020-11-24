package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	k5client "github.com/jcmturner/gokrb5/v8/client"
	k5config "github.com/jcmturner/gokrb5/v8/config"
	"github.com/jcmturner/gokrb5/v8/credentials"
	"github.com/jcmturner/gokrb5/v8/gssapi"
	"github.com/jcmturner/gokrb5/v8/spnego"
	"github.com/pkg/errors"
)

type RemctlClient struct {
	client *k5client.Client
	spn    string
}

type RemctlGSSAPIToken struct {
	Init         bool
	Resp         bool
	NegTokenInit spnego.NegTokenInit
	NegTokenResp spnego.NegTokenResp
	context      context.Context
}

func (r *RemctlClient) InitSecContext() (gssapi.ContextToken, error) {
	tkt, key, err := r.client.GetServiceTicket(r.spn)
	if err != nil {
		return &RemctlGSSAPIToken{}, err
	}
	negTokenInit, err := spnego.NewNegTokenInitKRB5(r.client, tkt, key)
	if err != nil {
		return &RemctlGSSAPIToken{}, errors.Wrap(err, "could not create NegTokenInit")
	}

	return &RemctlGSSAPIToken{
		Init:         true,
		NegTokenInit: negTokenInit,
	}, nil
}

func main() {

	remoteHost := "ccpphook.in2p3.fr"
	realm := "CC.IN2P3.FR"

	krb5ccname := os.Getenv("KRB5CCNAME")
	if krb5ccname == "" {
		panic("missing KRB5CCNAME")
	}

	fields := strings.SplitN(krb5ccname, ":", 2)
	if fields[0] != "FILE" {
		panic("only FILE ccache type is supported")
	}

	ccache, err := credentials.LoadCCache(fields[1])
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", ccache)

	krb5Config, err := k5config.Load("/etc/krb5.conf")
	if err != nil {
		panic(err)
	}

	krb5Client, err := k5client.NewFromCCache(
		ccache,
		krb5Config,
	)
	if err != nil {
		panic(err)
	}

	rClient := RemctlClient{
		client: krb5Client,
		spn:    fmt.Sprintf("host/%s@%s", remoteHost, realm),
	}

	gctx, err := rClient.InitSecContext()
	if err != nil {
		panic(err)
	}
}
