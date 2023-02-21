// Copyright (C) 2023  Tricorder Observability
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/tricorder/src/api-server/pb"
	"github.com/tricorder/src/pb/module"
)

func (*Deployer) Add(context.Context, *module.Module) (*pb.ModuleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Add not implemented")
}

func (*Deployer) Delete(context.Context, *pb.ModuleRequest) (*pb.ModuleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}

func (*Deployer) List(context.Context, *pb.ListQuery) (*pb.ModuleListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}

func (*Deployer) Deploy(context.Context, *pb.ModuleRequest) (*pb.ModuleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Deploy not implemented")
}

func (*Deployer) Undeploy(context.Context, *pb.ModuleRequest) (*pb.ModuleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Undeploy not implemented")
}
