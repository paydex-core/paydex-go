package horizon

import (
	"compress/flate"
	"context"
	"database/sql"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	chimiddleware "github.com/go-chi/chi/middleware"
	metrics "github.com/rcrowley/go-metrics"
	"github.com/rs/cors"
	"github.com/sebest/xff"

	"github.com/paydex-core/go-throttled"
	"github.com/paydex-core/paydex-go/exp/orderbook"
	"github.com/paydex-core/paydex-go/services/horizon/internal/actions"
	"github.com/paydex-core/paydex-go/services/horizon/internal/db2"
	"github.com/paydex-core/paydex-go/services/horizon/internal/db2/core"
	"github.com/paydex-core/paydex-go/services/horizon/internal/db2/history"
	"github.com/paydex-core/paydex-go/services/horizon/internal/ledger"
	"github.com/paydex-core/paydex-go/services/horizon/internal/paths"
	hProblem "github.com/paydex-core/paydex-go/services/horizon/internal/render/problem"
	"github.com/paydex-core/paydex-go/services/horizon/internal/render/sse"
	"github.com/paydex-core/paydex-go/services/horizon/internal/txsub/sequence"
	"github.com/paydex-core/paydex-go/support/db"
	"github.com/paydex-core/paydex-go/support/log"
	"github.com/paydex-core/paydex-go/support/render/problem"
)

const (
	LRUCacheSize            = 50000
	maxAssetsForPathFinding = 15
)

// Web contains the http server related fields for horizon: the router,
// rate limiter, etc.
type web struct {
	appCtx             context.Context
	router             *chi.Mux
	rateLimiter        *throttled.HTTPRateLimiter
	sseUpdateFrequency time.Duration
	staleThreshold     uint
	ingestFailedTx     bool

	historyQ *history.Q
	coreQ    *core.Q

	requestTimer metrics.Timer
	failureMeter metrics.Meter
	successMeter metrics.Meter
}

func init() {
	// register problems
	problem.RegisterError(sql.ErrNoRows, problem.NotFound)
	problem.RegisterError(sequence.ErrNoMoreRoom, hProblem.ServerOverCapacity)
	problem.RegisterError(db2.ErrInvalidCursor, problem.BadRequest)
	problem.RegisterError(db2.ErrInvalidLimit, problem.BadRequest)
	problem.RegisterError(db2.ErrInvalidOrder, problem.BadRequest)
	problem.RegisterError(sse.ErrRateLimited, hProblem.RateLimitExceeded)
	problem.RegisterError(context.DeadlineExceeded, hProblem.Timeout)
	problem.RegisterError(context.Canceled, hProblem.ServiceUnavailable)
	problem.RegisterError(db.ErrCancelled, hProblem.ServiceUnavailable)
}

// mustInitWeb installed a new Web instance onto the provided app object.
func mustInitWeb(ctx context.Context, hq *history.Q, cq *core.Q, updateFreq time.Duration, threshold uint, ingestFailedTx bool) *web {
	if hq == nil {
		log.Fatal("missing history DB for installing the web instance")
	}
	if cq == nil {
		log.Fatal("missing core DB for installing the web instance")
	}

	return &web{
		appCtx:             ctx,
		router:             chi.NewRouter(),
		historyQ:           hq,
		coreQ:              cq,
		sseUpdateFrequency: updateFreq,
		staleThreshold:     threshold,
		ingestFailedTx:     ingestFailedTx,
		requestTimer:       metrics.NewTimer(),
		failureMeter:       metrics.NewMeter(),
		successMeter:       metrics.NewMeter(),
	}
}

// mustInstallMiddlewares installs the middleware stack used for horizon onto the
// provided app.
// Note that a request will go through the middlewares from top to bottom.
func (w *web) mustInstallMiddlewares(app *App, connTimeout time.Duration) {
	if w == nil {
		log.Fatal("missing web instance for installing middlewares")
	}

	r := w.router
	r.Use(chimiddleware.StripSlashes)

	//TODO: remove this middleware
	r.Use(appContextMiddleware(app))

	r.Use(requestCacheHeadersMiddleware)
	r.Use(chimiddleware.RequestID)
	r.Use(contextMiddleware)
	r.Use(xff.Handler)
	r.Use(loggerMiddleware)
	r.Use(timeoutMiddleware(connTimeout))
	r.Use(requestMetricsMiddleware)
	r.Use(recoverMiddleware)
	r.Use(chimiddleware.Compress(flate.DefaultCompression, "application/hal+json"))

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"*"},
	})
	r.Use(c.Handler)

	r.Use(w.RateLimitMiddleware)
}

func installPathFindingRoutes(
	findPaths FindPathsHandler,
	findFixedPaths FindFixedPathsHandler,
	r *chi.Mux,
	expIngest bool,
	requiresExperimentalIngestion *ExperimentalIngestionMiddleware,
) {
	r.Group(func(r chi.Router) {
		r.Use(acceptOnlyJSON)
		if expIngest {
			r.Use(requiresExperimentalIngestion.Wrap)
		}
		r.Method("GET", "/paths", findPaths)
		r.Method("GET", "/paths/strict-receive", findPaths)
		r.Method("GET", "/paths/strict-send", findFixedPaths)
	})
}

