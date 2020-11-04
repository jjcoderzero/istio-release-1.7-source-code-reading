package status

import "istio.io/pkg/monitoring"

var (
	typeTag = monitoring.MustCreateLabel("type")

	// scrapeErrors records total number of failed scrapes.
	scrapeErrors = monitoring.NewSum(
		"scrape_failures_total",
		"The total number of failed scrapes.",
		monitoring.WithLabels(typeTag),
	)
	envoyScrapeErrors = scrapeErrors.With(typeTag.Value(ScrapeTypeEnvoy))
	appScrapeErrors   = scrapeErrors.With(typeTag.Value(ScrapeTypeApp))
	agentScrapeErrors = scrapeErrors.With(typeTag.Value(ScrapeTypeAgent))

	// scrapeErrors records total number of scrapes.
	scrapeTotals = monitoring.NewSum(
		"scrapes_total",
		"The total number of scrapes.",
	)
)

var (
	ScrapeTypeEnvoy = "envoy"
	ScrapeTypeApp   = "application"
	ScrapeTypeAgent = "agent"
)

func init() {
	monitoring.MustRegister(
		scrapeTotals,
		scrapeErrors,
	)
}
