package dal

import (
	"time"

	"github.com/tdycwym/edgex_admin/caller"
	"github.com/tdycwym/edgex_admin/logs"
	"gorm.io/gorm"
)

// EdgexServiceItem ...
type EdgexServiceItem struct {
	ID           int64     `gorm:"column:id" json:"id"`
	UserID       int64     `gorm:"column:user_id" json:"user_id"`
	EdgexName    string    `gorm:"column:edgex_name" json:"edgex_name"`
	Prefix       string    `gorm:"column:prefix" json:"prefix"`
	Status       int32     `gorm:"column:status" json:"status"`
	Deleted      int32     `gorm:"column:deleted" json:"deleted"`
	Address      string    `gorm:"column:address" json:"address"`
	CreatedTime  time.Time `gorm:"column:created_time" json:"created_time"`
	ModifiedTime time.Time `gorm:"column:modified_time" json:"modified_time"`
	Description  string    `gorm:"column:description" json:"description"`
	Location     string    `gorm:"column:location" json:"location"`
	Extra        string    `gorm:"column:extra" json:"extra"`
}

// AddEdgex ...
func AddEdgex(db *gorm.DB, edgex *EdgexServiceItem) error {
	dbRes := db.Debug().Model(&EdgexServiceItem{}).Create(edgex)
	if dbRes.Error != nil {
		logs.Error("[AddEdgex] create edgex failed: edgex=%+v, err=%v", edgex, dbRes.Error)
		return dbRes.Error
	}
	return nil
}

// UpdateEdgex ...
func UpdateEdgex(db *gorm.DB, edgexID int64, fieldsMap map[string]interface{}) error {
	dbRes := db.Debug().Model(&EdgexServiceItem{}).Where("id = ?", edgexID).Updates(fieldsMap)
	if dbRes.Error != nil {
		logs.Error("[UpdateEdgex] update edgex failed: edgexID=%+v, filedsMap=%+v, err=%v", edgexID, fieldsMap, dbRes.Error)
		return dbRes.Error
	}
	return nil
}

// GetEdgexList ...
func GetEdgexList(edgexIDs []int64, userIDs []int64, keyword string, status int, offset int, count int) (edgexList []*EdgexServiceItem, err error) {
	edgexList = make([]*EdgexServiceItem, 0)

	db := caller.EdgexDB.Debug().Model(&EdgexServiceItem{}).Where("deleted = 0")
	if len(edgexIDs) > 0 {
		db = db.Where("id IN (?)", edgexIDs)
	}

	if len(userIDs) > 0 {
		db = db.Where("user_id IN (?)", userIDs)
	}

	if keyword != "" {
		likeKey := "%" + keyword + "%"
		db = db.Where("edgex_name LIKE ? OR prefix LIKE ?", likeKey, likeKey)
	}

	statusList := getStatusList(status)
	dbRes := db.Where("status in (?)", statusList).Offset(offset).Limit(count).Find(&edgexList)
	if dbRes.Error != nil {
		logs.Error("[GetEdgexList] get edgexList failed: edgexIDs=%+v, userIDs=%+v, key=%v, status=%+v, offset=%v, count=%v",
			edgexIDs, userIDs, keyword, statusList, offset, count)
		err = dbRes.Error
		return
	}
	return
}

func getStatusList(status int) []int {
	switch status {
	case 1:
		return []int{0}
	case 2:
		return []int{1}
	}
	return []int{0, 1}
}
