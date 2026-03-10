package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"net/http"

	"git.zyner.org/meta/xyter/internal/api"
	"git.zyner.org/meta/xyter/internal/config"
	"git.zyner.org/meta/xyter/internal/database"
	"github.com/dromara/carbon/v2"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	k := config.Load()
	db := database.Open(k)
	database.Migrate(db)
	carbon.SetDefault(carbon.Default{
		Layout:       carbon.DateTimeLayout,
		Timezone:     k.String("tz"),
		Locale:       "en",
		WeekStartsAt: carbon.Monday,
		WeekendDays:  []carbon.Weekday{carbon.Saturday, carbon.Sunday},
	})
	api.Start(k, db)

	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":2112", nil)

	log.Println("API is running. Press CTRL-C to exit.")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop
}
