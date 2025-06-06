package service

import (
	"fmt"
	"goiam/internal/logger"
	"goiam/internal/model"
	"time"
)

const (
	// applicationCacheTTL is the time-to-live for application cache entries
	applicationCacheTTL = 5 * time.Second
)

// cachedApplicationService implements ApplicationService with caching
type cachedApplicationService struct {
	applicationService ApplicationService
	cache              CacheService
}

// NewCachedApplicationService creates a new cached application service
func NewCachedApplicationService(applicationService ApplicationService, cache CacheService) ApplicationService {
	return &cachedApplicationService{
		applicationService: applicationService,
		cache:              cache,
	}
}

// getApplicationCacheKey returns a cache key in the format /<tenant>/<realm>/application/<client_id>
func (s *cachedApplicationService) getApplicationCacheKey(tenant, realm, clientId string) string {
	return fmt.Sprintf("/%s/%s/application/%s", tenant, realm, clientId)
}

func (s *cachedApplicationService) GetApplication(tenant, realm, clientId string) (*model.Application, bool) {
	// Try to get from cache first
	cacheKey := s.getApplicationCacheKey(tenant, realm, clientId)
	if cached, exists := s.cache.Get(cacheKey); exists {
		if app, ok := cached.(*model.Application); ok {
			return app, true
		}
	}

	// If not in cache, get from service
	app, exists := s.applicationService.GetApplication(tenant, realm, clientId)
	if !exists {
		return nil, false
	}

	// Cache the result
	err := s.cache.Cache(cacheKey, app, applicationCacheTTL, 1)
	if err != nil {
		// Log error but continue - caching is not critical
		logger.InfoNoContext("Failed to cache application: %v", err)
	}

	return app, true
}

// Direct pass-through methods (no caching)
func (s *cachedApplicationService) ListApplications(tenant, realm string) ([]model.Application, error) {
	return s.applicationService.ListApplications(tenant, realm)
}

// Direct pass-through methods (no caching)
func (s *cachedApplicationService) ListAllApplications() ([]model.Application, error) {
	return s.applicationService.ListAllApplications()
}

func (s *cachedApplicationService) CreateApplication(tenant, realm string, application model.Application) error {
	err := s.applicationService.CreateApplication(tenant, realm, application)
	if err != nil {
		return err
	}

	// Invalidate cache
	s.invalidateCache(tenant, realm, application.ClientId)
	return nil
}

func (s *cachedApplicationService) UpdateApplication(tenant, realm string, application model.Application) error {
	err := s.applicationService.UpdateApplication(tenant, realm, application)
	if err != nil {
		return err
	}

	// Invalidate cache
	s.invalidateCache(tenant, realm, application.ClientId)
	return nil
}

func (s *cachedApplicationService) DeleteApplication(tenant, realm, clientId string) error {
	err := s.applicationService.DeleteApplication(tenant, realm, clientId)
	if err != nil {
		return err
	}

	// Invalidate cache
	s.invalidateCache(tenant, realm, clientId)
	return nil
}

func (s *cachedApplicationService) RegenerateClientSecret(tenant, realm, clientId string) (string, error) {
	secret, err := s.applicationService.RegenerateClientSecret(tenant, realm, clientId)
	if err != nil {
		return "", err
	}

	// Invalidate cache since client secret has changed
	s.invalidateCache(tenant, realm, clientId)
	return secret, nil
}

func (s *cachedApplicationService) VerifyClientSecret(tenant, realm, clientId, clientSecret string) (bool, error) {
	return s.applicationService.VerifyClientSecret(tenant, realm, clientId, clientSecret)
}

// invalidateCache invalidates the application cache entry
func (s *cachedApplicationService) invalidateCache(tenant, realm, clientId string) {
	cacheKey := s.getApplicationCacheKey(tenant, realm, clientId)
	s.cache.Invalidate(cacheKey)
}
