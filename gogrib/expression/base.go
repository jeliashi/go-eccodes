package expression

import "gogrib/accessor"

type GribExpression struct {
	name   string
	size   uint64
	inited bool
	GribExpressionInterface
}
type GribExpressionInterface interface {
	init_class() error
	init() error
	destroy() error
	print() error
	add_dependency(observer *accessor.GribAccessor) error
	native_type() (int, error)
	get_name() (string, error)
	evaluate_long(l int64) (int, error)
	evaluate_double(d float64) (int, error)
	evaluate_string(s string) (int, error)
}
