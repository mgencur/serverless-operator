/*
Copyright 2020 The Knative Authors

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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
	v1alpha1 "knative.dev/eventing/pkg/apis/configs/v1alpha1"
)

// FakeConfigMapPropagations implements ConfigMapPropagationInterface
type FakeConfigMapPropagations struct {
	Fake *FakeConfigsV1alpha1
	ns   string
}

var configmappropagationsResource = schema.GroupVersionResource{Group: "configs.internal.knative.dev", Version: "v1alpha1", Resource: "configmappropagations"}

var configmappropagationsKind = schema.GroupVersionKind{Group: "configs.internal.knative.dev", Version: "v1alpha1", Kind: "ConfigMapPropagation"}

// Get takes name of the configMapPropagation, and returns the corresponding configMapPropagation object, and an error if there is any.
func (c *FakeConfigMapPropagations) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.ConfigMapPropagation, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(configmappropagationsResource, c.ns, name), &v1alpha1.ConfigMapPropagation{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ConfigMapPropagation), err
}

// List takes label and field selectors, and returns the list of ConfigMapPropagations that match those selectors.
func (c *FakeConfigMapPropagations) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.ConfigMapPropagationList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(configmappropagationsResource, configmappropagationsKind, c.ns, opts), &v1alpha1.ConfigMapPropagationList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.ConfigMapPropagationList{ListMeta: obj.(*v1alpha1.ConfigMapPropagationList).ListMeta}
	for _, item := range obj.(*v1alpha1.ConfigMapPropagationList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested configMapPropagations.
func (c *FakeConfigMapPropagations) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(configmappropagationsResource, c.ns, opts))

}

// Create takes the representation of a configMapPropagation and creates it.  Returns the server's representation of the configMapPropagation, and an error, if there is any.
func (c *FakeConfigMapPropagations) Create(ctx context.Context, configMapPropagation *v1alpha1.ConfigMapPropagation, opts v1.CreateOptions) (result *v1alpha1.ConfigMapPropagation, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(configmappropagationsResource, c.ns, configMapPropagation), &v1alpha1.ConfigMapPropagation{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ConfigMapPropagation), err
}

// Update takes the representation of a configMapPropagation and updates it. Returns the server's representation of the configMapPropagation, and an error, if there is any.
func (c *FakeConfigMapPropagations) Update(ctx context.Context, configMapPropagation *v1alpha1.ConfigMapPropagation, opts v1.UpdateOptions) (result *v1alpha1.ConfigMapPropagation, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(configmappropagationsResource, c.ns, configMapPropagation), &v1alpha1.ConfigMapPropagation{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ConfigMapPropagation), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeConfigMapPropagations) UpdateStatus(ctx context.Context, configMapPropagation *v1alpha1.ConfigMapPropagation, opts v1.UpdateOptions) (*v1alpha1.ConfigMapPropagation, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(configmappropagationsResource, "status", c.ns, configMapPropagation), &v1alpha1.ConfigMapPropagation{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ConfigMapPropagation), err
}

// Delete takes name of the configMapPropagation and deletes it. Returns an error if one occurs.
func (c *FakeConfigMapPropagations) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(configmappropagationsResource, c.ns, name), &v1alpha1.ConfigMapPropagation{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeConfigMapPropagations) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(configmappropagationsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.ConfigMapPropagationList{})
	return err
}

// Patch applies the patch and returns the patched configMapPropagation.
func (c *FakeConfigMapPropagations) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.ConfigMapPropagation, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(configmappropagationsResource, c.ns, name, pt, data, subresources...), &v1alpha1.ConfigMapPropagation{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ConfigMapPropagation), err
}