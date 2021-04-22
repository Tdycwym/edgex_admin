package dal

import (
	"time"

	"github.com/nju-iot/edgex_admin/caller"
	"github.com/nju-iot/edgex_admin/logs"
	"github.com/nju-iot/edgex_admin/utils"
	"gorm.io/gorm"
)

const (
	Status_Follow   = 1
	Status_UnFollow = 2
)

type EdgexRelatedUser struct {
	ID           int64     `gorm:"column:id" json:"id"`
	UserID       int64     `gorm:"column:user_id" json:"user_id"`
	Username     string    `gorm:"column:username" json:"username"`
	EdgexID      int64     `gorm:"column:edgex_id" json:"edgex_id"`
	EdgexName    string    `gorm:"column:edgex_name" json:"edgex_name"`
	Status       int32     `gorm:"column:status" json:"status"`
	CreatedTime  time.Time `gorm:"column:created_time" json:"created_time"`
	ModifiedTime time.Time `gorm:"column:modified_time" json:"modified_time"`
}

// AddEdgexRelatedUser ...
func AddEdgexRelatedUser(db *gorm.DB, item *EdgexRelatedUser) error {
	dbRes := db.Debug().Model(&EdgexRelatedUser{}).Create(item)
	if dbRes.Error != nil {
		logs.Error("[AddEdgexRelatedUser] create AddEdgexRelatedUser record failed: item=%+v, err=%v", item, dbRes.Error)
		return dbRes.Error
	}
	return nil
}

// UpdateEdgexRelatedUser ...
func UpdateEdgexRelatedUser(db *gorm.DB, id int64, fieldsMap map[string]interface{}) error {
	dbRes := db.Debug().Model(&EdgexRelatedUser{}).Where("id = ?", id).Updates(fieldsMap)
	if dbRes.Error != nil {
		logs.Error("[UpdateEdgexRelatedUser] update EdgexRelatedUser failed: id=%+v, filedsMap=%+v, err=%v", id, fieldsMap, dbRes.Error)
		return dbRes.Error
	}
	return nil
}

// FindEdgexRelatedUserByUserIDAndEdgexID ...
func FindEdgexRelatedUserByUserIDAndEdgexID(edgexID int64, userID int64) (entity *EdgexRelatedUser, err error) {

	itemList := make([]*EdgexRelatedUser, 0)
	dbRes := caller.EdgexDB.Debug().Model(&EdgexRelatedUser{}).
		Where("user_id = ? AND edgex_id = ?", userID, edgexID).
		Find(&itemList)
	if dbRes.Error != nil {
		logs.Error("[FindEdgexRelatedUserByUserIDAndEdgexID] find EdgexRelatedUser failed: user_id=%v, edgex_id=%v", userID, edgexID)
		err = dbRes.Error
		return
	}
	if len(itemList) > 0 {
		entity = itemList[0]
	}
	return
}

// GetFollowEdgexIDs ...
func GetFollowEdgexIDs(userID int64) (edgexIDs []int64, err error) {
	edgexIDs = make([]int64, 0)
	itemList := make([]*EdgexRelatedUser, 0)
	dbRes := caller.EdgexDB.Debug().Model(&EdgexRelatedUser{}).
		Where("user_id = ? and status = 1", userID).
		Find(&itemList)
	if dbRes.Error != nil {
		logs.Error("[GetFollowEdgexIDs] select from database failed: userID=%v, status=1, err=%v", userID, err)
		err = dbRes.Error
		return
	}

	for _, item := range itemList {
		if item == nil {
			continue
		}
		edgexIDs = append(edgexIDs, item.EdgexID)
	}
	edgexIDs = utils.DeduplicationI64List(edgexIDs)
	return

}
