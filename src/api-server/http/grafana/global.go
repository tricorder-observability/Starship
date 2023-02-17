package grafana

// TODO(zhihui): Consider put these into a new type
// type Client struct {
//     ... base url, createAuthKey URI etc.
// }
var (
	BaseURL            string
	CreateAuthKeysURI  string
	CreateDashBoardURI string
	CreateDatabaseURI  string
	GetDashboardURI    string
	BasicAuth          string
)

// InitGrafanaConfig initializes global configurations for how to connect with Grafana.
func InitGrafanaConfig(baseURL, userName, password string) {
	// userName password come from command line flags
	BasicAuth = userName + ":" + password
	BaseURL = baseURL
	CreateAuthKeysURI = BaseURL + "/api/auth/keys"
	CreateDashBoardURI = BaseURL + "/api/dashboards/db"
	CreateDatabaseURI = BaseURL + "/api/datasources"
	GetDashboardURI = BaseURL + "/api/dashboards/uid/"
}
