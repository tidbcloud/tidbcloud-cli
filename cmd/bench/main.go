package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/time/rate"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/config/store"
	"github.com/tidbcloud/tidbcloud-cli/internal/iostream"
	"github.com/tidbcloud/tidbcloud-cli/internal/prop"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/aws/s3"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless/cluster"

	"github.com/spf13/viper"
	"github.com/zalando/go-keyring"
)

const (
	defaultProjectID     = "1369847559694040868"
	defaultRegion        = "regions/aws-us-east-1"
	defaultNamePrefix    = "keep--1h"
	defaultSpendingLimit = 10
	defaultConcurrency   = 5
	defaultTotal         = 100
	defaultRPS           = 2.0
	waitInterval         = 2 * time.Second
	waitTimeout          = 10 * time.Minute
)

type benchConfig struct {
	concurrency   int
	rps           float64
	total         int
	projectID     string
	region        string
	namePrefix    string
	spendingLimit int
	minRcu        int
	maxRcu        int
	encryption    bool
	disablePub    bool
	waitReady     bool
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	initBenchConfig()
	config.SetActiveProfile(viper.GetString(prop.CurProfile))

	cfg := parseFlags()
	h := newHelper()

	client, err := h.Client()
	if err != nil {
		log.Fatalf("init client: %v", err)
	}

	runBench(ctx, client, cfg)
}

func parseFlags() benchConfig {
	cfg := benchConfig{
		concurrency:   defaultConcurrency,
		rps:           defaultRPS,
		total:         defaultTotal,
		projectID:     defaultProjectID,
		region:        defaultRegion,
		namePrefix:    defaultNamePrefix,
		spendingLimit: defaultSpendingLimit,
	}

	flag.IntVar(&cfg.concurrency, "concurrency", cfg.concurrency, "number of concurrent workers")
	flag.Float64Var(&cfg.rps, "rps", cfg.rps, "requests per second")
	flag.IntVar(&cfg.total, "total", cfg.total, "total number of clusters to create")
	flag.StringVar(&cfg.projectID, "project-id", cfg.projectID, "project id")
	flag.StringVar(&cfg.region, "region", cfg.region, "region name")
	flag.StringVar(&cfg.namePrefix, "name-prefix", cfg.namePrefix, "prefix of the cluster name")
	flag.IntVar(&cfg.spendingLimit, "spending-limit", cfg.spendingLimit, "monthly spending limit in USD cents, Starter only")
	flag.IntVar(&cfg.minRcu, "min-rcu", 0, "minimum RCU, Essential only")
	flag.IntVar(&cfg.maxRcu, "max-rcu", 0, "maximum RCU, Essential only")
	flag.BoolVar(&cfg.encryption, "encryption", false, "enable enhanced encryption")
	flag.BoolVar(&cfg.disablePub, "disable-public-endpoint", false, "disable public endpoint")
	flag.BoolVar(&cfg.waitReady, "wait-ready", true, "wait for cluster to be ACTIVE")
	flag.Parse()

	if cfg.total <= 0 {
		log.Fatalf("total must be positive")
	}

	if cfg.concurrency <= 0 {
		log.Fatalf("concurrency must be positive")
	}

	if cfg.rps <= 0 {
		log.Fatalf("rps must be positive")
	}

	if (cfg.minRcu > 0 || cfg.maxRcu > 0) && cfg.minRcu > cfg.maxRcu {
		log.Fatalf("min-rcu cannot exceed max-rcu")
	}

	return cfg
}

func initBenchConfig() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("get home: %v", err)
	}
	path := filepath.Join(home, config.HomePath)
	if err := os.MkdirAll(path, 0700); err != nil {
		log.Fatalf("init config dir: %v", err)
	}

	viper.AddConfigPath(path)
	viper.SetConfigType("toml")
	viper.SetConfigName("config")
	viper.SetConfigPermissions(0600)
	if err := viper.SafeWriteConfig(); err != nil {
		var existErr viper.ConfigFileAlreadyExistsError
		if !errors.As(err, &existErr) {
			log.Fatalf("write config: %v", err)
		}
	}
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("read config: %v", err)
	}
}

