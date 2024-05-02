package models

import "time"

type NotificationRepo struct {
	repo SqlExecutor
}

func InitializeNotificationRepo(db SqlExecutor) *NotificationRepo {
	return &NotificationRepo{repo: db}
}

type NotificationModel interface {
	GetNotificationList(lineUserId string) (NotificationList, error)
}

type NotificationList struct {
	LineUserId  string `json:"lineUserId"`
	DosageName  string `json:"dosageName"`
	DosageAmout string `json:"dosageAmout"`
}

func (nr *NotificationRepo) GetNotificationList(t time.Time) ([]NotificationList, error) {
	formattedTime := t.Format("15:04")
	stmt, err := nr.repo.Prepare(
		`
			SELECT 
				u.line_user_id, 
				d.name AS dosage_name, 
				d.amount AS dosage_amount
			FROM 
				users u
			JOIN 
				dosage d ON u.id = d.user_id
			WHERE 
				((d.morning_flg = true AND u.morning_medication_time = $1) OR 
				(d.afternoon_flg = true AND u.afternoon_medication_time = $1) OR 
				(d.evening_flg = true AND u.evening_medication_time = $1))
			`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(formattedTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if err := rows.Err(); err != nil {
		return nil, err
	}

	var notificationList []NotificationList

	for rows.Next() {
		var notification NotificationList
		err := rows.Scan(&notification.LineUserId, &notification.DosageName, &notification.DosageAmout)
		if err != nil {
			return nil, err
		}
		notificationList = append(notificationList, notification)
	}
	return notificationList, nil
}
