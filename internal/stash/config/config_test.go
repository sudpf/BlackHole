package config

import "testing"

func TestValidateMetricsConfig(t *testing.T) {
	tests := []struct {
		name    string
		metrics metricsConfig
		wantErr bool
	}{
		{
			name: "disabled metrics skips optional fields",
		},
		{
			name: "enabled metrics",
			metrics: metricsConfig{
				Enabled: true,
				Listen:  "http://127.0.0.1:9102",
				Path:    "/metrics",
			},
		},
		{
			name: "enabled metrics requires valid url",
			metrics: metricsConfig{
				Enabled: true,
				Listen:  "127.0.0.1:9102",
				Path:    "/metrics",
			},
			wantErr: true,
		},
		{
			name: "enabled metrics only supports http",
			metrics: metricsConfig{
				Enabled: true,
				Listen:  "https://127.0.0.1:9102",
				Path:    "/metrics",
			},
			wantErr: true,
		},
		{
			name: "enabled metrics path must be absolute",
			metrics: metricsConfig{
				Enabled: true,
				Listen:  "http://127.0.0.1:9102",
				Path:    "metrics",
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg := validTestConfig()
			cfg.Metrics = test.metrics

			err := validate(cfg)
			if test.wantErr && err == nil {
				t.Fatal("validate() error = nil, want error")
			}
			if !test.wantErr && err != nil {
				t.Fatalf("validate() error = %v", err)
			}
		})
	}
}

func TestMetricsAccessors(t *testing.T) {
	cfg := validTestConfig()
	cfg.Metrics = metricsConfig{
		Enabled: true,
		Listen:  "http://0.0.0.0:9102",
		Path:    "/custom-metrics",
	}

	if !cfg.MetricsEnabled() {
		t.Fatal("MetricsEnabled() = false, want true")
	}
	if got := cfg.MetricsListenAddress(); got != "0.0.0.0:9102" {
		t.Fatalf("MetricsListenAddress() = %q, want 0.0.0.0:9102", got)
	}
	if got := cfg.MetricsPath(); got != "/custom-metrics" {
		t.Fatalf("MetricsPath() = %q, want /custom-metrics", got)
	}
}

func validTestConfig() *Config {
	return &Config{
		App: appConfig{
			Listen:          "http://127.0.0.1:8002",
			RequestTimeout:  1,
			ShutdownTimeout: 1,
		},
		Log: logConfig{
			Level: "info",
			Size:  "256m",
			Dir:   "logs",
		},
		GracePeriod: 1,
	}
}
