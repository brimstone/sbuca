package server

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	//"github.com/brimstone/sbuca/x509util"

	"github.com/brimstone/sbuca/ca"
)

var config map[string]string

func Run(myConfig map[string]string) {

	config = myConfig
	m := martini.Classic()
	m.Use(render.Renderer())

	//FIXME
	ca.NewCA(config["root-dir"])

	m.Group("", func(r martini.Router) {
		r.Get("/", getRoot)
		r.Get("/ca/certificate", getCA)
		r.Get("/certificates/:id", getCertificates)
		r.Post("/certificates", authorizeSigning, postCertificates)
	})

	m.RunOnAddr(config["address"])

}
