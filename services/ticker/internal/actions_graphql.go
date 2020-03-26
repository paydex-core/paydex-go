package ticker

import (
	"github.com/paydex-core/paydex-go/services/ticker/internal/gql"
	"github.com/paydex-core/paydex-go/services/ticker/internal/tickerdb"
	hlog "github.com/paydex-core/paydex-go/support/log"
)

func StartGraphQLServer(s *tickerdb.TickerSession, l *hlog.Entry, port string) {
	graphql := gql.New(s, l)

	graphql.Serve(port)
}
