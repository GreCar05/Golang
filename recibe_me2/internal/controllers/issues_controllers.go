package controllers

import (
	"encoding/json"
	"net/http"
	"recibe_me/configs/constants"
	"recibe_me/internal/helpers"
	"recibe_me/internal/models"
)

// IssueAdd adds a Issue
func IssueAdd(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var issueData models.IssueModel
	err := decoder.Decode(&issueData)

	if err != nil {
		helpers.Response(w, http.StatusInternalServerError, constants.ERR_DECODE, err, nil)
		return
	}

	defer r.Body.Close()

	err = helpers.IssuesCollection.Insert(issueData)

	if err != nil {
		helpers.Response(w, http.StatusBadRequest, constants.ERR_INSERT_DATA, err, nil)
		return
	}

	helpers.Response(w, http.StatusOK, constants.SUCCESS, err, nil)
}
