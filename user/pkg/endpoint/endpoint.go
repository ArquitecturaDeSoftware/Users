package endpoint

import (
	"context"

	endpoint "github.com/go-kit/kit/endpoint"
	io "github.com/jscastelblancoh/users_service/user/pkg/io"
	service "github.com/jscastelblancoh/users_service/user/pkg/service"
)

// GetbyIdRequest collects the request parameters for the GetbyId method.
type GetbyIdRequest struct {
	Id string `json:"id"`
}

// GetbyIdResponse collects the response parameters for the GetbyId method.
type GetbyIdResponse struct {
	T   []io.User `json:"t"`
	Err error     `json:"err"`
}

// MakeGetbyIdEndpoint returns an endpoint that invokes GetbyId on the service.
func MakeGetbyIdEndpoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetbyIdRequest)
		t, err := s.GetbyId(ctx, req.Id)
		return GetbyIdResponse{
			Err: err,
			T:   t,
		}, nil
	}
}

// Failed implements Failer.
func (r GetbyIdResponse) Failed() error {
	return r.Err
}

// PostRequest collects the request parameters for the Post method.
type PostRequest struct {
	Statistic io.User `json:"statistic"`
}

// PostResponse collects the response parameters for the Post method.
type PostResponse struct {
	T   io.User `json:"t"`
	Err error   `json:"err"`
}

// MakePostEndpoint returns an endpoint that invokes Post on the service.
func MakePostEndpoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(PostRequest)
		t, err := s.Post(ctx, req.Statistic)
		return PostResponse{
			Err: err,
			T:   t,
		}, nil
	}
}

// Failed implements Failer.
func (r PostResponse) Failed() error {
	return r.Err
}

// DeleteRequest collects the request parameters for the Delete method.
type DeleteRequest struct {
	Id string `json:"id"`
}

// DeleteResponse collects the response parameters for the Delete method.
type DeleteResponse struct {
	Err error `json:"err"`
}

// MakeDeleteEndpoint returns an endpoint that invokes Delete on the service.
func MakeDeleteEndpoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteRequest)
		err := s.Delete(ctx, req.Id)
		return DeleteResponse{Err: err}, nil
	}
}

// Failed implements Failer.
func (r DeleteResponse) Failed() error {
	return r.Err
}

// Failer is an interface that should be implemented by response types.
// Response encoders can check if responses are Failer, and if so they've
// failed, and if so encode them using a separate write path based on the error.
type Failure interface {
	Failed() error
}

// GetbyId implements Service. Primarily useful in a client.
func (e Endpoints) GetbyId(ctx context.Context, id string) (t []io.User, err error) {
	request := GetbyIdRequest{Id: id}
	response, err := e.GetbyIdEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(GetbyIdResponse).T, response.(GetbyIdResponse).Err
}

// Post implements Service. Primarily useful in a client.
func (e Endpoints) Post(ctx context.Context, statistic io.User) (t io.User, err error) {
	request := PostRequest{Statistic: statistic}
	response, err := e.PostEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(PostResponse).T, response.(PostResponse).Err
}

// Delete implements Service. Primarily useful in a client.
func (e Endpoints) Delete(ctx context.Context, id string) (err error) {
	request := DeleteRequest{Id: id}
	response, err := e.DeleteEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(DeleteResponse).Err
}

// PutRequest collects the request parameters for the Put method.
type PutRequest struct {
	Id   string  `json:"id"`
	User io.User `json:"user"`
}

// PutResponse collects the response parameters for the Put method.
type PutResponse struct {
	Error error `json:"error"`
}

// MakePutEndpoint returns an endpoint that invokes Put on the service.
func MakePutEndpoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(PutRequest)
		error := s.Put(ctx, req.Id, req.User)
		return PutResponse{Error: error}, nil
	}
}

// Failed implements Failer.
func (r PutResponse) Failed() error {
	return r.Error
}

// Put implements Service. Primarily useful in a client.
func (e Endpoints) Put(ctx context.Context, id string, user io.User) (error error) {
	request := PutRequest{
		Id:   id,
		User: user,
	}
	response, err := e.PutEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(PutResponse).Error
}
