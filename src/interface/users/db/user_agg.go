package usersdb

import "go.mongodb.org/mongo-driver/bson"

// pass 0 into limit to remove the limit
func UserAgg(preFilter bson.M, postFilter bson.M, limit int64) bson.A {
	limitFilter := bson.M{"$match": bson.M{}}
	if limit > 0 {
		limitFilter = bson.M{"$limit": limit}
	}
	return bson.A{
		bson.M{"$match": preFilter},
		bson.M{
			"$lookup": bson.M{
				"from":         "mediaFiles",
				"localField":   "imageId",
				"foreignField": "id",
				"as":           "image",
			},
		},
		bson.M{
			"$addFields": bson.M{
				"image": bson.M{
					"$arrayElemAt": bson.A{"$image", 0},
				},
			},
		},
		bson.M{"$match": postFilter},
		limitFilter,
	}
}
