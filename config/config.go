// Environment file for getting variables
// Currently the only thing it does is set the master password
// Should probably have it take over functions from OS such as port and mongodb connection details
// Reads from the config/environments/dev.yaml file by default
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// setting dev by default
func getEnv() string {
	env := os.Getenv("NETMAKER_ENV")
	if len(env) == 0 {
		return "dev"
	}
	return env
}

// Config : application config stored as global variable
var Config *EnvironmentConfig = &EnvironmentConfig{}
var SetupErr error

// EnvironmentConfig - environment conf struct
type EnvironmentConfig struct {
	Server ServerConfig `yaml:"server"`
	SQL    SQLConfig    `yaml:"sql"`
}

// ServerConfig - server conf struct
type ServerConfig struct {
	CoreDNSAddr          string `yaml:"corednsaddr"`
	APIConnString        string `yaml:"apiconn"`
	APIHost              string `yaml:"apihost"`
	APIPort              string `yaml:"apiport"`
	Broker               string `yam:"broker"`
	ServerBrokerEndpoint string `yaml:"serverbrokerendpoint"`
	BrokerType           string `yaml:"brokertype"`
	EmqxRestEndpoint     string `yaml:"emqxrestendpoint"`
	MasterKey            string `yaml:"masterkey"`
	DNSKey               string `yaml:"dnskey"`
	AllowedOrigin        string `yaml:"allowedorigin"`
	NodeID               string `yaml:"nodeid"`
	RestBackend          string `yaml:"restbackend"`
	MessageQueueBackend  string `yaml:"messagequeuebackend"`
	DNSMode              string `yaml:"dnsmode"`
	DisableRemoteIPCheck string `yaml:"disableremoteipcheck"`
	Version              string `yaml:"version"`
	SQLConn              string `yaml:"sqlconn"`
	Platform             string `yaml:"platform"`
	Database             string `yaml:"database"`
	Verbosity            int32  `yaml:"verbosity"`
	AuthProvider         string `yaml:"authprovider"`
	OIDCIssuer           string `yaml:"oidcissuer"`
	ClientID             string `yaml:"clientid"`
	ClientSecret         string `yaml:"clientsecret"`
	FrontendURL          string `yaml:"frontendurl"`
	DisplayKeys          string `yaml:"displaykeys"`
	AzureTenant          string `yaml:"azuretenant"`
	Telemetry            string `yaml:"telemetry"`
	HostNetwork          string `yaml:"hostnetwork"`
	Server               string `yaml:"server"`
	PublicIPService      string `yaml:"publicipservice"`
	MQPassword           string `yaml:"mqpassword"`
	MQUserName           string `yaml:"mqusername"`
	MetricsExporter      string `yaml:"metrics_exporter"`
	BasicAuth            string `yaml:"basic_auth"`
	LicenseValue         string `yaml:"license_value"`
	NetmakerAccountID    string `yaml:"netmaker_account_id"`
	IsEE                 string `yaml:"is_ee"`
	StunPort             int    `yaml:"stun_port"`
	StunHost             string `yaml:"stun_host"`
	Proxy                string `yaml:"proxy"`
	NodeLimit            int    `yaml:"node_limit"`
	UserLimit            int    `yaml:"user_limit"`
	ClientLimit          int    `yaml:"client_limit"`
	NetworkLimit         int    `yaml:"network_limit"`
	DeployedByOperator   bool   `yaml:"deployed_by_operator"`
}

// SQLConfig - Generic SQL Config
type SQLConfig struct {
	Host     string `yaml:"host"`
	Port     int32  `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DB       string `yaml:"db"`
	SSLMode  string `yaml:"sslmode"`
}

// reading in the env file
func ReadConfig(absolutePath string) (*EnvironmentConfig, error) {
	if len(absolutePath) == 0 {
		absolutePath = fmt.Sprintf("environments/%s.yaml", getEnv())
	}
	f, err := os.Open(absolutePath)
	var cfg EnvironmentConfig
	if err != nil {
		return &cfg, err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	if decoder.Decode(&cfg) != nil {
		return &cfg, err
	}
	return &cfg, err
}
