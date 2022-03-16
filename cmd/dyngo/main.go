package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	externalip "github.com/glendc/go-external-ip"
	"github.com/nrdcg/goinwx"
)

func get_env_value(env_var string) string {
	ret := os.Getenv(env_var)
	if ret != "" {
		return ret
	}
	fmt.Printf("Error! Did not find required environment variable %v\n", env_var)
	os.Exit(1)
	return ""
}

func determineExtIP() net.IP {
	consensus := externalip.DefaultConsensus(nil, nil)
	consensus.UseIPProtocol(4)

	ip, err := consensus.ExternalIP()
	if err != nil {
		fmt.Println("Error retrieving external IP!")
		os.Exit(1)
	}
	log.Printf("Found external IP: %v\n", ip.To4().String())
	return ip
}

func removeSubdomain(fullTLD string) string {
	domain := strings.Split(fullTLD, ".")
	ret := domain[1] + "." + domain[2]
	return ret
}

func getSubdomainRecord(subdomain string, client *goinwx.Client) (goinwx.NameserverRecord, error) {
	var req = &goinwx.NameserverInfoRequest{
		Domain: removeSubdomain(subdomain),
	}

	resp, _ := client.Nameservers.Info(req) // FIXME: Error handling
	for _, child := range resp.Records {
		if child.Name == subdomain {
			log.Printf("Found Subdomain: %v\n", child)
			return child, nil
		}
	}
	return *&goinwx.NameserverRecord{}, errors.New("Did not find subdomain!")
}

func updateRecord(user string, pass string, record string, address string) {
	client := goinwx.NewClient(user, pass, &goinwx.ClientOptions{})
	_, err := client.Account.Login()
	if err != nil {
		log.Fatalf("INWX: Failed login:\n%v\n", err)
	}
	defer func() {
		if err := client.Account.Logout(); err != nil {
			log.Fatalf("INWX: Failed logout:\n%v\n", err)
		}
	}()

	subdomain, err := getSubdomainRecord(record, client)
	if err != nil {
		log.Fatalf("INWX: Error while finding subdomain: %v\n", err)
		os.Exit(1)
	}

	var request = &goinwx.NameserverRecordRequest{
		Name:    subdomain.Name,
		Type:    subdomain.Type,
		Content: address,
		TTL:     301,
	}

	err = client.Nameservers.UpdateRecord(subdomain.ID, request)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	username := get_env_value("INWX_USERNAME")
	password := get_env_value("INWX_PASSWORD")
	domain_record := get_env_value("INWX_DOMAIN_RECORD")
	ip_v4 := determineExtIP().To4().String()

	log.Printf("Found data: %v for %v - IP: %v\n", username, domain_record, ip_v4)

	updateRecord(username, password, domain_record, ip_v4)
	log.Printf("Updated record for %v to %v", domain_record, ip_v4)
}