func installAccountOfferRoute(
	offersAction actions.GetAccountOffersHandler,
	streamHandler sse.StreamHandler,
	enableExperimentalIngestion bool,
	r *chi.Mux,
	requiresExperimentalIngestion *ExperimentalIngestionMiddleware,
) {
	path := "/accounts/{account_id}/offers"
	if enableExperimentalIngestion {
		r.With(requiresExperimentalIngestion.Wrap).Method(
			http.MethodGet,
			path,
			streamablePageHandler(offersAction, streamHandler),
		)
	} else {
		r.Get(path, OffersByAccountAction{}.Handle)
	}
}

// mustInstallActions installs the routing configuration of horizon onto the
// provided app.  All route registration should be implemented here.
func (w *web) mustInstallActions(
	config Config,
	pathFinder paths.Finder,
	orderBookGraph *orderbook.OrderBookGraph,
	requiresExperimentalIngestion *ExperimentalIngestionMiddleware,
) {
	if w == nil {
		log.Fatal("missing web instance for installing web actions")
	}

	r := w.router
	r.Get("/", RootAction{}.Handle)
	r.Get("/metrics", MetricsAction{}.Handle)

	// ledger actions
	r.Route("/ledgers", func(r chi.Router) {
		r.Get("/", LedgerIndexAction{}.Handle)
		r.Route("/{ledger_id}", func(r chi.Router) {
			r.Get("/", LedgerShowAction{}.Handle)
			r.Get("/transactions", w.streamIndexActionHandler(w.getTransactionPage, w.streamTransactions))
			r.Get("/operations", OperationIndexAction{}.Handle)
			r.Get("/payments", OperationIndexAction{OnlyPayments: true}.Handle)
			r.Get("/effects", EffectIndexAction{}.Handle)
		})
	})

	// account actions
	r.Route("/accounts", func(r chi.Router) {
		r.With(requiresExperimentalIngestion.Wrap).
			Method(
				http.MethodGet,
				"/",
				restPageHandler(actions.GetAccountsHandler{}),
			)
		r.Route("/{account_id}", func(r chi.Router) {
			r.Get("/", w.streamShowActionHandler(w.getAccountInfo, true))
			r.Get("/transactions", w.streamIndexActionHandler(w.getTransactionPage, w.streamTransactions))
			r.Get("/operations", OperationIndexAction{}.Handle)
			r.Get("/payments", OperationIndexAction{OnlyPayments: true}.Handle)
			r.Get("/effects", EffectIndexAction{}.Handle)
			r.Get("/trades", TradeIndexAction{}.Handle)
			r.Get("/data/{key}", DataShowAction{}.Handle)
		})
	})

	streamHandler := sse.StreamHandler{
		RateLimiter:  w.rateLimiter,
		LedgerSource: ledger.NewHistoryDBSource(w.sseUpdateFrequency),
	}

	installAccountOfferRoute(
		actions.GetAccountOffersHandler{},
		streamHandler,
		config.EnableExperimentalIngestion,
		r,
		requiresExperimentalIngestion,
	)

	// transaction history actions
	r.Route("/transactions", func(r chi.Router) {
		r.Get("/", w.streamIndexActionHandler(w.getTransactionPage, w.streamTransactions))
		r.Route("/{tx_id}", func(r chi.Router) {
			r.Get("/", showActionHandler(w.getTransactionResource))
			r.Get("/operations", OperationIndexAction{}.Handle)
			r.Get("/payments", OperationIndexAction{OnlyPayments: true}.Handle)
			r.Get("/effects", EffectIndexAction{}.Handle)
		})
	})

	// operation actions
	r.Route("/operations", func(r chi.Router) {
		r.Get("/", OperationIndexAction{}.Handle)
		r.Get("/{id}", OperationShowAction{}.Handle)
		r.Get("/{op_id}/effects", EffectIndexAction{}.Handle)
	})

	// payment actions
	r.Get("/payments", OperationIndexAction{OnlyPayments: true}.Handle)

	// effect actions
	r.Get("/effects", EffectIndexAction{}.Handle)

	// trading related endpoints
	r.Get("/trades", TradeIndexAction{}.Handle)
	r.Get("/trade_aggregations", TradeAggregateIndexAction{}.Handle)

	r.Route("/offers", func(r chi.Router) {
		r.With(requiresExperimentalIngestion.Wrap).
			Method(
				http.MethodGet,
				"/",
				restPageHandler(actions.GetOffersHandler{}),
			)
		r.With(acceptOnlyJSON, requiresExperimentalIngestion.Wrap).
			Method(
				http.MethodGet,
				"/{id}",
				objectActionHandler{actions.GetOfferByID{}},
			)
		r.Get("/{offer_id}/trades", TradeIndexAction{}.Handle)
	})

	if config.EnableExperimentalIngestion {
		r.With(requiresExperimentalIngestion.Wrap).Method(
			http.MethodGet,
			"/order_book",
			streamableObjectActionHandler{
				streamHandler: streamHandler,
				action: actions.GetOrderbookHandler{
					OrderBookGraph: orderBookGraph,
				},
			},
		)
	} else {
		r.Get("/order_book", OrderBookShowAction{}.Handle)
	}

	// Transaction submission API
	r.Post("/transactions", TransactionCreateAction{}.Handle)

	findPaths := FindPathsHandler{
		staleThreshold:       config.StaleThreshold,
		checkHistoryIsStale:  !config.EnableExperimentalIngestion,
		setLastLedgerHeader:  config.EnableExperimentalIngestion,
		maxPathLength:        config.MaxPathLength,
		maxAssetsParamLength: maxAssetsForPathFinding,
		pathFinder:           pathFinder,
		coreQ:                w.coreQ,
	}
	findFixedPaths := FindFixedPathsHandler{
		maxPathLength:        config.MaxPathLength,
		setLastLedgerHeader:  config.EnableExperimentalIngestion,
		maxAssetsParamLength: maxAssetsForPathFinding,
		pathFinder:           pathFinder,
		coreQ:                w.coreQ,
	}
	installPathFindingRoutes(
		findPaths,
		findFixedPaths,
		w.router,
		config.EnableExperimentalIngestion,
		requiresExperimentalIngestion,
	)

	if config.EnableExperimentalIngestion {
		r.With(requiresExperimentalIngestion.Wrap).Method(
			http.MethodGet,
			"/assets",
			restPageHandler(actions.AssetStatsHandler{}),
		)
	} else if config.EnableAssetStats {
		r.Get("/assets", AssetsAction{}.Handle)
	}

	// Network state related endpoints
	r.Get("/fee_stats", FeeStatsAction{}.Handle)

	// friendbot
	if config.FriendbotURL != nil {
		redirectFriendbot := func(w http.ResponseWriter, r *http.Request) {
			redirectURL := config.FriendbotURL.String() + "?" + r.URL.RawQuery
			http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
		}
		r.Post("/friendbot", redirectFriendbot)
		r.Get("/friendbot", redirectFriendbot)
	}

	r.NotFound(NotFoundAction{}.Handle)
}

