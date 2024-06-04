package sessionsdb

import "go.mongodb.org/mongo-driver/bson"

func SessionAgg(preFilter bson.M, postFilter bson.M, limit int64) bson.A {
	limitFilter := bson.M{"$match": bson.M{}}
	if limit > 0 {
		limitFilter = bson.M{"$limit": limit}
	}
	return bson.A{
		bson.M{"$match": preFilter},
		bson.M{"$lookup": bson.M{
			"from":         "users",
			"localField":   "userId",
			"foreignField": "id",
			"as":           "user",
		}},
		bson.M{"$addFields": bson.M{
			"user": bson.M{
				"$arrayElemAt": bson.A{"$user", 0},
			},
		}},
		bson.M{"$lookup": bson.M{
			"from":         "mediaFiles",
			"localField":   "user.imageId",
			"foreignField": "id",
			"as":           "user.image",
		}},
		bson.M{"$addFields": bson.M{
			"user.image": bson.M{
				"$arrayElemAt": bson.A{"$user.image", 0},
			},
		}},
		bson.M{"$match": postFilter},
		limitFilter,
	}
}
