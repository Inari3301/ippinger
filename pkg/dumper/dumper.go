package dumper

type Dumper interface {
	Dump(v any) error
	FromDump(v any) error
}
