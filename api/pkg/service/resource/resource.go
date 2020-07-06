package resource

import (
	"context"
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"

	"github.com/tektoncd/hub/api/gen/resource"
	"github.com/tektoncd/hub/api/pkg/app"
	"github.com/tektoncd/hub/api/pkg/db/model"
)

type service struct {
	logger *zap.SugaredLogger
	db     *gorm.DB
}

var replaceStrings = strings.NewReplacer("github.com", "raw.githubusercontent.com", "/tree/", "/")

// Errors
var (
	fetchError    = resource.MakeInternalError(fmt.Errorf("Failed to fetch resources"))
	notFoundError = resource.MakeNotFound(fmt.Errorf("Resource not found"))
)

// New returns the resource service implementation.
func New(api app.Config) resource.Service {
	return &service{api.Logger(), api.DB()}
}

// Find resources based on name, type or both
func (s *service) Query(ctx context.Context, p *resource.QueryPayload) (res resource.ResourceCollection, err error) {

	q := s.db.Order("rating DESC, name").Limit(p.Limit).
		Preload("Catalog").
		Preload("Versions", func(db *gorm.DB) *gorm.DB {
			return db.Order("string_to_array(version, '.')::int[];")
		}).
		Preload("Tags", func(db *gorm.DB) *gorm.DB {
			return db.Order("tags.name ASC")
		})

	if p.Type != "" {
		q = q.Where("LOWER(type) = ?", p.Type)
	}

	if p.Name != "" {
		name := "%" + strings.ToLower(p.Name) + "%"
		q = q.Where("LOWER(name) LIKE ?", name)
	}

	var all []model.Resource
	if err := q.Find(&all).Error; err != nil {
		s.logger.Error(err)
		return nil, fetchError
	}

	return resourceCollection(all)
}

func resourceCollection(rs []model.Resource) (resource.ResourceCollection, error) {
	if len(rs) == 0 {
		return nil, notFoundError
	}

	res := resource.ResourceCollection{}
	for _, r := range rs {
		res = append(res, initResource(r))
	}

	return res, nil
}

func initResource(r model.Resource) *resource.Resource {
	res := &resource.Resource{}
	res.ID = r.ID
	res.Name = r.Name
	res.Catalog = &resource.Catalog{
		ID:   r.Catalog.ID,
		Type: r.Catalog.Type,
	}
	res.Type = r.Type
	res.Rating = r.Rating

	lv := (r.Versions)[len(r.Versions)-1]
	res.LatestVersion = &resource.Version{
		ID:                  lv.ID,
		Version:             lv.Version,
		Description:         lv.Description,
		DisplayName:         lv.DisplayName,
		MinPipelinesVersion: lv.MinPipelinesVersion,
		WebURL:              lv.URL,
		RawURL:              replaceStrings.Replace(lv.URL),
		UpdatedAt:           lv.UpdatedAt.UTC().String(),
	}
	for _, tag := range r.Tags {
		res.Tags = append(res.Tags, &resource.Tag{
			ID:   tag.ID,
			Name: tag.Name,
		})
	}

	return res
}