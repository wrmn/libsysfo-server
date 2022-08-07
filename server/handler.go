package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"libsysfo-server/database"
	"libsysfo-server/utility/email"
	"net/http"
)

func newFeedback(w http.ResponseWriter, r *http.Request) {
	var e database.Feedback
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&e)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			badRequest(w, "Wrong Type provided for field "+unmarshalErr.Field)
		} else {
			badRequest(w, err.Error())
		}
		return
	}
	database.DB.Save(&e)
	if e.Email != nil {
		content := fmt.Sprintf("<html><head></head><body><p>Hello %s,</p>Terima kasih untuk feedback anda untuk aplikasi libsysfo</body>	</html>",
			e.Name,
		)

		emailBody := email.Content{
			Subject:     "Feedback received",
			HtmlContent: content,
		}

		receiver := email.ToData{
			Name:  "Libsysfo user",
			Email: *e.Email,
		}

		err = emailBody.SendEmail(receiver)
		if err != nil {
			fmt.Println("feedback email invalid")
		}
	}

	response{
		Status:      http.StatusOK,
		Reason:      "Ok",
		Description: "Feedback send",
	}.responseFormatter(w)
}

func getNotification(w http.ResponseWriter, r *http.Request) {
	data, invalid := checkToken(r, w)
	if invalid {
		return
	}

	notification := []database.Notification{}
	db := database.DB.Where("user_id = ?", data.ID).Order("created_at DESC").Find(&notification)
	if invalid := databaseException(w, db); invalid {
		return
	}

	response{
		Data: responseBody{
			Notification: notification,
		},
		Status:      http.StatusOK,
		Reason:      "Ok",
		Description: "success",
	}.responseFormatter(w)
}

func readNotification(w http.ResponseWriter, r *http.Request) {
	data, invalid := checkToken(r, w)
	if invalid {
		return
	}

	notification := []database.Notification{}
	db := database.DB.Where("user_id = ?", data.ID).Order("created_at DESC").Find(&notification)
	if invalid := databaseException(w, db); invalid {
		return
	}

	for _, n := range notification {
		n.Read = true
		database.DB.Save(&n)
	}

	response{
		Status:      http.StatusOK,
		Reason:      "Ok",
		Description: "success",
	}.responseFormatter(w)
}
