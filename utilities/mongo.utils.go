package utilities

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//InterfaceToBsonM func is to convert an interface to bson.M
func InterfaceToBsonM(val interface{}) (bson.M, error) {
	var result bson.M
	byteInterface, err := bson.Marshal(val)
	if err != nil {
		return nil, err
	}
	if err := bson.Unmarshal(byteInterface, &result); err != nil {
		return nil, err
	}
	return result, nil
}
func ConvertStringIdToObjectID(id string) (*primitive.ObjectID, error) {
	if ok := primitive.IsValidObjectID(id); !ok {
		return nil, errors.New("invalid id")
	}
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	return &objectId, nil

}
