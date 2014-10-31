package services

import (
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

func (this *ActivitiesService) SaveActivity(activity *models.Activity) error {
	var err error
	var existingActivity *models.Activity

	if activity.Id == -1 {
		existingActivity, err = this.db.TryGetActivity(activity.UserId, activity.Day, activity.Text)
		if err != nil {
			return err
		}
	}

	if existingActivity != nil {
		existingActivity.Duration += activity.Duration
		return this.db.SaveActivity(existingActivity)
	}

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
