package models

const (
	//导出文件前缀
	CONSUMER_FILE_PREFIX = "consumer_list_"

	//默认pageNumber, pageSize
	DEFAULT_PAGENUMBER = 1
	DEFAULT_PAGESIZE   = 20
)

var (
	ConsumerHeader = []string{
		"客户ID",
		"客户名称",
		"企业名称",
		"客户分组",
		"更新时间",
		"更新人",
	}
)
