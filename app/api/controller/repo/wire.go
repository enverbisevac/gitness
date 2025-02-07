// Copyright 2023 Harness, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package repo

import (
	"github.com/harness/gitness/app/auth/authz"
	"github.com/harness/gitness/app/services/importer"
	"github.com/harness/gitness/app/store"
	"github.com/harness/gitness/app/url"
	"github.com/harness/gitness/gitrpc"
	"github.com/harness/gitness/store/database/dbtx"
	"github.com/harness/gitness/types"
	"github.com/harness/gitness/types/check"

	"github.com/google/wire"
)

// WireSet provides a wire set for this package.
var WireSet = wire.NewSet(
	ProvideController,
)

func ProvideController(config *types.Config, tx dbtx.Transactor, urlProvider url.Provider,
	uidCheck check.PathUID, authorizer authz.Authorizer, repoStore store.RepoStore,
	spaceStore store.SpaceStore, pipelineStore store.PipelineStore,
	principalStore store.PrincipalStore, rpcClient gitrpc.Interface,
	importer *importer.Repository,
) *Controller {
	return NewController(config.Git.DefaultBranch, tx, urlProvider,
		uidCheck, authorizer, repoStore,
		spaceStore, pipelineStore, principalStore, rpcClient,
		importer)
}
