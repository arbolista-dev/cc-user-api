package ds

import (
  "github.com/arbolista-dev/cc-user-api/app/models"
  "time"
  "upper.io/db.v2"
  "encoding/json"
)

func UpdateActions(userID uint, update models.ActionUpdate) (err error) {
	var user models.User
	err = userSource.Find(db.Cond{"user_id": userID}).One(&user)
	if err != nil {
		return
	}

  var userActionExists bool
  var action models.Action
  err = actionSource.Find(db.Cond{"user_id": userID}, db.Cond{"key": update.Key}).One(&action)
  if err != nil {
    userActionExists = false
  } else {
    userActionExists = true
  }

  // Details only exist if status == 'pledged' || 'completed'
  var details map[string]interface{}
  if update.Details != nil {
    err = json.Unmarshal([]byte(update.Details), &details)
    if err != nil {
      details = nil
    }
  }

  if !userActionExists {

    // user action needs to be newly created!
    action.Key = update.Key
    action.Status = update.Status
    action.UserID = userID
    action.CreatedAt = time.Now()

    if status := update.Status; status == "pledged" || status == "completed" {
      action.TonsSaved = details["tons_saved"].(float64)
      action.DollarsSaved = details["dollars_saved"].(float64)
      action.UpfrontCost = details["upfront_cost"].(float64)

      actionSource.Insert(action)
      return
    } else if status == "not_relevant" {

      actionSource.Insert(action)
      return
    }
    return

  } else {

    switch status := update.Status; status {
      case "pledged":
        action.Status = update.Status
        action.TonsSaved = details["tons_saved"].(float64)
        action.DollarsSaved = details["dollars_saved"].(float64)
        action.UpfrontCost = details["upfront_cost"].(float64)

        err = actionSource.Find(db.Cond{"user_id": userID}, db.Cond{"key": update.Key}).Update(action)
        return

      case "completed":
        action.Status = update.Status
        action.TonsSaved = details["tons_saved"].(float64)
        action.DollarsSaved = details["dollars_saved"].(float64)
        action.UpfrontCost = details["upfront_cost"].(float64)

        err = actionSource.Find(db.Cond{"user_id": userID}, db.Cond{"key": update.Key}).Update(action)
        return
      case "unpledged":
        if action.Status == "pledged" {
          err = actionSource.Find(db.Cond{"user_id": userID}, db.Cond{"key": update.Key}).Delete()
          return
        }
        return
      case "not_relevant":
        action.Status = update.Status
        action.TonsSaved = 0
        action.DollarsSaved = 0
        action.UpfrontCost = 0

        err = actionSource.Find(db.Cond{"user_id": userID}, db.Cond{"key": update.Key}).Update(action)
      case "relevant":
        if action.Status == "not_relevant" {
          err = actionSource.Find(db.Cond{"user_id": userID}, db.Cond{"key": update.Key}).Delete()
          return
        }
        return
      case "uncompleted":
        if action.Status == "completed" {
          action.Status = update.Status
          err = actionSource.Find(db.Cond{"user_id": userID}, db.Cond{"key": update.Key}).Update(action)
          return
        }
        return
    }
  }
	return
}

func DeleteUserActions(userID uint) (err error) {
  err = actionSource.Find(db.Cond{"user_id": userID}).Delete()
  return
}
