package server

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	//"github.com/brimstone/sbuca/x509util"
	"net/http"
	"strconv"

	"github.com/brimstone/sbuca/ca"
	"github.com/brimstone/sbuca/pkix"
)

func getCertificates(req *http.Request, params martini.Params, r render.Render) {

	format := req.URL.Query().Get("format")

	newCA, err := ca.NewCA(config["root-dir"])
	if err != nil {
		panic(err)
	}

	id := params["id"]
	idInt, err := strconv.Atoi(id)
	if err != nil {
		r.JSON(401, map[string]interface{}{
			"result": "wrong id",
		})
		return
	}
	cert, err := newCA.GetCertificate(int64(idInt))
	if err != nil {
		r.JSON(401, map[string]interface{}{
			"result": "cannot get cert",
		})
		return
	}

	pem, err := cert.ToPEM()
	if err != nil {
		r.JSON(401, map[string]interface{}{
			"result": "cannot get cert",
		})
		return
	}

	if format == "file" {
		r.Data(200, pem)
	} else {
		r.JSON(200, map[string]interface{}{
			"certificate": map[string]interface{}{
				"id":  cert.GetSerialNumber().Int64(),
				"crt": string(pem),
				//"csr": csr,
			},
		})
	}

}

func postCertificates(req *http.Request, params martini.Params, r render.Render) {

	csrString := req.PostFormValue("csr")
	format := req.URL.Query().Get("format")

	csr, err := pkix.NewCertificateRequestFromPEM([]byte(csrString))
	if err != nil {
		panic(err)
	}

	newCA, err := ca.NewCA(config["root-dir"])
	if err != nil {
		panic(err)
	}

	cert, err := newCA.IssueCertificate(csr)
	if err != nil {
		panic(err)
	}

	certPem, err := cert.ToPEM()
	if err != nil {
		panic(err)
	}
	if format == "file" {
		r.Data(200, certPem)
	} else {
		r.JSON(200, map[string]interface{}{
			"certificate": map[string]interface{}{
				"id":  cert.GetSerialNumber().Int64(),
				"crt": string(certPem),
				//"csr": csr,
			},
		})
	}
}
