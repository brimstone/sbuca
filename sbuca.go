package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/brimstone/sbuca/pkix"
	"github.com/brimstone/sbuca/server"
	"github.com/codegangsta/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = "sbuca"
	app.Usage = "Simple But Useful CA"
	app.Commands = []cli.Command{

		{
			Name:  "server",
			Usage: "Run a CA server",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "dir",
					Value: ".",
					Usage: "Root directory for certificates",
				},
				cli.StringFlag{
					Name:  "address",
					Value: os.Getenv("HOST") + ":8600",
					Usage: "Token for cert signing functions",
				},
				cli.StringFlag{
					Name:  "admin-token",
					Value: "",
					Usage: "Token for administrative functions",
				},
				cli.StringFlag{
					Name:  "sign-token",
					Value: "",
					Usage: "Token for cert signing functions",
				},
			},
			Action: func(c *cli.Context) {
				config := map[string]string{
					"address":     c.String("address"),
					"root-dir":    c.String("dir"),
					"admin-token": c.String("admin-token"),
					"sign-token":  c.String("sign-token"),
				}
				server.Run(config)
			},
		},

		{
			Name:  "genkey",
			Usage: "Generate a RSA Private Key to STDOUT",
			Action: func(c *cli.Context) {
				key, err := pkix.NewKey()
				if err != nil {
					panic(err)
				}

				pem, err := key.ToPEM()
				if err != nil {
					panic(err)
				}

				fmt.Print(string(pem))
			},
		},

		{
			Name:  "gencsr",
			Usage: "Generate a Certificate Request to STDOUT",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "key",
					Usage: "RSA Private Key",
				},
			},
			Action: func(c *cli.Context) {
				keyName := c.String("key")
				if keyName == "" {
					fmt.Fprintln(os.Stderr, "[ERROR] Requere private key as parameter")
					return
				}

				key, err := pkix.NewKeyFromPrivateKeyPEMFile(keyName)
				if err != nil {
					fmt.Fprintln(os.Stderr, "[ERROR] Failed to generate CSR: "+err.Error())
					return
				}

				csr, err := pkix.NewCertificateRequest(key)
				if err != nil {
					fmt.Fprintln(os.Stderr, "[ERROR] Failed to generate CSR: "+err.Error())
					return
				}

				pem, err := csr.ToPEM()
				if err != nil {
					fmt.Fprintln(os.Stderr, "[ERROR] Failed to generate CSR: "+err.Error())
					return
				}

				fmt.Print(string(pem))
			},
		},

		{
			Name:  "submitcsr",
			Usage: "Submit a Certificate Request to CA",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "host",
					Usage: "Host ip & port",
				},
				cli.StringFlag{
					Name:  "format",
					Value: "cert",
					Usage: "output cert or id",
				},
				cli.StringFlag{
					Name:  "token",
					Usage: "Authorization Token",
				},
			},
			Action: func(c *cli.Context) {
				host := c.String("host")
				if host == "" {
					fmt.Fprintln(os.Stderr, "[ERROR] Requere host as parameter")
					return
				}

				format := c.String("format")
				if format != "cert" && format != "id" {
					fmt.Fprintln(os.Stderr, "[ERROR] format should be 'cert' or 'id'")
					return
				}

				args := c.Args()
				if len(args) == 0 {
					fmt.Fprintln(os.Stderr, "[ERROR] Should provide csr")
					return
				}
				csrName := c.Args().First()

				csr, err := pkix.NewCertificateRequestFromPEMFile(csrName)
				if err != nil {
					fmt.Fprintln(os.Stderr, "[ERROR] Failed to parse CSR: "+err.Error())
					return
				}

				//resp, err := http.Post("http://example.com/upload", "application/json", &buf)
				//var data interface{}
				pem, err := csr.ToPEM()
				if err != nil {
					fmt.Fprintln(os.Stderr, "[ERROR] Failed to parse CSR: "+err.Error())
					return
				}

				data := make(url.Values)
				data.Add("csr", string(pem))

				//resp, err := http.PostForm("http://"+host+"/certificates", data)
				client := &http.Client{}
				req, _ := http.NewRequest("POST", "http://"+host+"/certificates", strings.NewReader(data.Encode()))
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				if c.String("token") != "" {
					req.Header.Set("X-API-KEY", c.String("token"))
				}
				resp, err := client.Do(req)
				if err != nil {
					fmt.Fprintln(os.Stderr, "[ERROR] Failed to request: "+err.Error())
					return
				}
				decoder := json.NewDecoder(resp.Body)
				respData := make(map[string]map[string]interface{})
				if err := decoder.Decode(&respData); err != nil {
					panic(err)
				}

				if format == "cert" {
					fmt.Print(respData["certificate"]["crt"])
				}
				if format == "id" {
					fmt.Println(respData["certificate"]["id"])
				}
			},
		},

		{
			Name: "getcrt",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "host",
					Usage: "Host ip & port",
				},
			},
			Usage: "Get a Certificate from CA and output to STDOUT",
			Action: func(c *cli.Context) {

				host := c.String("host")
				if host == "" {
					fmt.Fprintln(os.Stderr, "[ERROR] Requere host as parameter")
					return
				}

				args := c.Args()
				if len(args) == 0 {
					fmt.Fprintln(os.Stderr, "[ERROR] Should provide id (same as serial number)")
					return
				}
				id := c.Args().First()

				resp, err := http.Get("http://" + host + "/certificates/" + id)
				if err != nil {
					fmt.Fprintln(os.Stderr, "[ERROR] Failed to request: "+err.Error())
					return
				}

				decoder := json.NewDecoder(resp.Body)
				respData := make(map[string]map[string]interface{})
				if err := decoder.Decode(&respData); err != nil {
					panic(err)
				}

				fmt.Print(respData["certificate"]["crt"])

			},
		},

		{
			Name: "getcacrt",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "host",
					Usage: "Host ip & port",
				},
			},
			Usage: "Get CA's Certificate and output to STDOUT",
			Action: func(c *cli.Context) {

				host := c.String("host")
				if host == "" {
					fmt.Fprintln(os.Stderr, "[ERROR] Require host as parameter")
					return
				}

				resp, err := http.Get("http://" + host + "/ca/certificate")
				if err != nil {
					fmt.Fprintln(os.Stderr, "[ERROR] Failed to request CA cert: "+err.Error())
					os.Exit(1)
				}

				decoder := json.NewDecoder(resp.Body)
				if resp.StatusCode != 200 {
					fmt.Fprintln(os.Stderr, "[ERROR] Failed to request CA cert: "+resp.Status)
					os.Exit(resp.StatusCode)
				}
				respData := make(map[string]map[string]interface{})
				if err := decoder.Decode(&respData); err != nil {
					panic(err)
				}

				fmt.Print(respData["ca"]["crt"])

			},
		},
		{
			Name: "oneshot",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "host",
					Usage: "Host ip & port",
				},
				cli.StringFlag{
					Name:  "key",
					Usage: "Path to private key file",
				},
				cli.StringFlag{
					Name:  "crt",
					Usage: "Path to public cert file",
				},
				cli.StringFlag{
					Name:  "ca",
					Usage: "Path to public ca file",
				},
				cli.StringFlag{
					Name:  "token",
					Usage: "Authorization Token",
				},
			},
			Usage: "Generate key, request, and submit to a server, all in one shot",
			Action: func(c *cli.Context) {
				keypath := c.String("key")
				if keypath == "" {
					fmt.Fprintln(os.Stderr, "Path to private key required")
					return
				}
				crtpath := c.String("crt")
				if crtpath == "" {
					fmt.Fprintln(os.Stderr, "Path to public certificate required")
					return
				}
				capath := c.String("ca")
				if capath == "" {
					fmt.Fprintln(os.Stderr, "Path to central authority required")
					return
				}
				host := c.String("host")
				if host == "" {
					fmt.Fprintln(os.Stderr, "Location of host required")
					return
				}
				keyfile, err := os.Create(keypath)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				defer keyfile.Close()

				crtfile, err := os.Create(crtpath)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				defer crtfile.Close()

				cafile, err := os.Create(capath)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				defer cafile.Close()

				// genkey
				key, err := pkix.NewKey()
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				pem, err := key.ToPEM()
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				n, err := io.WriteString(keyfile, string(pem))
				if err != nil {
					fmt.Println(n, err)
				}
				fmt.Println("Key written to", keypath)

				// gencsr
				csr, err := pkix.NewCertificateRequest(key)
				if err != nil {
					fmt.Fprintln(os.Stderr, "[ERROR] Failed to generate CSR: "+err.Error())
					return
				}
				csrpem, err := csr.ToPEM()
				// submitcsr
				data := make(url.Values)
				data.Add("csr", string(csrpem))

				//resp, err := http.PostForm("http://"+host+"/certificates", data)
				client := &http.Client{}
				req, _ := http.NewRequest("POST", "http://"+host+"/certificates", strings.NewReader(data.Encode()))
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				if c.String("token") != "" {
					req.Header.Set("X-API-KEY", c.String("token"))
				}
				resp, err := client.Do(req)
				if err != nil {
					fmt.Fprintln(os.Stderr, "[ERROR] Failed to request: "+err.Error())
					return
				}
				decoder := json.NewDecoder(resp.Body)
				respData := make(map[string]map[string]interface{})
				if err := decoder.Decode(&respData); err != nil {
					panic(err)
				}

				n, err = io.WriteString(crtfile, string(respData["certificate"]["crt"].(string)))
				if err != nil {
					fmt.Println(n, err)
				}
				fmt.Println("Certificate written to", crtpath)

				// getcacrt
				resp, err = http.Get("http://" + host + "/ca/certificate")
				if err != nil {
					fmt.Fprintln(os.Stderr, "[ERROR] Failed to request CA cert: "+err.Error())
					os.Exit(1)
				}

				decoder = json.NewDecoder(resp.Body)
				if resp.StatusCode != 200 {
					fmt.Fprintln(os.Stderr, "[ERROR] Failed to request CA cert: "+resp.Status)
					os.Exit(resp.StatusCode)
				}
				respData = make(map[string]map[string]interface{})
				if err := decoder.Decode(&respData); err != nil {
					panic(err)
				}

				n, err = io.WriteString(cafile, string(respData["ca"]["crt"].(string)))
				if err != nil {
					fmt.Println(n, err)
				}
				fmt.Println("Certificate Authority written to", capath)
			},
		},
	}

	app.Run(os.Args)

}
