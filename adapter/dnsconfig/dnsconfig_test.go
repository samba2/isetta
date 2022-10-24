package dnsconfig

import (
	"math/rand"
	"os"
	"path"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/ini.v1"
)

func TestWslConfCreationIfNotExists(t *testing.T) {
	tmpFileName := buildTmpFileName()
	defer os.Remove(tmpFileName)

	disableResolvConfGenerationForFile(tmpFileName)

	content, err := os.ReadFile(tmpFileName)
	assert.NoError(t, err)
	assert.Regexp(t, "generateResolvConf = false", string(content))
}

func TestDisableResolvConfGenerationForFile(t *testing.T) {
	file, err := os.CreateTemp("", "isetta")
	assert.NoError(t, err)
	defer os.Remove(file.Name())

	file.WriteString(`
	[miau]
	foo  = bar
	`)

	disableResolvConfGenerationForFile(file.Name())

	content, _ := os.ReadFile(file.Name())
	assert.Regexp(t, "generateResolvConf = false", string(content))
}

func TestSetToFalse(t *testing.T) {
	testCases := []struct {
		name          string
		contentBefore string
	}{
		{
			name: "already false",
			contentBefore: `
			[network]
			generateResolvConf  = false
			`,
		},
		{
			name: "currently true",
			contentBefore: `
			[network]
			generateResolvConf  = true
			`,
		},
		{
			name: "missing entry",
			contentBefore: `
			[miau]
			foo  = bar
			`,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			cfg, _ := ini.Load([]byte(tC.contentBefore))

			setToFalse(cfg, "foo")

			result, err := cfg.Section("network").Key("generateResolvConf").Bool()
			assert.NoError(t, err)
			assert.False(t, result)
		})
	}
}

func TestSetSetDnsServer(t *testing.T) {
	tmpFileName := buildTmpFileName()
	defer os.Remove(tmpFileName)

	setServer(tmpFileName, "1.2.3.4")

	content, err := os.ReadFile(tmpFileName)
	assert.NoError(t, err)
	assert.Contains(t, string(content), "nameserver 1.2.3.4\n")
}

func TestSetInvalidAddress(t *testing.T) {
	tmpFileName := buildTmpFileName()
	defer os.Remove(tmpFileName)

	err := isIpValid("invalid")
	assert.Error(t, err)
}

func TestIsDnsServerSet(t *testing.T) {
	tmpFileName := buildTmpFileName()
	defer os.Remove(tmpFileName)

	testCases := []struct {
		content       string
		searchAddress string
		desc          string
		result        bool
	}{
		{
			desc:          "valid name server 1",
			content:       "nameserver 8.8.8.8\n",
			searchAddress: "8.8.8.8",
			result:        true,
		},
		{
			desc:          "valid name server 2",
			content:       "#foo\nnameserver 8.8.8.8",
			searchAddress: "8.8.8.8",
			result:        true,
		},
		{
			desc:          "invalid name server",
			content:       "foo",
			searchAddress: "8.8.8.8",
			result:        false,
		},
		{
			desc:          "disabled name server",
			content:       "#nameserver 8.8.8.8\n",
			searchAddress: "8.8.8.8",
			result:        false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			os.WriteFile(tmpFileName, []byte(tC.content), 0644)
			assert.Equal(t, tC.result, isDnsServerSet(tC.searchAddress, tmpFileName))
		})
	}

}

func TestIsDnsServerSetErrorReadingFile(t *testing.T) {
	_, err := readResolveConf("/not/existing/path")
	assert.Error(t, err)
}

func buildTmpFileName() string {
	return path.Join(os.TempDir(), "isetta-"+strconv.Itoa(rand.Int()))
}
