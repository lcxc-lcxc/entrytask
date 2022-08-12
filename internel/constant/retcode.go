package constant

var (
	Success = RetCode{0, "success"}

	ServerError   = RetCode{10000000, "Server Error"}
	InvalidParams = RetCode{10000001, "Invalid Params"}
	NotFound      = RetCode{10000002, "Not Found"}

	SessionError       = RetCode{20000000, "Session Error"}
	UserLoginFailed    = RetCode{20000001, "User Login Failed"}
	UserLoginRequired  = RetCode{20000002, "User Login Required"}
	UserRegisterFailed = RetCode{20000003, "User Register Failed"}

	ProductListGetFailed   = RetCode{30000001, "Product List Get Failed"}
	ProductSearchFailed    = RetCode{30000002, "Product Search Get Failed"}
	ProductDetailGetFailed = RetCode{30000003, "Product Detail Failed"}

	CommentDetailGetFailed   = RetCode{40000001, "Comment Detail Get Failed"}
	CommentCreateFailed      = RetCode{40000002, "Comment Create Failed"}
	CommentReplyCreateFailed = RetCode{40001001, "Comment Reply Create Failed"}
)

var RetcodeMap = map[int]RetCode{
	0:        Success,
	10000000: ServerError,
	10000001: InvalidParams,
	10000002: NotFound,
	20000000: SessionError,
	20000001: UserLoginFailed,
	20000002: UserLoginRequired,
	20000003: UserRegisterFailed,
	30000001: ProductListGetFailed,
	30000002: ProductSearchFailed,
	30000003: ProductDetailGetFailed,
	40000001: CommentDetailGetFailed,
	40000002: CommentCreateFailed,
	40001001: CommentReplyCreateFailed,
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
