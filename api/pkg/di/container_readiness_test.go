package di

import (
	"context"
	"testing"
	"time"
)

func TestReadinessReportFailsWhenCoreDependenciesAreMissing(t *testing.T) {
	t.Setenv("ENV", "local")
	t.Setenv("DATABASE_URL", "")
	t.Setenv("DATABASE_URL_DEDICATED", "")
	t.Setenv("REDIS_URL", "")
	t.Setenv("HTTPSMS_READY_OWNER", "")
	t.Setenv("HTTPSMS_READY_PHONE", "")

	report := NewLiteContainer().ReadinessReport(context.Background())
	if report.Status != "unhealthy" {
		t.Fatalf("expected unhealthy report, got %q", report.Status)
	}

	checks := map[string]readinessCheck{}
	for _, check := range report.Checks {
		checks[check.Name] = check
	}
	for _, name := range []string{"database", "dedicated_database", "redis"} {
		if checks[name].Status != "error" {
			t.Fatalf("expected %s to be error, got %#v", name, checks[name])
		}
	}
	if checks["phone_heartbeat"].Status != "skipped" {
		t.Fatalf("expected phone heartbeat to be skipped without HTTPSMS_READY_OWNER, got %#v", checks["phone_heartbeat"])
	}
}

func TestReadinessReportRequiresDatabaseForConfiguredPhoneHeartbeat(t *testing.T) {
	t.Setenv("ENV", "local")
	t.Setenv("DATABASE_URL", "")
	t.Setenv("DATABASE_URL_DEDICATED", "")
	t.Setenv("REDIS_URL", "")
	t.Setenv("HTTPSMS_READY_OWNER", "+919769011110")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	report := NewLiteContainer().ReadinessReport(ctx)

	var phone readinessCheck
	for _, check := range report.Checks {
		if check.Name == "phone_heartbeat" {
			phone = check
			break
		}
	}
	if phone.Status != "error" {
		t.Fatalf("expected phone heartbeat config error, got %#v", phone)
	}
	if phone.Error != "database unavailable for phone heartbeat check" {
		t.Fatalf("expected DB-unavailable phone readiness error, got %q", phone.Error)
	}
}
