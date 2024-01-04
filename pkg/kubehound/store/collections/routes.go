package collections

type Route struct {
}

var _ Collection = (*Route)(nil) // Ensure interface compliance

func (c Route) Name() string {
	return RouteName
}

func (c Route) BatchSize() int {
	return DefaultBatchSize
}
