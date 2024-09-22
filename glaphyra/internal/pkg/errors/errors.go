package errors

type ObjectNotFoundError struct {
	// TODO хз че добавить мб сделать отдельно на уровне каждого пакета
}

func (e *ObjectNotFoundError) Error() string {
	return "Object not found"
}
