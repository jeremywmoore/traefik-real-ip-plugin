package traefik_real_ip_plugin

import (
	"github.com/stretchr/testify/assert"
	"net"
	"net/http"
	"testing"
)

func TestHeaderRetriever(t *testing.T) {
	testCases := []struct {
		description string
		header     []string
		want       string
	}{
		{
			"Basic case",
			[]string{"127.0.0.1"},
			"127.0.0.1",
		},
		{
			"Multiple header values",
			[]string{"127.0.0.1", "10.0.0.1, 10.0.0.2"},
			"127.0.0.1",
		},
		{
			"Invalid leading IP",
			[]string{"127.0.0.x", "25.0.0.1"},
			"25.0.0.1",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			header := &http.Header{"X-Forwarded-For": tc.header}
			retriever := &HeaderRetriever{"X-Forwarded-For"}
			retrieved := retriever.Retrieve(*header)
			assert.Equal(t, tc.want, retrieved.String())
		})
	}
}
  

func TestProxyCountRetriever(t *testing.T) {
	testCases := []struct {
		description string
		proxyCount int
		header     []string
		want       string
	}{
		{
			"Proxy count 0",
			0,
			[]string{"127.0.0.1, 25.0.0.1, 10.0.0.1, 10.0.0.2"},
			"",
		},
		{
			"Proxy count -1",
			-1,
			[]string{"127.0.0.1, 25.0.0.1, 10.0.0.1, 10.0.0.2"},
			"",
		},
		{
			"Proxy count 1",
			1,
			[]string{"127.0.0.1, 25.0.0.1, 10.0.0.1, 10.0.0.2"},
			"10.0.0.1",
		},
		{
			"Proxy count 1 with multiple header values",
			1,
			[]string{"127.0.0.1, 25.0.0.1", "10.0.0.1, 10.0.0.2"},
			"10.0.0.1",
		},
		{
			"Proxy count 2",
			2,
			[]string{"127.0.0.1, 25.0.0.1, 10.0.0.1, 10.0.0.2"},
			"25.0.0.1",
		},
		{
			"Proxy count 4",
			4,
			[]string{"127.0.0.1, 25.0.0.1, 10.0.0.1, 10.0.0.2"},
			"",
		},
		{
			"Proxy count 5",
			5,
			[]string{"127.0.0.1, 25.0.0.1, 10.0.0.1, 10.0.0.2"},
			"",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			header := &http.Header{"X-Forwarded-For": tc.header}
			retriever := &ProxyCountRetriever{"X-Forwarded-For", tc.proxyCount}
			retrieved := retriever.Retrieve(*header)
			if(tc.want == "") {
				assert.Nil(t, retrieved)
			} else {
				assert.Equal(t, tc.want, retrieved.String())
			}
		})
	}

}		

func TestProxyCIDRRetriever(t *testing.T) {
	testCases := []struct {
		description string
		proxyCIDRs []string
		header     []string
		want       string
	}{
		{
			"Multiple untrusted IPs",
			[]string{"10.0.0.0/30"},
			[]string{"127.0.0.1, 25.0.0.1, 10.0.0.1, 10.0.0.2"},
			"25.0.0.1",
		},
		{
			"Multiple untrusted IPs and multiple header values",
			[]string{"10.0.0.0/30"},
			[]string{"127.0.0.1, 25.0.0.1", "10.0.0.1, 10.0.0.2"},
			"25.0.0.1",
		},
		{
			"Empty Proxy CIDRs",
			[]string{},
			[]string{"127.0.0.1, 25.0.0.1, 10.0.0.1, 10.0.0.2"},
			"10.0.0.2",
		},
		{
			"Empty Proxy CIDRs and no header value",
			[]string{},
			[]string{},
			"",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			cidrs := make([]*net.IPNet, 0, len(tc.proxyCIDRs))
			for _, c := range tc.proxyCIDRs {
				_, cidr, err := net.ParseCIDR(c)
				if err != nil {
					assert.NoError(t, err)
				}
				cidrs = append(cidrs, cidr)
			}
			header := &http.Header{"X-Forwarded-For": tc.header}
			retriever := &ProxyCIDRRetriever{"X-Forwarded-For", cidrs}
			retrieved := retriever.Retrieve(*header)
			if(tc.want == "") {
				assert.Nil(t, retrieved)
			} else {
				assert.Equal(t, tc.want, retrieved.String())
			}
		})
	}
}
