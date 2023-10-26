/*
 * Copyright 2018-2020 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package mslibpak

import (
	"github.com/buildpacks/libcnb"

	"github.com/paketo-buildpacks/libpak/bard"
)

type buildDelegate struct {
	delegate libcnb.Builder
}

func (b buildDelegate) Build(context libcnb.BuildContext) (libcnb.BuildResult, error) {
	result, err := b.delegate.Build(context)
	if err != nil {
		err = bard.IdentifiableError{
			Name:        context.Buildpack.Info.Name,
			Description: context.Buildpack.Info.Version,
			Err:         err,
		}
	}

	return result, err
}
