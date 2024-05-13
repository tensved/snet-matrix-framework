package lib

import "errors"

type MaxSizeConstraint struct {
	MaxSizeBytes int
}

func (c *MaxSizeConstraint) Check(value interface{}) error {

	val, ok := value.(string)
	if !ok {
		return errors.New("value is not a string")
	}

	if len(val) > c.MaxSizeBytes {
		return errors.New("value is too big")
	}

	return nil
}

type MinSizeConstraint struct {
	MinSizeBytes int
}

func (c *MinSizeConstraint) Check(value interface{}) error {

	val, ok := value.([]byte)
	if !ok || len(val) < c.MinSizeBytes {
		return errors.New("value is too small")
	}

	return nil
}
