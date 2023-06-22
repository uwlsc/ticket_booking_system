package cache

import (
	"ticket-api/pkg"

	"github.com/dgraph-io/ristretto"
)

var cacheManager *ristretto.Cache

func GetCacheManager() *ristretto.Cache {
	if cacheManager == nil {
		cacheManager = SetUpCacheStore(pkg.GetLogger())
	}
	return cacheManager
}

func SetUpCacheStore(logger pkg.Logger) *ristretto.Cache {
	ristrettoCache, err := ristretto.NewCache(
		&ristretto.Config{
			NumCounters: 1000,
			MaxCost:     50000000,
			BufferItems: 64,
		},
	)
	if err != nil {
		logger.Fatalf("Unable to initialize cache in project")
	}

	return ristrettoCache
}
