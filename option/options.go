package option

type DB struct {
	Driver     string `default:"mysql"`
	DataSource string
	DBName     string `json:"dbname"`
	UserName   string `json:"username"`
	Password   string `json:"password"`
	Host       string `json:"host"`
	Port       int    `json:"port"`
	ReadHost   string `json:"readhost"`

	MaxIdleConns    int `json:"maxidleconns"`
	MaxOpenConns    int `json:"maxopenconns"`
	ConnMaxLifetime int `json:"connMaxLifetime"`

	ReadTimeout  int `json:"readtimeout"`
	WriteTimeout int `json:"writetimeout"`
	Timeout      int `json:"timeout"`
}

type Postgresql struct {
}

type Options struct {
	// host:port address.
	Addr string
	// Optional password. Must match the password specified in the
	// requirepass server configuration option.
	Password string
	// Database to be selected after connecting to the server.
	DB int
	// Dial timeout for establishing new connections.
	// Default is 5 seconds.
	DialTimeout int
	// Amount of time client waits for connection if all connections
	// are busy before returning an error.
	// Default is ReadTimeout + 1 second.
	PoolTimeout int
	// Maximum number of socket connections.
	// Default is 10 connections per every CPU as reported by runtime.NumCPU.
	PoolSize int
	// Amount of time after which client closes idle connections.
	// Should be less than server's timeout.
	// Default is 5 minutes. -1 disables idle timeout check.
	IdleTimeout int64
	// Timeout for socket reads. If reached, commands will fail
	// with a timeout instead of blocking. Use value -1 for no timeout and 0 for default.
	// Default is 3 seconds.
	ReadTimeout int
	// Timeout for socket writes. If reached, commands will fail
	// with a timeout instead of blocking.
	// Default is ReadTimeout.
	WriteTimeout int
	// Maximum number of retries before giving up.
	// Default is to not retry failed commands.
	MaxRetries int
	//
	IsCluster bool
}
