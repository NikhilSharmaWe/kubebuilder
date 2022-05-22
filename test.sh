#!/usr/bin/env bash

# Copyright 2018 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# prow calls this file currently, but we can just use `make test` to test
# the set of things we want.
CHECK_DOCS_ONLY=$(make test-docs-only)
if [ -z "$CHECK_DOCS_ONLY" ]; then
    make test
    exit 0
fi
echo "WARNING: The tests were skipped because only changes on the docs were faced"
