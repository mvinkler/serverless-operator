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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/discovery"
	fakediscovery "k8s.io/client-go/discovery/fake"
	"k8s.io/client-go/testing"
	clientset "knative.dev/eventing-kafka/pkg/client/clientset/versioned"
	bindingsv1alpha1 "knative.dev/eventing-kafka/pkg/client/clientset/versioned/typed/bindings/v1alpha1"
	fakebindingsv1alpha1 "knative.dev/eventing-kafka/pkg/client/clientset/versioned/typed/bindings/v1alpha1/fake"
	bindingsv1beta1 "knative.dev/eventing-kafka/pkg/client/clientset/versioned/typed/bindings/v1beta1"
	fakebindingsv1beta1 "knative.dev/eventing-kafka/pkg/client/clientset/versioned/typed/bindings/v1beta1/fake"
	kafkav1alpha1 "knative.dev/eventing-kafka/pkg/client/clientset/versioned/typed/kafka/v1alpha1"
	fakekafkav1alpha1 "knative.dev/eventing-kafka/pkg/client/clientset/versioned/typed/kafka/v1alpha1/fake"
	messagingv1beta1 "knative.dev/eventing-kafka/pkg/client/clientset/versioned/typed/messaging/v1beta1"
	fakemessagingv1beta1 "knative.dev/eventing-kafka/pkg/client/clientset/versioned/typed/messaging/v1beta1/fake"
	sourcesv1alpha1 "knative.dev/eventing-kafka/pkg/client/clientset/versioned/typed/sources/v1alpha1"
	fakesourcesv1alpha1 "knative.dev/eventing-kafka/pkg/client/clientset/versioned/typed/sources/v1alpha1/fake"
	sourcesv1beta1 "knative.dev/eventing-kafka/pkg/client/clientset/versioned/typed/sources/v1beta1"
	fakesourcesv1beta1 "knative.dev/eventing-kafka/pkg/client/clientset/versioned/typed/sources/v1beta1/fake"
)

// NewSimpleClientset returns a clientset that will respond with the provided objects.
// It's backed by a very simple object tracker that processes creates, updates and deletions as-is,
// without applying any validations and/or defaults. It shouldn't be considered a replacement
// for a real clientset and is mostly useful in simple unit tests.
func NewSimpleClientset(objects ...runtime.Object) *Clientset {
	o := testing.NewObjectTracker(scheme, codecs.UniversalDecoder())
	for _, obj := range objects {
		if err := o.Add(obj); err != nil {
			panic(err)
		}
	}

	cs := &Clientset{tracker: o}
	cs.discovery = &fakediscovery.FakeDiscovery{Fake: &cs.Fake}
	cs.AddReactor("*", "*", testing.ObjectReaction(o))
	cs.AddWatchReactor("*", func(action testing.Action) (handled bool, ret watch.Interface, err error) {
		gvr := action.GetResource()
		ns := action.GetNamespace()
		watch, err := o.Watch(gvr, ns)
		if err != nil {
			return false, nil, err
		}
		return true, watch, nil
	})

	return cs
}

// Clientset implements clientset.Interface. Meant to be embedded into a
// struct to get a default implementation. This makes faking out just the method
// you want to test easier.
type Clientset struct {
	testing.Fake
	discovery *fakediscovery.FakeDiscovery
	tracker   testing.ObjectTracker
}

func (c *Clientset) Discovery() discovery.DiscoveryInterface {
	return c.discovery
}

func (c *Clientset) Tracker() testing.ObjectTracker {
	return c.tracker
}

var _ clientset.Interface = &Clientset{}

// BindingsV1alpha1 retrieves the BindingsV1alpha1Client
func (c *Clientset) BindingsV1alpha1() bindingsv1alpha1.BindingsV1alpha1Interface {
	return &fakebindingsv1alpha1.FakeBindingsV1alpha1{Fake: &c.Fake}
}

// BindingsV1beta1 retrieves the BindingsV1beta1Client
func (c *Clientset) BindingsV1beta1() bindingsv1beta1.BindingsV1beta1Interface {
	return &fakebindingsv1beta1.FakeBindingsV1beta1{Fake: &c.Fake}
}

// KafkaV1alpha1 retrieves the KafkaV1alpha1Client
func (c *Clientset) KafkaV1alpha1() kafkav1alpha1.KafkaV1alpha1Interface {
	return &fakekafkav1alpha1.FakeKafkaV1alpha1{Fake: &c.Fake}
}

// MessagingV1beta1 retrieves the MessagingV1beta1Client
func (c *Clientset) MessagingV1beta1() messagingv1beta1.MessagingV1beta1Interface {
	return &fakemessagingv1beta1.FakeMessagingV1beta1{Fake: &c.Fake}
}

// SourcesV1alpha1 retrieves the SourcesV1alpha1Client
func (c *Clientset) SourcesV1alpha1() sourcesv1alpha1.SourcesV1alpha1Interface {
	return &fakesourcesv1alpha1.FakeSourcesV1alpha1{Fake: &c.Fake}
}

// SourcesV1beta1 retrieves the SourcesV1beta1Client
func (c *Clientset) SourcesV1beta1() sourcesv1beta1.SourcesV1beta1Interface {
	return &fakesourcesv1beta1.FakeSourcesV1beta1{Fake: &c.Fake}
}
