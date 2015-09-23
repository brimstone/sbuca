package pkix

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"math/big"
	"strings"
)

type Certificate struct {
	DerBytes []byte

	Crt *x509.Certificate
}

func GenSubject(organization string) pkix.Name {
	return pkix.Name{
		Organization: []string{organization},
	}
}
func NewCertificateFromDER(derBytes []byte) (*Certificate, error) {

	crt, err := x509.ParseCertificate(derBytes)
	if err != nil {
		return nil, err
	}

	cert := &Certificate{
		DerBytes: derBytes,
		Crt:      crt,
	}

	return cert, nil
}
func NewCertificateFromPEM(pemBytes []byte) (*Certificate, error) {

	pemBlock, _ := pem.Decode(pemBytes)
	if pemBlock == nil {
		return nil, errors.New("PEM decode failed")
	}

	crt, err := x509.ParseCertificate(pemBlock.Bytes)
	if err != nil {
		return nil, err
	}

	cert := &Certificate{
		DerBytes: pemBlock.Bytes,
		Crt:      crt,
	}

	return cert, nil
}
func NewCertificateFromPEMFile(filename string) (*Certificate, error) {

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return NewCertificateFromPEM(data)
}

func (certificate *Certificate) ToPEM() ([]byte, error) {

	pemBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certificate.DerBytes,
	}

	pemBytes := pem.EncodeToMemory(pemBlock)

	return pemBytes, nil
}
func (certificate *Certificate) ToPEMFile(filename string) error {
	pemBytes, err := certificate.ToPEM()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, pemBytes, 0400)
}
func (certificate *Certificate) GetSerialNumber() *big.Int {
	return certificate.Crt.SerialNumber
}

func Marshal(name pkix.Name) (string, error) {
	var output []string
	if name.CommonName != "" {
		output = append(output, "CN="+name.CommonName)
	}
	if len(name.Country) > 0 {
		for i := range name.Country {
			output = append(output, "C="+name.Country[i])
		}
	}
	if len(name.Locality) > 0 {
		for i := range name.Locality {
			output = append(output, "L="+name.Locality[i])
		}
	}
	if len(name.Province) > 0 {
		for i := range name.Province {
			output = append(output, "ST="+name.Province[i])
		}
	}
	if len(name.StreetAddress) > 0 {
		for i := range name.StreetAddress {
			output = append(output, "SA="+name.StreetAddress[i])
		}
	}
	if len(name.Organization) > 0 {
		for i := range name.Organization {
			output = append(output, "O="+name.Organization[i])
		}
	}
	if len(name.OrganizationalUnit) > 0 {
		for i := range name.OrganizationalUnit {
			output = append(output, "OU="+name.OrganizationalUnit[i])
		}
	}
	return strings.Join(output, ","), nil
}

func Unmarshal(dn string) (pkix.Name, error) {
	var output pkix.Name
	segments := strings.Split(dn, ",")
	for segment := range segments {
		identifier := strings.SplitN(segments[segment], "=", 2)
		if identifier[0] == "CN" {
			output.CommonName = identifier[1]
		} else if identifier[0] == "C" {
			output.Country = append(output.Country, identifier[1])
		} else if identifier[0] == "L" {
			output.Locality = append(output.Locality, identifier[1])
		} else if identifier[0] == "ST" {
			output.Province = append(output.Province, identifier[1])
		} else if identifier[0] == "SA" {
			output.StreetAddress = append(output.StreetAddress, identifier[1])
		} else if identifier[0] == "O" {
			output.Organization = append(output.Organization, identifier[1])
		} else if identifier[0] == "OU" {
			output.OrganizationalUnit = append(output.OrganizationalUnit, identifier[1])
		}
	}
	return output, nil
}
