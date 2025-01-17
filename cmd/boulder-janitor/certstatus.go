package main

import (
	"github.com/jmhodges/clock"
	blog "github.com/letsencrypt/boulder/log"
)

// newCertificateStatusJob returns a batchedDBJob configured to delete expired
// rows from the certificateStatus table.
func newCertificateStatusJob(
	db janitorDB,
	log blog.Logger,
	clk clock.Clock,
	config Config) *batchedDBJob {
	purgeBefore := config.Janitor.CertificateStatus.GracePeriod.Duration
	workQuery := `SELECT id, notAfter AS expires FROM certificateStatus
		 WHERE
		   id > :startID
		 LIMIT :limit`
	log.Debugf("Creating CertificateStatus job from config: %#v\n", config.Janitor.CertificateStatus)
	return &batchedDBJob{
		db:          db,
		log:         log,
		clk:         clk,
		purgeBefore: purgeBefore,
		workSleep:   config.Janitor.CertificateStatus.WorkSleep.Duration,
		batchSize:   config.Janitor.CertificateStatus.BatchSize,
		maxDPS:      config.Janitor.CertificateStatus.MaxDPS,
		parallelism: config.Janitor.CertificateStatus.Parallelism,
		table:       "certificateStatus",
		workQuery:   workQuery,
	}
}
