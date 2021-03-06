/*
Copyright ArxanFintech Technology Ltd. 2017-2018 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

                 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package wallet

// SNRequest struct for create SN
type SNRequest struct {
	Id          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Owner       string `json:"owner,omitempty"`
	ActiveCount int    `json:"active_count,omitempty"`
	ExpireTime  int64  `json:"expire_time,omitempty"`
}

// SNMetadata service number metadata struct definition
type SNMetadata struct {
	ActivedCount   int `json:"actived_count,omitempty"`
	MaxActiveCount int `json:"max_active_count,omitempty"`
}
