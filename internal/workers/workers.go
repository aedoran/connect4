package workers

// Msg mimics the message passed to job handlers.
type Msg struct{ args []interface{} }

// Args returns job arguments.
func (m *Msg) Args() []interface{} { return m.args }

// NewMsg constructs a Msg from args.
func NewMsg(args []interface{}) *Msg { return &Msg{args: args} }

// JobFunc is a handler for a job.
type JobFunc func(*Msg)

func Configure(map[string]string) {}

func Process(queue string, fn JobFunc, concurrency int) {}

func Run() {}

func Quit() {}
