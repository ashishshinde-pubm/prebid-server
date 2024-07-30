package katana

import (
	"encoding/json"

	"github.com/prebid/prebid-server/v2/modules/moduledeps"
	"github.com/prebid/prebid-server/v2/modules/prebid/katana/config"
)

type Katana struct {
	cfg config.Config
}

func InitKatana(rawCfg json.RawMessage, moduleDeps moduledeps.ModuleDeps) (Katana, error) {
	cfg := config.Config{}

	err := json.Unmarshal(rawCfg, &cfg)
	if err != nil {
		return Katana{}, err
	}

	kt := Katana{
		cfg: cfg,
	}

	return kt, nil

}
