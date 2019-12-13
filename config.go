package higgs

type (
	Configuration struct {
		Database DatabaseConfig
		Web      HttpConfig
		App      AppConfig
	}

	DatabaseConfig struct {
		URI        string
		Database   string
	}

	HttpConfig struct {
		UserAgent  string
		TimeoutSec int
	}

	AppConfig struct {
		MaxRoutines int
	}
)
