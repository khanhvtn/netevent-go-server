package scalars

import (
	"fmt"
	"io"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Lets redefine the base ID type to use an id from an external library
func MarshalID(id primitive.ObjectID) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, strconv.Quote(id.Hex()))
	})
}

// And the same for the unmarshaler
func UnmarshalID(v interface{}) (primitive.ObjectID, error) {
	if _, ok := v.(string); !ok {
		return primitive.NilObjectID, fmt.Errorf("ID must be strings")
	}
	objectID, err := primitive.ObjectIDFromHex(v.(string))
	return objectID, err
}
