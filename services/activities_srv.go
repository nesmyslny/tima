package services

import (
	"log"
	"time"

	"github.com/nesmyslny/tima/dbaccess"
	"github.com/nesmyslny/tima/models"
)

type ActivitiesService struct {
	db *DbAccess.Db
}

func NewActivitiesService(db *DbAccess.Db) *ActivitiesService {
	return &ActivitiesService{db}
}

func (this *ActivitiesService) GetActivities(userId int, day time.Time) ([]models.Activity, error) {
	activities, err := this.db.GetActivities(userId, day)
	if err != nil {
		return nil, err
	}
	return activities, nil
}

func (this *ActivitiesService) AddActivity(activity *models.Activity) error {
	existingActivity, err := this.db.TryGetActivity(activity.UserId, activity.Day, activity.Text)
	if err != nil {
		log.Print(err.Error())
		return err
	}
	if existingActivity != nil {
		existingActivity.Duration += activity.Duration
		return this.db.SaveActivity(existingActivity)
	}
	activity.Id = -1
	return this.db.SaveActivity(activity)
}

func (this *ActivitiesService) DeleteActivity(id int) error {
	activity, err := this.db.GetActivity(id)
	if err != nil {
		return err
	}

	err = this.db.DeleteActivity(activity)
	if err != nil {
		return err
	}

	return nil
}
