package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"flamingo.me/flamingo/v3/framework/opencensus"
	"github.com/tessig/flamingo-mysql/db"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"go.opencensus.io/trace"

	"github.com/tessig/flamingo-product-rating/src/app/domain"
)

type (
	// RatingRepository implements the MySQL connection for Rating persistence
	RatingRepository struct {
		db db.DB
	}
)

var (
	_ domain.RatingRepository = new(RatingRepository)

	stat            = stats.Int64("rating/metrics/amount", "Amount of ratings", stats.UnitDimensionless)
	keyProductID, _ = tag.NewKey("productID")
)

func init() {
	_ = opencensus.View("rating/metrics/amount", stat, view.Count(), keyProductID)
}

// Inject dependencies
func (r *RatingRepository) Inject(db db.DB) {
	r.db = db
}

// Count returns the number of all ratings
func (r *RatingRepository) Count(ctx context.Context) (int64, error) {
	ctx, span := trace.StartSpan(ctx, "app/ratingrepository/Count")
	defer span.End()

	return r.count(ctx, "")
}

// CountByProductID returns the number of all ratings for a given product ID
func (r *RatingRepository) CountByProductID(ctx context.Context, pid int) (int64, error) {
	ctx, span := trace.StartSpan(ctx, "app/ratingrepository/CountByProductID")
	defer span.End()

	return r.count(ctx, "where product_id=?", pid)
}

