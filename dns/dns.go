package dns

import (
	"errors"
	"net"
	"strings"
)

type DnsResolver interface {
	GetSPFRecord(string) (string, error)
	GetARecords(string) ([]string, error)
	GetMXRecords(string) ([]*net.MX, error)
}

type GoSPFDNS struct {
}

func IsSPF(record string) bool {
	if strings.HasPrefix(record, "v=spf") {
		return true
	}

	return false
}

func IsSupportedProtocol(record string) bool {
	if strings.HasPrefix(record, "v=spf1") {
		return true
	}

	return false
}

func (dns *GoSPFDNS) GetARecords(name string) ([]string, error) {
	return net.LookupHost(name)
}

func (dns *GoSPFDNS) GetMXRecords(name string) ([]*net.MX, error) {
	return net.LookupMX(name)
}

func (dns *GoSPFDNS) GetSPFRecord(name string) (string, error) {

	records, err := net.LookupTXT(name)
	if err != nil {
		return "", err
	}

	for _, record := range records {
		if !IsSPF(record) {
			continue
		}

		if !IsSupportedProtocol(record) {
			return "", errors.New("Unsupported SPF record: " + record)
		}

		return record, nil

	}

	return "", errors.New("No SPF record found for " + name)

}