func maybeInitWebRateLimiter(rateQuota *throttled.RateQuota) *throttled.HTTPRateLimiter {
	// Disabled
	if rateQuota == nil {
		return nil
	}

	rateLimiter, err := throttled.NewGCRARateLimiter(LRUCacheSize, *rateQuota)
	if err != nil {
		log.Fatalf("unable to create RateLimiter: %v", err)
	}

	return &throttled.HTTPRateLimiter{
		RateLimiter:   rateLimiter,
		DeniedHandler: &RateLimitExceededAction{Action{}},
		VaryBy:        VaryByRemoteIP{},
	}
}

type VaryByRemoteIP struct{}

func (v VaryByRemoteIP) Key(r *http.Request) string {
	return remoteAddrIP(r)
}

func remoteAddrIP(r *http.Request) string {
	// To support IPv6
	lastSemicolon := strings.LastIndex(r.RemoteAddr, ":")
	if lastSemicolon == -1 {
		return r.RemoteAddr
	} else {
		return r.RemoteAddr[0:lastSemicolon]
	}
}

// horizonSession returns a new session that loads data from the horizon
// database. The returned session is bound to `ctx`.
func (w *web) horizonSession(ctx context.Context) (*db.Session, error) {
	err := errorIfHistoryIsStale(w.isHistoryStale())
	if err != nil {
		return nil, err
	}

	return &db.Session{DB: w.historyQ.Session.DB, Ctx: ctx}, nil
}

// coreSession returns a new session that loads data from the paydex core
// database. The returned session is bound to `ctx`.
func (w *web) coreSession(ctx context.Context) *db.Session {
	return &db.Session{DB: w.coreQ.Session.DB, Ctx: ctx}
}

// isHistoryStale returns true if the latest history ledger is more than
// `StaleThreshold` ledgers behind the latest core ledger
func (w *web) isHistoryStale() bool {
	if w.staleThreshold == 0 {
		return false
	}

	ls := ledger.CurrentState()
	return (ls.CoreLatest - ls.HistoryLatest) > int32(w.staleThreshold)
}

// errorIfHistoryIsStale returns a formatted error if isStale is true.
func errorIfHistoryIsStale(isStale bool) error {
	if !isStale {
		return nil
	}

	ls := ledger.CurrentState()
	err := hProblem.StaleHistory
	err.Extras = map[string]interface{}{
		"history_latest_ledger": ls.HistoryLatest,
		"core_latest_ledger":    ls.CoreLatest,
	}
	return err
}
