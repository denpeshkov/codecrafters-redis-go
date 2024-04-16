package proto

type Type byte

const (
	String     = '+' // +<string>\r\n
	Error      = '-' // -<string>\r\n
	Integer    = ':' // :[<+|->]<value>\r\n
	BulkString = '$' // $<length>\r\n<data>\r\n
	Array      = '*' // *<number-of-elements>\r\n<element-1>...<element-n>
	Null       = '_' // _\r\n
	Boolean    = '#' // #<t|f>\r\n
	Double     = ',' // ,[<+|->]<integral>[.<fractional>][<E|e>[sign]<exponent>]\r\n
	BigNumber  = '(' // ([+|-]<number>\r\n
	BulkError  = '!' // !<length>\r\n<error>\r\n
	VerbString = '=' // =<length>\r\n<encoding>:<data>\r\n
	Map        = '%' // %<number-of-entries>\r\n<key-1><value-1>...<key-n><value-n>
	Set        = '~' // ~<number-of-elements>\r\n<element-1>...<element-n>
	Push       = '>' // ><number-of-elements>\r\n<element-1>...<element-n>
)
