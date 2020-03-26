package test

import (
	"github.com/paydex-core/paydex-go/services/horizon/internal/test/scenarios"
)

func loadScenario(scenarioName string, includeHorizon bool) {
	paydexCorePath := scenarioName + "-core.sql"

	scenarios.Load(PaydexCoreDatabaseURL(), paydexCorePath)

	if includeHorizon {
		horizonPath := scenarioName + "-horizon.sql"
		scenarios.Load(DatabaseURL(), horizonPath)
	}
}
