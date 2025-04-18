package client

import (
	"net"
	"testing"

	"github.com/console-dns/spec/models"
)

func TestAll(t *testing.T) {
	client := NewConsoleDnsClient("http://127.0.0.1:8090", "ce50401e-1924-4d5e-be7b-a8af8b147f6d")
	zones, _, err := client.ListZones()
	if err != nil {
		t.Fatal(err)
	}
	listZones := zones.ListZones()
	firstZone := listZones[0]
	_, _, err = client.ListZone(firstZone)
	if err != nil {
		t.Fatal(err)
	}
	a := models.RecordA{
		Ttl: 1200,
		Ip:  net.ParseIP("10.0.1.1"),
	}
	_, err = client.CreateRecord(firstZone, "www", "A", a)
	if err != nil {
		t.Fatal(err)
	}
	b := a.Clone()
	b.Ip = net.ParseIP("10.0.1.2")
	_, err = client.UpdateRecord(firstZone, "www", "A", a, b)
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.DeleteRecord(firstZone, "www", "A", b)
	if err != nil {
		t.Fatal(err)
	}
}
