// Copyright © 2019 Banzai Cloud
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

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"time"

	v1alpha1 "github.com/banzaicloud/kafka-operator/api/kafka/v1alpha1"
	scheme "github.com/banzaicloud/kafka-operator/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// KafkaUsersGetter has a method to return a KafkaUserInterface.
// A group's client should implement this interface.
type KafkaUsersGetter interface {
	KafkaUsers(namespace string) KafkaUserInterface
}

// KafkaUserInterface has methods to work with KafkaUser resources.
type KafkaUserInterface interface {
	Create(*v1alpha1.KafkaUser) (*v1alpha1.KafkaUser, error)
	Update(*v1alpha1.KafkaUser) (*v1alpha1.KafkaUser, error)
	UpdateStatus(*v1alpha1.KafkaUser) (*v1alpha1.KafkaUser, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.KafkaUser, error)
	List(opts v1.ListOptions) (*v1alpha1.KafkaUserList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.KafkaUser, err error)
	KafkaUserExpansion
}

// kafkaUsers implements KafkaUserInterface
type kafkaUsers struct {
	client rest.Interface
	ns     string
}

// newKafkaUsers returns a KafkaUsers
func newKafkaUsers(c *KafkaV1alpha1Client, namespace string) *kafkaUsers {
	return &kafkaUsers{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the kafkaUser, and returns the corresponding kafkaUser object, and an error if there is any.
func (c *kafkaUsers) Get(name string, options v1.GetOptions) (result *v1alpha1.KafkaUser, err error) {
	result = &v1alpha1.KafkaUser{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("kafkausers").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of KafkaUsers that match those selectors.
func (c *kafkaUsers) List(opts v1.ListOptions) (result *v1alpha1.KafkaUserList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.KafkaUserList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("kafkausers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested kafkaUsers.
func (c *kafkaUsers) Watch(opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("kafkausers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a kafkaUser and creates it.  Returns the server's representation of the kafkaUser, and an error, if there is any.
func (c *kafkaUsers) Create(kafkaUser *v1alpha1.KafkaUser) (result *v1alpha1.KafkaUser, err error) {
	result = &v1alpha1.KafkaUser{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("kafkausers").
		Body(kafkaUser).
		Do().
		Into(result)
	return
}

// Update takes the representation of a kafkaUser and updates it. Returns the server's representation of the kafkaUser, and an error, if there is any.
func (c *kafkaUsers) Update(kafkaUser *v1alpha1.KafkaUser) (result *v1alpha1.KafkaUser, err error) {
	result = &v1alpha1.KafkaUser{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("kafkausers").
		Name(kafkaUser.Name).
		Body(kafkaUser).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *kafkaUsers) UpdateStatus(kafkaUser *v1alpha1.KafkaUser) (result *v1alpha1.KafkaUser, err error) {
	result = &v1alpha1.KafkaUser{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("kafkausers").
		Name(kafkaUser.Name).
		SubResource("status").
		Body(kafkaUser).
		Do().
		Into(result)
	return
}

// Delete takes name of the kafkaUser and deletes it. Returns an error if one occurs.
func (c *kafkaUsers) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("kafkausers").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *kafkaUsers) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("kafkausers").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched kafkaUser.
func (c *kafkaUsers) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.KafkaUser, err error) {
	result = &v1alpha1.KafkaUser{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("kafkausers").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}