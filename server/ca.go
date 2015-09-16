package server

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	//"github.com/brimstone/sbuca/x509util"
	"net/http"

	"github.com/brimstone/sbuca/ca"
)

func getCA(req *http.Request, params martini.Params, r render.Render) {

	format := req.URL.Query().Get("format")

	newCA, err := ca.NewCA(config["root-dir"])
	if err != nil {
		panic(err)
	}

	pem, err := newCA.Certificate.ToPEM()
	if err != nil {
		panic(err)
	}

	if format == "file" {
		r.Data(200, pem)
	} else {
		r.JSON(200, map[string]interface{}{
			"ca": map[string]interface{}{
				"crt": string(pem),
			},
		})
	}
}
