package pkix

import (
	"crypto/x509/pkix"
	"reflect"
	"testing"
)

func Test_Marshal(t *testing.T) {
	testdata := make(map[string]pkix.Name)

	testdata["O=null"] = pkix.Name{
		Organization: []string{"null"},
	}

	testdata["CN=myserver,O=myorg"] = pkix.Name{
		Organization: []string{"myorg"},
		CommonName:   "myserver",
	}
	testdata["CN=myserver,C=merica,L=City,ST=State,SA=StreetAddress,O=myorg,OU=myorgunit"] = pkix.Name{
		Organization:       []string{"myorg"},
		OrganizationalUnit: []string{"myorgunit"},
		Country:            []string{"merica"},
		Locality:           []string{"City"},
		Province:           []string{"State"},
		StreetAddress:      []string{"StreetAddress"},
		CommonName:         "myserver",
	}

	for expected, input := range testdata {
		output, err := Marshal(input)
		if err != nil {
			t.Error(err)
		}
		if output != expected {
			t.Error("Failed got " + output + " expected " + expected)
		}
	}
}

func Test_Unmarshal(t *testing.T) {
	testdata := make(map[string]pkix.Name)

	testdata["O=null"] = pkix.Name{
		Organization: []string{"null"},
	}

	testdata["CN=myserver,O=myorg"] = pkix.Name{
		Organization: []string{"myorg"},
		CommonName:   "myserver",
	}
	testdata["CN=myserver,C=merica,L=City,ST=State,SA=StreetAddress,O=myorg,OU=myorgunit"] = pkix.Name{
		Organization:       []string{"myorg"},
		OrganizationalUnit: []string{"myorgunit"},
		Country:            []string{"merica"},
		Locality:           []string{"City"},
		Province:           []string{"State"},
		StreetAddress:      []string{"StreetAddress"},
		CommonName:         "myserver",
	}

	for dnstring, expected := range testdata {
		pkixName, err := Unmarshal(dnstring)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(pkixName, expected) {
			t.Error("Failed got ", pkixName, " expected ", expected)
		}
	}
}
