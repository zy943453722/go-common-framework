package table

import (
	"strconv"

	"github.com/beego/beego/v2/client/orm"
	"go-common-framework/services/mysql"
	"go-common-framework/util"
)

const (
	//用户分组
	USER_GROUP_EXTERNAL     = "external"
	USER_GROUP_RD_QA        = "rd&qa"
	USER_GROUP_MANAGER      = "manager"
	USER_GROUP_OTHER        = "other"
	USER_GROUP_UNCLASSIFIED = "unclassified"
)

var (
	UserGroupMap = map[string]string{
		USER_GROUP_EXTERNAL:     "外部客户",
		USER_GROUP_RD_QA:        "测试研发",
		USER_GROUP_MANAGER:      "管理节点",
		USER_GROUP_OTHER:        "其他",
		USER_GROUP_UNCLASSIFIED: "未分组",
	}
)

type Consumer struct {
	Id         int
	Cid        string
	Name       string
	Company    string
	UserGroup  string
	CreatedAt  int64
	CreatedErp string
	UpdatedAt  int64
	UpdatedErp string
	DeletedAt  int64
	IsDel      int
}

func init() {
	orm.RegisterModel(new(Consumer))
}

func CreateConsumer(requestId string, consumerList []*Consumer) error {
	o := orm.NewOrmUsingDB("write")
	mysql.InitOrmLogFunc(requestId)

	var err error
	for _, consumer := range consumerList {
		if _, _, err = o.ReadOrCreate(consumer, "source_id", "source"); err != nil {
			return err
		}
	}
	return nil
}

func GetConsumerList(requestId string, params map[string]interface{}, limit, offset, all int) ([]*Consumer, int, error) {
	consumerList := make([]*Consumer, 0)
	count := make([]orm.Params, 0)
	mysql.InitOrmLogFunc(requestId)
	o := orm.NewOrmUsingDB("read")

	sql := "select * from " + CONSUMER_TABLE_NAME + " where is_del = 0"
	sqlCount := "select count(*) as num from " + CONSUMER_TABLE_NAME + " where is_del = 0"

	sql += buildConsumerParam(params)
	sqlCount += buildConsumerParam(params)
	_, err := o.Raw(sqlCount).Values(&count)
	if err != nil {
		return nil, 0, err
	}
	if all == ALL {
		sql += " order by created_at desc"
	} else if limit > 0 && offset >= 0 {
		sql += " order by created_at desc limit " + strconv.Itoa(offset) + "," + strconv.Itoa(limit) + ""
	}
	_, err = o.Raw(sql).QueryRows(&consumerList)
	if err != nil {
		if err == orm.ErrNoRows {
			return nil, 0, nil
		} else {
			return nil, 0, err
		}
	} else {
		if len(count) > 0 && count[0]["num"] != nil {
			num, _ := strconv.Atoi(count[0]["num"].(string))
			return consumerList, num, nil
		} else {
			return consumerList, 0, nil
		}
	}
}

func buildConsumerParam(param map[string]interface{}) string {
	if param == nil {
		return ""
	}
	sql := ""
	if v, ok := param["userGroup"]; ok && v.(string) != "" {
		sql += " and user_group='" + v.(string) + "'"
	}
	if v, ok := param["userGroup_in"]; ok && v.(string) != "" {
		userGroups := util.GetInSqlStr(v.(string))
		sql += " and user_group in (" + userGroups + ")"
	}
	if v, ok := param["name"]; ok && v.(string) != "" {
		sql += " and name like '%" + v.(string) + "%'"
	}
	if v, ok := param["company"]; ok && v.(string) != "" {
		sql += " and company like '%" + v.(string) + "%'"
	}

	return sql
}

func GetConsumerInfo(requestId, cid string) (*Consumer, error) {
	mysql.InitOrmLogFunc(requestId)
	o := orm.NewOrmUsingDB("read")
	consumerInfo := new(Consumer)
	if err := o.QueryTable(CONSUMER_TABLE_NAME).Filter("is_del", NOT_DELETED).Filter("cid", cid).One(consumerInfo); err != nil {
		if err == orm.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return consumerInfo, nil
}

func UpdateConsumer(requestId, cid string, params orm.Params) error {
	o := orm.NewOrmUsingDB("write")
	mysql.InitOrmLogFunc(requestId)
	if _, err := o.QueryTable(CONSUMER_TABLE_NAME).Filter("cid", cid).Update(params); err != nil {
		return err
	}
	return nil
}
