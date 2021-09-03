package migrate
/*
Copyright 2021 The Knative Authors

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

import (
	"context"
	"fmt"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CheckStoredVersions verifies that all storedVersions from the CRD are listed in the spec
// with storage: true. It means the CRD has been migrated and previous/unused API versions
// can be safely removed from the spec.
func CheckStoredVersions(ctx context.Context, apiextensions apiextensionsv1.ApiextensionsV1Client, gr schema.GroupResource) error {
	crdClient := apiextensions.CustomResourceDefinitions()

	crd, err := crdClient.Get(ctx, gr.String(), metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("unable to fetch crd %s - %w", gr, err)
	}

	for _, stored := range crd.Status.StoredVersions {
		for _, v := range crd.Spec.Versions {
			if !v.Storage {
				return fmt.Errorf("\"%s\" is invalid: spec.versions.storage must be true for \"%s\"", crd.Name, v.Name)
			}
		}
	}

	return nil
}

