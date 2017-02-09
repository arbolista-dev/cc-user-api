package ds

import (
  "github.com/arbolista-dev/cc-user-api/app/models"
  "time"
  "upper.io/db.v2"
  "encoding/json"
)

func UpdateUserGoals(userID uint, update models.UserGoalUpdate) (err error) {
	var user models.User
	err = userSource.Find(db.Cond{"user_id": userID}).One(&user)
	if err != nil {
		return
	}

  var userGoalExists bool
  var userGoal models.UserGoal
  err = userGoalSource.Find(db.Cond{"user_id": userID}, db.Cond{"key": update.Key}).One(&userGoal)
  if err != nil {
    userGoalExists = false
  } else {
    userGoalExists = true
  }

  // Details only exist if status == 'pledged' || 'completed'
  var details map[string]interface{}
  if update.Details != nil {
    err = json.Unmarshal([]byte(update.Details), &details)
    if err != nil {
      details = nil
    }
  }

  if !userGoalExists {

    // user goal needs to be newly created!
    userGoal.Key = update.Key
    userGoal.Status = update.Status
    userGoal.UserID = userID
    userGoal.CreatedAt = time.Now()

    if status := update.Status; status == "pledged" || status == "completed" {
      userGoal.TonsSaved = details["tons_saved"].(float64)
      userGoal.DollarsSaved = details["dollars_saved"].(float64)
      userGoal.UpfrontCost = details["upfront_cost"].(float64)

      userGoalSource.Insert(userGoal)
      return
    } else if status == "not_relevant" {

      userGoalSource.Insert(userGoal)
      return
    }
    return

  } else {

    switch status := update.Status; status {
      case "pledged":
        userGoal.Status = update.Status
        userGoal.TonsSaved = details["tons_saved"].(float64)
        userGoal.DollarsSaved = details["dollars_saved"].(float64)
        userGoal.UpfrontCost = details["upfront_cost"].(float64)

        err = userGoalSource.Find(db.Cond{"user_id": userID}, db.Cond{"key": update.Key}).Update(userGoal)
        return

      case "completed":
        userGoal.Status = update.Status
        userGoal.TonsSaved = details["tons_saved"].(float64)
        userGoal.DollarsSaved = details["dollars_saved"].(float64)
        userGoal.UpfrontCost = details["upfront_cost"].(float64)

        err = userGoalSource.Find(db.Cond{"user_id": userID}, db.Cond{"key": update.Key}).Update(userGoal)
        return
      case "unpledged":
        if userGoal.Status == "pledged" {
          err = userGoalSource.Find(db.Cond{"user_id": userID}, db.Cond{"key": update.Key}).Delete()
          return
        }
        return
      case "not_relevant":
        userGoal.Status = update.Status
        userGoal.TonsSaved = 0
        userGoal.DollarsSaved = 0
        userGoal.UpfrontCost = 0

        err = userGoalSource.Find(db.Cond{"user_id": userID}, db.Cond{"key": update.Key}).Update(userGoal)
      case "relevant":
        if userGoal.Status == "not_relevant" {
          err = userGoalSource.Find(db.Cond{"user_id": userID}, db.Cond{"key": update.Key}).Delete()
          return
        }
        return
      case "uncompleted":
        if userGoal.Status == "completed" {
          userGoal.Status = update.Status
          err = userGoalSource.Find(db.Cond{"user_id": userID}, db.Cond{"key": update.Key}).Update(userGoal)
          return
        }
        return
    }
  }
	return
}

func RetrieveUserGoals(userID uint) (userGoals models.UserGoalsList, err error) {
  query = userGoalSource.Find(db.Cond{"user_id": userID})
  count, err := query.Count()
  if err != nil {
    return
  }
  userGoals.TotalCount = count

  err = query.All(&userGoals.List)
  if err != nil {
    return
  }
  return
}

func DeleteUserGoals(userID uint) (err error) {
  err = userGoalSource.Find(db.Cond{"user_id": userID}).Delete()
  return
}
