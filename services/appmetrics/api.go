package appmetrics

import (
	"context"

	"github.com/pborman/uuid"

	"github.com/status-im/status-go/appmetrics"
	"github.com/status-im/status-go/logutils"
)

func NewAPI(db *appmetrics.Database) *API {
	return &API{db: db, sessionID: uuid.NewRandom().String()}
}

type API struct {
	db        *appmetrics.Database
	sessionID string
}

func (api *API) ValidateAppMetrics(ctx context.Context, appMetrics []appmetrics.AppMetric) error {
	logutils.ZapLogger().Debug("[AppMetricsAPI::ValidateAppMetrics]")
	return api.db.ValidateAppMetrics(appMetrics)
}

func (api *API) SaveAppMetrics(ctx context.Context, appMetrics []appmetrics.AppMetric) error {
	logutils.ZapLogger().Debug("[AppMetricsAPI::SaveAppMetrics]")
	return api.db.SaveAppMetrics(appMetrics, api.sessionID)
}

func (api *API) GetAppMetrics(ctx context.Context, limit int, offset int) (appmetrics.Page, error) {
	logutils.ZapLogger().Debug("[AppMetricsAPI::GetAppMetrics]")
	return api.db.GetAppMetrics(limit, offset)
}
