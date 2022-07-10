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
	var e database.ProfileFeedback
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
	fmt.Println(*e.Email)
	if *e.Email != "" {
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
