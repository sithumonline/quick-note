package note

import (
	"gorm.io/gorm"

	log "github.com/sirupsen/logrus"

	"github.com/sithumonline/quick-note/models"
)

type NoteRepo interface {
	Save(note *models.Note) error
	GetList(uid string) ([]models.Note, error)
	Get(id string, uid string) (*models.Note, error)
	Delete(id string, uid string) error
	Update(note *models.Note, id string, uid string) error
	Migrate() error
}

type Note struct {
	db *gorm.DB
}

func NewNoteRepo(db *gorm.DB) *Note {
	return &Note{
		db: db,
	}
}

func (p *Note) Save(note *models.Note) error {
	if result := p.db.Create(&note); result.Error != nil {
		log.Errorf("failed to create note: %+v: %v", p, result.Error)
		return result.Error
	}

	return nil
}

func (p *Note) GetList(uid string) ([]models.Note, error) {
	list := make([]models.Note, 0)

	if result := p.db.Find(&list, "user_id = ?", uid); result.Error != nil {
		log.Errorf("failed to find note: %+v: %v", p, result.Error)
		return nil, result.Error
	}

	return list, nil
}

func (p *Note) Get(id string, uid string) (*models.Note, error) {
	t := &models.Note{}

	if result := p.db.Where("id = ? AND user_id = ?", id, uid).First(t); result.Error != nil {
		log.Errorf("failed to find note: %+v: %v", p, result.Error)
		return t, result.Error
	}

	return t, nil
}

func (p *Note) Delete(id string, uid string) error {
	if result := p.db.Model(&p).Where("id = ? AND user_id = ?", id, uid).Delete(p); result.Error != nil {
		log.Errorf("failed to delete note: %+v: %v", p, result.Error)
		return result.Error
	}

	return nil
}

func (p *Note) Update(note *models.Note, id string, uid string) error {
	if !note.Public {
		if result := p.db.Model(&note).Where("id = ? AND user_id = ?", id, uid).Updates(map[string]interface{}{
			"public": false,
		}); result.Error != nil {
			log.Errorf("failed to update note: %+v: %v", p, result.Error)
			return result.Error
		}

		return nil
	}

	if result := p.db.Model(&note).Where("id = ? AND user_id = ?", id, uid).Updates(note); result.Error != nil {
		log.Errorf("failed to update note: %+v: %v", p, result.Error)
		return result.Error
	}

	return nil
}

func (p *Note) Migrate() error {
	p.db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	if err := p.db.AutoMigrate(models.Note{}); err != nil {
		log.Errorf("failed to migrate note: %+v: %v", models.Note{}, err)
		return err
	}

	return nil
}
