package db

type _UpdateSpec struct {
	Set interface{} `bson:"$set"`
}
