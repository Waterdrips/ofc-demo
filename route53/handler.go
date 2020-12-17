package function

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/route53"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/openfaas/openfaas-cloud/sdk"
)

func Handle(w http.ResponseWriter, r *http.Request) {

	readSecrets()
	sess := session.Must(session.NewSession())
	client := route53.New(sess)
	z, err := client.ListHostedZones(&route53.ListHostedZonesInput{})
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusOK)
		return
	}
	var records []*route53.ResourceRecordSet
	for _, zone := range z.HostedZones {

		r, err := client.ListResourceRecordSets(&route53.ListResourceRecordSetsInput{HostedZoneId: zone.Id})
		if err != nil {
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusOK)
			return
		}
		for _, s := range r.ResourceRecordSets {
			records = append(records, s)
		}
	}

	s := ""

	for _, r := range records {
		s = fmt.Sprintf("%s\n%s", s, r.String())
	}
	w.Write([]byte(s))
}

func readSecrets() {
	id, _ := sdk.ReadSecret("access-key-id")
	key, _ := sdk.ReadSecret("secret-access-key")

	os.Setenv("AWS_ACCESS_KEY_ID", id)
	os.Setenv("AWS_SECRET_ACCESS_KEY", key)

}
