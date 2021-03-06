/*
Copyright Red Hat, Inc.

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

package internalversion

import (
	"time"

	operators "github.com/operator-framework/operator-lifecycle-manager/pkg/api/apis/operators"
	scheme "github.com/operator-framework/operator-lifecycle-manager/pkg/api/client/clientset/internalversion/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// OperatorGroupsGetter has a method to return a OperatorGroupInterface.
// A group's client should implement this interface.
type OperatorGroupsGetter interface {
	OperatorGroups(namespace string) OperatorGroupInterface
}

// OperatorGroupInterface has methods to work with OperatorGroup resources.
type OperatorGroupInterface interface {
	Create(*operators.OperatorGroup) (*operators.OperatorGroup, error)
	Update(*operators.OperatorGroup) (*operators.OperatorGroup, error)
	UpdateStatus(*operators.OperatorGroup) (*operators.OperatorGroup, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*operators.OperatorGroup, error)
	List(opts v1.ListOptions) (*operators.OperatorGroupList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *operators.OperatorGroup, err error)
	OperatorGroupExpansion
}

// operatorGroups implements OperatorGroupInterface
type operatorGroups struct {
	client rest.Interface
	ns     string
}

// newOperatorGroups returns a OperatorGroups
func newOperatorGroups(c *OperatorsClient, namespace string) *operatorGroups {
	return &operatorGroups{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the operatorGroup, and returns the corresponding operatorGroup object, and an error if there is any.
func (c *operatorGroups) Get(name string, options v1.GetOptions) (result *operators.OperatorGroup, err error) {
	result = &operators.OperatorGroup{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("operatorgroups").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of OperatorGroups that match those selectors.
func (c *operatorGroups) List(opts v1.ListOptions) (result *operators.OperatorGroupList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &operators.OperatorGroupList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("operatorgroups").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested operatorGroups.
func (c *operatorGroups) Watch(opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("operatorgroups").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a operatorGroup and creates it.  Returns the server's representation of the operatorGroup, and an error, if there is any.
func (c *operatorGroups) Create(operatorGroup *operators.OperatorGroup) (result *operators.OperatorGroup, err error) {
	result = &operators.OperatorGroup{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("operatorgroups").
		Body(operatorGroup).
		Do().
		Into(result)
	return
}

// Update takes the representation of a operatorGroup and updates it. Returns the server's representation of the operatorGroup, and an error, if there is any.
func (c *operatorGroups) Update(operatorGroup *operators.OperatorGroup) (result *operators.OperatorGroup, err error) {
	result = &operators.OperatorGroup{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("operatorgroups").
		Name(operatorGroup.Name).
		Body(operatorGroup).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *operatorGroups) UpdateStatus(operatorGroup *operators.OperatorGroup) (result *operators.OperatorGroup, err error) {
	result = &operators.OperatorGroup{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("operatorgroups").
		Name(operatorGroup.Name).
		SubResource("status").
		Body(operatorGroup).
		Do().
		Into(result)
	return
}

// Delete takes name of the operatorGroup and deletes it. Returns an error if one occurs.
func (c *operatorGroups) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("operatorgroups").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *operatorGroups) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("operatorgroups").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched operatorGroup.
func (c *operatorGroups) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *operators.OperatorGroup, err error) {
	result = &operators.OperatorGroup{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("operatorgroups").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
