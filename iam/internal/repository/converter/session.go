package converter

import "github.com/google/uuid"

func ValueToUserId(value []byte) (uuid.UUID, error) {
	userID, err := uuid.Parse(string(value))

	if err != nil {
		return uuid.Nil, err
	}
	return userID, nil
}
