package service

import (
	"context"

	"github.com/jscastelblancoh/users_service/user/pkg/db"
	"github.com/jscastelblancoh/users_service/user/pkg/io"
	"gopkg.in/mgo.v2/bson"
)

// UserService describes the service.
type UserService interface {
	GetbyId(ctx context.Context, id string) (t []io.User, err error)
	Post(ctx context.Context, statistic io.User) (t io.User, err error)
	Delete(ctx context.Context, id string) (err error)
	Put(ctx context.Context, id string, user io.User) (error error)
}

type basicUserService struct{}

func (b *basicUserService) GetbyId(ctx context.Context, id string) (t []io.User, err error) {
	session, err2 := db.GetMongoSession()
	if err2 != nil {
		return t, err2
	}
	defer session.Close()
	c := session.DB("user_service").C("users")
	err = c.Find(bson.M{"_id": bson.ObjectIdHex(id)}).All(&t)
	return t, err
}
func (b *basicUserService) Post(ctx context.Context, statistic io.User) (t io.User, err error) {
	statistic.ID = bson.NewObjectId()
	session, err := db.GetMongoSession()
	if err != nil {
		return t, err
	}
	defer session.Close()
	c := session.DB("user_service").C("users")
	err = c.Insert(&statistic)
	return statistic, err
}
func (b *basicUserService) Delete(ctx context.Context, id string) (err error) {
	session, err := db.GetMongoSession()
	if err != nil {
		return err
	}
	defer session.Close()
	c := session.DB("user_service").C("users")
	return c.Remove(bson.M{"_id": bson.ObjectIdHex(id)})
}
func (b *basicUserService) Put(ctx context.Context, id string, user io.User) (error error) {
	session, err := db.GetMongoSession()
	if err != nil {
		return err
	}
	defer session.Close()
	c := session.DB("user_service").C("users")
	c.Update(bson.M{"_id": bson.ObjectIdHex(id)}, bson.M{"$set": bson.M{"name": user.Name}})
	c.Update(bson.M{"_id": bson.ObjectIdHex(id)}, bson.M{"$set": bson.M{"cedula": user.Cedula}})
	c.Update(bson.M{"_id": bson.ObjectIdHex(id)}, bson.M{"$set": bson.M{"lunchroom_id": user.LunchroomID}})
	return c.Update(bson.M{"_id": bson.ObjectIdHex(id)}, bson.M{"$set": bson.M{"active_ticket": user.ActiveTicket}})
}

// NewBasicUserService returns a naive, stateless implementation of UserService.
func NewBasicUserService() UserService {
	return &basicUserService{}
}

// New returns a UserService with all of the expected middleware wired in.
func New(middleware []Middleware) UserService {
	var svc UserService = NewBasicUserService()
	for _, m := range middleware {
		svc = m(svc)
	}
	return svc
}
