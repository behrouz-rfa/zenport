//go:build e2e

package e2e

import (
	"context"
	"database/sql"
	"zenport/gates/gatesclient"

	"github.com/cucumber/godog"
	"github.com/go-openapi/strfmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"zenport/gates/gatesclient/models"
	"zenport/gates/gatesclient/time"
)

type timeKey struct{}

type gatesFeature struct {
	db     *sql.DB
	client *gatesclient.TimeProcessing
}

var _ feature = (*gatesFeature)(nil)

func (c *gatesFeature) init(cfg featureConfig) (err error) {
	if cfg.useMonoDB {
		c.db, err = sql.Open("pgx", "postgres://zenports_user:zenports_pass@localhost:5432/zenports?sslmode=disable")
	} else {
		c.db, err = sql.Open("pgx", "postgres://gates_user:gates_pass@localhost:5432/gates?sslmode=disable&search_path=gates,public")
	}
	if err != nil {
		return
	}
	c.client = gatesclient.New(cfg.transport, strfmt.Default)

	return
}

func (c *gatesFeature) register(ctx *godog.ScenarioContext) {
	//ctx.Step(`^a valid store owner$`, c.noop)
	ctx.Step(`^the store (?:called )?"([^"]*)" already exists$`, c.iAskedTimeCalled)
	//ctx.Step(`^I create (?:the|a) store called "([^"]*)"$`, c.iCreateTheStoreCalled)
	//ctx.Step(`^(?:I )?(?:ensure |expect )?the store (?:was|is) created$`, c.expectTheStoreWasCreated)
	//ctx.Step(`^(?:I )?(?:ensure |expect )?a store called "([^"]*)" (?:to )?exists?$`, c.expectAStoreCalledToExist)
	//ctx.Step(`^(?:I )?(?:ensure |expect )?no store called "([^"]*)" (?:to )?exists?$`, c.expectNoStoreCalledToExist)
	//
	//ctx.Step(`^I create the product called "([^"]*)"$`, c.iCreateTheProductCalled)
	//ctx.Step(`^I create the product called "([^"]*)" with price "([^"]*)"$`, c.iCreateTheProductCalledWithPrice)
	//ctx.Step(`^(?:I )?(?:ensure |expect )?the product (?:was|is) created$`, c.expectTheProductWasCreated)
	//ctx.Step(`^(?:I )?(?:ensure |expect )?a product called "([^"]*)" (?:to )?exists?$`, c.expectAProductCalledToExist)
	//ctx.Step(`^(?:I )?(?:ensure |expect )?no product called "([^"]*)" (?:to )?exists?$`, c.expectNoProductCalledToExist)
	//
	//ctx.Step(`^a store has the following items$`, c.aStoreHasTheFollowingItems)
}

func (c *gatesFeature) noop() {
	// noop
}
func (c *gatesFeature) reset() {
	//truncate := func(tableName string) {
	//	_, _ = c.db.Exec(fmt.Sprintf("TRUNCATE %s", tableName))
	//}
	//
	//truncate("gates.gates")
	//truncate("gates.products")
	//truncate("gates.events")
	//truncate("gates.snapshots")
	//truncate("gates.inbox")
	//truncate("gates.outbox")
}

func (c *gatesFeature) iAskedTimeCalled(ctx context.Context, name string) context.Context {
	resp, err := c.client.Time.GetTime(time.NewGetTimeParams().WithBody(&models.GatespbGetTimeRequest{
		Ask: "What time is it?",
	}))
	ctx = setLastResponseAndError(ctx, resp, err)
	if err != nil {
		return ctx
	}
	return context.WithValue(ctx, timeKey{}, resp.Payload.Time)
}

func (c *gatesFeature) expectTheStoreWasCreated(ctx context.Context) error {
	if err := lastResponseWas(ctx, &time.GetTimeOK{}); err != nil {
		return err
	}

	return nil
}
