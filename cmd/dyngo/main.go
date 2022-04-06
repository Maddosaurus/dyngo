package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

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
	if len(domain) < 3 {
		// we don't have a SUBdomain, we have a domain
		return fullTLD
	}
	ret := domain[1] + "." + domain[2]
	return ret
}

func getSubdomainRecord(subdomain string, client *goinwx.Client) (*goinwx.NameserverRecord, error) {
	var req = &goinwx.NameserverInfoRequest{
		Domain: removeSubdomain(subdomain),
	}

	resp, _ := client.Nameservers.Info(req) // FIXME: Error handling
	for _, child := range resp.Records {
		if child.Name == subdomain {
			log.Printf("Found Subdomain: %v\n", child)
			return &child, nil
		}
	}
	return nil, errors.New("Did not find subdomain!")
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

	if strings.Compare(subdomain.Content, address) == 0 {
		log.Printf("A record is still up to date - exiting!")
		return
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
	log.Printf("Updated record for %v to %v", record, address)
}

func main() {
	username := get_env_value("INWX_USERNAME")
	password := get_env_value("INWX_PASSWORD")
	domain_record := get_env_value("INWX_DOMAIN_RECORD")
	sleep_min_str := get_env_value("INWX_SLEEP_MINUTES")
	ip_v4 := ""

	sleep_min, err := strconv.Atoi(sleep_min_str)
	if err != nil {
		log.Fatalf("Could not convert INWX_SLEEP_MINUTES to integer: %v", err)
	}

	log.Printf("Running with user %v for %v\n", username, domain_record)

	for {
		ip_v4 = determineExtIP().To4().String()
		updateRecord(username, password, domain_record, ip_v4)
		time.Sleep(time.Duration(sleep_min) * time.Minute)
	}
}
