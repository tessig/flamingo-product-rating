package infrastructure

import (
	"errors"

	"github.com/tessig/flamingo-mysql/db"
	"github.com/tessig/flamingo-product-rating/src/app/domain"
)

type (
	// RatingRepository implements the MySQL connection for Rating persistence
	RatingRepository struct {
		db db.DB
	}
)

// Inject dependencies
func (r *RatingRepository) Inject(db db.DB) {
	r.db = db
}

// Count returns the number of all ratings
func (r *RatingRepository) Count() (int64, error) {
	return r.count("")
}

// CountByProductID returns the number of all ratings for a given product ID
func (r *RatingRepository) CountByProductID(pid int) (int64, error) {
	return r.count("where product_id=?", pid)
}

func (r *RatingRepository) count(where string, args ...interface{}) (int64, error) {
	var count int64
	err := r.db.Connection().Get(&count, "select COUNT(*) from ratings "+where, args...)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// List off all Ratings
func (r *RatingRepository) List() ([]*domain.Rating, error) {
	return r.list("")
}

// ListByProductID returns all ratings for a given product ID
func (r *RatingRepository) ListByProductID(pid int) ([]*domain.Rating, error) {
	return r.list("where product_id=?", pid)
}

func (r *RatingRepository) list(where string, args ...interface{}) ([]*domain.Rating, error) {
	var ratings []*domain.Rating
	err := r.db.Connection().Select(&ratings, "select * from ratings "+where+" order by created_at DESC", args...)
	if err != nil {
		return nil, err
	}

	return ratings, nil
}

// Get returns a specific rating
func (r *RatingRepository) Get(id int) (*domain.Rating, error) {
	rating := &domain.Rating{}
	err := r.db.Connection().Get(rating, "select * from ratings where id=?", id)
	if err != nil {
		return nil, err
	}

	return rating, nil
}

// Save creates or updates a given rating based on its ID
func (r *RatingRepository) Save(rating *domain.Rating) error {
	result, err := r.db.Connection().NamedExec(
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

	return nil
}

// Delete the given rating
func (r *RatingRepository) Delete(rating *domain.Rating) error {
	result, err := r.db.Connection().Exec("delete from ratings where id=$1", rating.ID)
	if err != nil {
		return err
	}
	if affected, _ := result.RowsAffected(); affected == 0 {
		return errors.New("no entity has been deleted")
	}

	return nil
}

// Amounts returns a map of product IDs an the corresponding number of ratings
func (r *RatingRepository) Amounts() (*domain.RatingAmounts, error) {
	var data []struct {
		ID     int
		Amount int
	}

	err := r.db.Connection().Select(&data, "select product_id as id,COUNT(*) as amount from ratings group by product_id order by product_id ASC;")
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
func (r *RatingRepository) Average() (*domain.RatingAverage, error) {
	return r.average("")
}

// AverageByProductID returns the average statistics of all data for a specific product ID
func (r *RatingRepository) AverageByProductID(pid int) (*domain.RatingAverage, error) {
	return r.average("where product_id=?", pid)
}

func (r *RatingRepository) average(where string, args ...interface{}) (*domain.RatingAverage, error) {
	var data struct {
		Avg    float64
		Amount int
	}
	c, err := r.count(where, args...)
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

	err = r.db.Connection().Get(&data, "select AVG(stars) as avg,COUNT(*) as amount from ratings "+where, args...)
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
func (r *RatingRepository) Breakdown() ([]*domain.RatingBreakdown, error) {
	return r.breakdown("")
}

// BreakdownByProductID returns the statistics by single stars for a given product
func (r *RatingRepository) BreakdownByProductID(pid int) ([]*domain.RatingBreakdown, error) {
	return r.breakdown("where product_id=?", pid)
}

func (r *RatingRepository) breakdown(where string, args ...interface{}) ([]*domain.RatingBreakdown, error) {
	var total int
	var data []*struct {
		Stars  int
		Amount int
	}
	err := r.db.Connection().Get(&total, "select count(*) from ratings "+where, args...)
	if err != nil {
		return nil, err
	}

	err = r.db.Connection().Select(&data, "select stars,count(*) as amount from ratings "+where+" group by stars order by stars DESC", args...)
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
