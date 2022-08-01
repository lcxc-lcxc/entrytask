/**
 @author: 15973
 @date: 2022/07/08
 @note:
**/
package constant

var (
	Success = RetCode{0, "success"}

	ServerError   = RetCode{10000000, "Server Error"}
	InvalidParams = RetCode{10000001, "Invalid Params"}
	NotFound      = RetCode{10000002, "Not Found"}

	UserLoginFailed    = RetCode{20010001, "User Login Failed"}
	UserLoginRequired  = RetCode{20010002, "User Login Required"}
	UserRegisterFailed = RetCode{20010003, "User Register Failed"}
)

var RetcodeMap = map[int]RetCode{
	0:        Success,
	10000000: ServerError,
	10000001: InvalidParams,
	10000002: NotFound,
	20010001: UserLoginFailed,
	20010002: UserLoginRequired,
	20010003: UserRegisterFailed,
}

type RetCode struct {
	retcode int
	msg     string
}

func (r RetCode) GetRetCode() int {
	return r.retcode
}

func (r RetCode) GetMsg() string {
	return r.msg
}