func (r *RatingRepository) count(ctx context.Context, where string, args ...interface{}) (int64, error) {
	ctx, span := trace.StartSpan(ctx, "app/ratingrepository/count")
	defer span.End()

	var count int64
	err := r.db.Connection().GetContext(ctx, &count, "select COUNT(*) from ratings "+where, args...)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// List off all Ratings
func (r *RatingRepository) List(ctx context.Context) ([]*domain.Rating, error) {
	ctx, span := trace.StartSpan(ctx, "app/ratingrepository/List")
	defer span.End()

	return r.list(ctx, "")
}

// ListByProductID returns all ratings for a given product ID
func (r *RatingRepository) ListByProductID(ctx context.Context, pid int) ([]*domain.Rating, error) {
	ctx, span := trace.StartSpan(ctx, "app/ratingrepository/ListByProductID")
	defer span.End()

	return r.list(ctx, "where product_id=?", pid)
}

func (r *RatingRepository) list(ctx context.Context, where string, args ...interface{}) ([]*domain.Rating, error) {
	ctx, span := trace.StartSpan(ctx, "app/ratingrepository/list")
	defer span.End()

	var ratings []*domain.Rating
	err := r.db.Connection().SelectContext(ctx, &ratings, "select * from ratings "+where+" order by created_at DESC", args...)
	if err != nil {
		return nil, err
	}

	return ratings, nil
}

// Get returns a specific rating
func (r *RatingRepository) Get(ctx context.Context, id int) (*domain.Rating, error) {
	ctx, span := trace.StartSpan(ctx, "app/ratingrepository/Get")
	defer span.End()

	rating := &domain.Rating{}
	err := r.db.Connection().GetContext(ctx, rating, "select * from ratings where id=?", id)
	if err != nil {
		return nil, err
	}

	return rating, nil
}

// Save creates or updates a given rating based on its ID
func (r *RatingRepository) Save(ctx context.Context, rating *domain.Rating) error {
	ctx, span := trace.StartSpan(ctx, "app/ratingrepository/Save")
	defer span.End()

	result, err := r.db.Connection().NamedExecContext(
		ctx,
		`INSERT into ratings 
    	(id, name, title, text, created_at, product_id, stars) 
        VALUES (:id,:name,:title,:text,:created_at,:product_id,:stars) 
        ON DUPLICATE KEY UPDATE 
        name=:name,title=:title,text=:text,created_at=:created_at,product_id=:product_id,stars=:stars`,
		rating,
	)
	if err != nil {
		return err
	}
	if affected, _ := result.RowsAffected(); affected == 0 {
		return errors.New("no entity has been deleted")
	}

	ctx, _ = tag.New(ctx, tag.Upsert(keyProductID, strconv.Itoa(rating.ProductID)), tag.Upsert(opencensus.KeyArea, "-"))
	stats.Record(ctx, stat.M(1))

	return nil
}

// Delete the given rating
func (r *RatingRepository) Delete(ctx context.Context, rating *domain.Rating) error {
	ctx, span := trace.StartSpan(ctx, "app/ratingrepository/Delete")
	defer span.End()

	result, err := r.db.Connection().ExecContext(ctx, "delete from ratings where id=$1", rating.ID)
	if err != nil {
		return err
	}
	if affected, _ := result.RowsAffected(); affected == 0 {
		return errors.New("no entity has been deleted")
	}

	ctx, _ = tag.New(ctx, tag.Upsert(keyProductID, strconv.Itoa(rating.ProductID)), tag.Upsert(opencensus.KeyArea, "-"))
	stats.Record(ctx, stat.M(-1))

	return nil
}

// Amounts returns a map of product IDs an the corresponding number of ratings
func (r *RatingRepository) Amounts(ctx context.Context) (*domain.RatingAmounts, error) {
	ctx, span := trace.StartSpan(ctx, "app/ratingrepository/Amounts")
	defer span.End()

	var data []struct {
		ID     int
		Amount int
	}

	err := r.db.Connection().SelectContext(ctx, &data, "select product_id as id,COUNT(*) as amount from ratings group by product_id order by product_id;")
	if err != nil {
		return nil, err
	}

	result := make(domain.RatingAmounts, len(data))
	for _, d := range data {
		result[d.ID] = d.Amount
	}

	return &result, nil
}

// Average returns the average statistics of all data
func (r *RatingRepository) Average(ctx context.Context) (*domain.RatingAverage, error) {
	ctx, span := trace.StartSpan(ctx, "app/ratingrepository/Average")
	defer span.End()

	return r.average(ctx, "")
}

// AverageByProductID returns the average statistics of all data for a specific product ID
func (r *RatingRepository) AverageByProductID(ctx context.Context, pid int) (*domain.RatingAverage, error) {
	ctx, span := trace.StartSpan(ctx, "app/ratingrepository/AverageByProductID")
	defer span.End()

	return r.average(ctx, "where product_id=?", pid)
}

func (r *RatingRepository) average(ctx context.Context, where string, args ...interface{}) (*domain.RatingAverage, error) {
	ctx, span := trace.StartSpan(ctx, "app/ratingrepository/average")
	defer span.End()

	var data struct {
		Avg    float64
		Amount int
	}
	c, err := r.count(ctx, where, args...)
	if err != nil {
		return nil, err
	}

	if c == 0 {
		return &domain.RatingAverage{
				Value:  0,
				Amount: 0,
				Max:    5,
				Stars:  0,
			},
			nil
	}

	err = r.db.Connection().GetContext(ctx, &data, "select AVG(stars) as avg,COUNT(*) as amount from ratings "+where, args...)
	if err != nil {
		return nil, err
	}

	return &domain.RatingAverage{
			Value:  data.Avg,
			Amount: data.Amount,
			Max:    5,
			Stars:  int(data.Avg),
		},
		nil
}

// Breakdown returns the statistics by single stars
func (r *RatingRepository) Breakdown(ctx context.Context) ([]*domain.RatingBreakdown, error) {
	ctx, span := trace.StartSpan(ctx, "app/ratingrepository/Breakdown")
	defer span.End()

	return r.breakdown(ctx, "")
}

// BreakdownByProductID returns the statistics by single stars for a given product
func (r *RatingRepository) BreakdownByProductID(ctx context.Context, pid int) ([]*domain.RatingBreakdown, error) {
	ctx, span := trace.StartSpan(ctx, "app/ratingrepository/BreakdownByProductID")
	defer span.End()

	return r.breakdown(ctx, "where product_id=?", pid)
}

func (r *RatingRepository) breakdown(ctx context.Context, where string, args ...interface{}) ([]*domain.RatingBreakdown, error) {
	ctx, span := trace.StartSpan(ctx, "app/ratingrepository/breakdown")
	defer span.End()

	var total int
	var data []*struct {
		Stars  int
		Amount int
	}
	err := r.db.Connection().GetContext(ctx, &total, "select count(*) from ratings "+where, args...)
	if err != nil {
		return nil, err
	}

	err = r.db.Connection().SelectContext(ctx, &data, "select stars,count(*) as amount from ratings "+where+" group by stars order by stars DESC", args...)
	if err != nil {
		return nil, err
	}

	breakdown := make([]*domain.RatingBreakdown, len(data))
	for i, d := range data {
		breakdown[i] = &domain.RatingBreakdown{
			Stars:      d.Stars,
			Amount:     d.Amount,
			Total:      total,
			Percentage: float64(d.Amount*100) / float64(total),
		}
	}

	return breakdown, nil
}
