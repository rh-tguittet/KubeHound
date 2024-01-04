package store

import (
	routev1 "github.com/openshift/api/route/v1"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Route struct {
	Id           primitive.ObjectID `bson:"_id"`
	IsNamespaced bool               `bson:"is_namespaced"`
	Namespace    string             `bson:"namespace"`
	Name         string             `bson:"name"`
	K8           routev1.Route      `bson:"k8"`
	Ownership    OwnershipInfo      `bson:"ownership"`
	Runtime      RuntimeInfo        `bson:"runtime"`
	// TODO[TG] Maybe we need to add more things here?
}
