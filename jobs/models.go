package jobs

import "time"

type WorkOrderModel struct {
	ID            int       `db:"id"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
	BusinessName  string    `db:"business_name"`  // 业务名称，比如创建数据库
	DbName        string    `db:"db_name"`        // 创建的新数据库名称，比如test
	ApplicantUser string    `db:"applicant_user"` // 申请人名称
	Remark        string    `db:"remark"`         // 备注
	Department    int       `db:"department"`     // 条线名称
	DbType        int       `db:"db_type"`        // 数据库类型
	Usage         int       `db:"usage"`          // 使用用途 1 开发 2 测试 3 压测
	TimeLimit     int       `db:"time_limit"`     // 数据库预估使用期限 1 使用数周 2 使用数月 3 长期使用
	Status        int       `db:"status"`         // 工单的状态，1 未完成 2 完成
	Role          int       `db:"role"`           // 数据库用户权限  1 管理员  2 普通用户(默认)  3 游客
	OperationUser string    `db:"operation_user"` // 操作人名称
	FinishedAt    time.Time `db:"finished_at"`    // 工单完成时间,这里需要使用自定义时间类型LocalTime,如果使用time.Time插入数据库的时候就会变成0000-00:00:00
}