func newHelper() *internal.Helper {
	return &internal.Helper{
		Client: func() (cloud.TiDBCloudClient, error) {
			publicKey, privateKey := config.GetPublicKey(), config.GetPrivateKey()
			serverlessEndpoint := config.GetServerlessEndpoint()
			if serverlessEndpoint == "" {
				serverlessEndpoint = cloud.DefaultServerlessEndpoint
			}
			iamEndpoint := config.GetIAMEndpoint()
			if iamEndpoint == "" {
				iamEndpoint = cloud.DefaultIAMEndpoint
			}

			if publicKey != "" && privateKey != "" {
				return cloud.NewClientDelegateWithApiKey(publicKey, privateKey, serverlessEndpoint, iamEndpoint)
			}

			if err := config.ValidateToken(); err != nil {
				return nil, err
			}
			token, err := config.GetAccessToken()
			if err != nil {
				if errors.Is(err, keyring.ErrNotFound) || errors.Is(err, store.ErrNotSupported) {
					return nil, err
				}
				return nil, err
			}
			return cloud.NewClientDelegateWithToken(token, serverlessEndpoint, iamEndpoint)
		},
		Uploader: func(client cloud.TiDBCloudClient) s3.Uploader {
			return s3.NewUploader(client)
		},
		QueryPageSize: internal.DefaultPageSize,
		IOStreams:     iostream.System(),
	}
}

func runBench(ctx context.Context, client cloud.TiDBCloudClient, cfg benchConfig) {
	limiter := rate.NewLimiter(rate.Limit(cfg.rps), int(math.Ceil(cfg.rps)))
	jobs := make(chan int, cfg.total)

	var success int64
	var failed int64

	var wg sync.WaitGroup

	timestamp := time.Now().Unix()
	for i := 0; i < cfg.concurrency; i++ {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()
			for idx := range jobs {
				if err := limiter.Wait(ctx); err != nil {
					log.Printf("worker %d rate wait err: %v", worker, err)
					continue
				}
				name := fmt.Sprintf("%s-%d-%d", cfg.namePrefix, timestamp, idx)
				start := time.Now()
				id, err := createOnce(ctx, client, cfg, name)
				if err != nil {
					atomic.AddInt64(&failed, 1)
					log.Printf("worker %d create %s failed: %v", worker, name, err)
					continue
				}

				if cfg.waitReady {
					if err := waitClusterReady(ctx, client, id); err != nil {
						atomic.AddInt64(&failed, 1)
						log.Printf("worker %d wait %s failed: %v", worker, id, err)
						continue
					}
				}

				atomic.AddInt64(&success, 1)
				log.Printf("worker %d create %s (id=%s) ok in %s", worker, name, id, time.Since(start))
			}
		}(i)
	}

	for i := 0; i < cfg.total; i++ {
		jobs <- i
	}
	close(jobs)

	wg.Wait()
	log.Printf("bench done: success=%d failed=%d", success, failed)
}

func createOnce(ctx context.Context, client cloud.TiDBCloudClient, cfg benchConfig, name string) (string, error) {
	payload := &cluster.TidbCloudOpenApiserverlessv1beta1Cluster{
		DisplayName: name,
		Region: cluster.Commonv1beta1Region{
			Name: &cfg.region,
		},
	}

	if cfg.projectID != "" {
		payload.Labels = &map[string]string{"tidb.cloud/project": cfg.projectID}
	}
	if cfg.spendingLimit > 0 {
		payload.SpendingLimit = &cluster.ClusterSpendingLimit{
			Monthly: toInt32Ptr(int32(cfg.spendingLimit)),
		}
	}
	if cfg.minRcu > 0 || cfg.maxRcu > 0 {
		payload.AutoScaling = &cluster.V1beta1ClusterAutoScaling{
			MinRcu: toInt64Ptr(int64(cfg.minRcu)),
			MaxRcu: toInt64Ptr(int64(cfg.maxRcu)),
		}
	}
	if cfg.encryption {
		payload.EncryptionConfig = &cluster.V1beta1ClusterEncryptionConfig{
			EnhancedEncryptionEnabled: &cfg.encryption,
		}
	}
	if cfg.disablePub {
		payload.Endpoints = &cluster.V1beta1ClusterEndpoints{
			Public: &cluster.EndpointsPublic{
				Disabled: &cfg.disablePub,
			},
		}
	}

	resp, err := client.CreateCluster(ctx, payload)
	if err != nil {
		return "", err
	}
	if resp.ClusterId == nil {
		return "", fmt.Errorf("empty cluster id")
	}
	return *resp.ClusterId, nil
}

func waitClusterReady(ctx context.Context, client cloud.TiDBCloudClient, clusterID string) error {
	ticker := time.NewTicker(waitInterval)
	defer ticker.Stop()
	timer := time.After(waitTimeout)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer:
			return fmt.Errorf("timeout waiting for cluster %s ready", clusterID)
		case <-ticker.C:
			c, err := client.GetCluster(ctx, clusterID, cluster.CLUSTERSERVICEGETCLUSTERVIEWPARAMETER_BASIC)
			if err != nil {
				return err
			}
			if c.State != nil && *c.State == cluster.COMMONV1BETA1CLUSTERSTATE_ACTIVE {
				return nil
			}
		}
	}
}

func toInt32Ptr(v int32) *int32 {
	return &v
}

func toInt64Ptr(v int64) *int64 {
	val := int64(v)
	return &val
}
