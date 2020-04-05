package crud

import (
	"api/models"
	"api/utils/channels"
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

type repositoryPostsCRUD struct {
	db *gorm.DB
}

func NewRepositoryPostsDB(db *gorm.DB) *repositoryPostsCRUD {
	return &repositoryPostsCRUD{db}
}

func (r *repositoryPostsCRUD) Save(post models.Post) (models.Post, error) {
	var err error
	done := make(chan bool)

	go func(ch chan<- bool) {
		defer close(ch)
		err = r.db.Debug().Model(&models.Post{}).Create(&post).Error
		if err != nil {
			ch <- false
			return
		}

		ch <- true
	}(done)

	if channels.OK(done) {
		return post, nil
	}

	return models.Post{}, err
}

func (r *repositoryPostsCRUD) FindAll() ([]models.Post, error) {
	var err error
	posts := []models.Post{}
	done := make(chan bool)

	go func(ch chan<- bool) {
		defer close(ch)
		err = r.db.Debug().Model(&models.Post{}).Limit(100).Find(&posts).Error
		if err != nil {
			ch <- false
			return
		}
		ch <- true
	}(done)

	if channels.OK(done) {
		return posts, nil
	}

	return nil, err
}

func (r *repositoryPostsCRUD) FindByID(uid uint32) (models.Post, error) {
	var err error
	post := models.Post{}
	done := make(chan bool)

	go func(ch chan<- bool) {
		defer close(ch)
		err := r.db.Debug().Model(&models.Post{}).Where("id = ?", uid).Find(&post).Error
		if err != nil {
			ch <- false
			return
		}
		ch <- true
	}(done)

	if channels.OK(done) {
		return post, nil
	}

	if gorm.IsRecordNotFoundError(err) {
		return post, errors.New("Post not found")
	}

	return post, err
}

func (r *repositoryPostsCRUD) Update(uid uint32, post models.Post) (int64, error) {
	var rs *gorm.DB
	done := make(chan bool)
	go func(ch chan<- bool) {
		defer close(ch)
		rs = r.db.Debug().Model(&models.Post{}).Where("id = ?", uid).Take(&models.Post{}).UpdateColumns(
			map[string]interface{}{
				"title":      post.Title,
				"content":    post.Content,
				"updated_at": time.Now(),
			},
		)
		ch <- true
	}(done)

	if gorm.IsRecordNotFoundError(rs.Error) {
		return 0, errors.New("Post not found")
	}

	// if channels.Ok(done) && rs.Error != nil {
	// 	return 0, rs.Error
	// }
	// else {
	// 	rs.RowsAffected, nil
	// }

	if channels.OK(done) {
		if rs.Error != nil {
			return 0, rs.Error
		}
		return rs.RowsAffected, nil
	}

	return 0, rs.Error

}

func (r *repositoryPostsCRUD) Delete(pid int64, uid uint32) (int64, error) {
	var rs *gorm.DB
	done := make(chan bool)

	go func(ch chan<- bool) {
		defer close(ch)
		rs = r.db.Debug().Model(&models.Post{}).Where("id = ? and author_id = ?", pid, uid).Take(&models.Post{}).Delete(&models.Post{})
		ch <- true
	}(done)

	if channels.OK(done) {
		if rs.Error != nil {
			return 0, rs.Error
		}
		return rs.RowsAffected, nil
	}

	return 0, rs.Error
}
