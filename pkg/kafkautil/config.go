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

package kafkautil

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"strconv"
	"strings"

	v1alpha1 "github.com/banzaicloud/kafka-operator/api/v1alpha1"
	"github.com/banzaicloud/kafka-operator/pkg/resources/kafka"
	"github.com/banzaicloud/kafka-operator/pkg/resources/pki"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Should I retain the option to run it as a standalone topic/user controller?
const (
	kafkaHostVar      = "KAFKA_BROKER"
	kafkaUseSSLVar    = "KAFKA_USE_SSL"
	kafkaSSLKeyVar    = "KAFKA_SSL_KEY_FILE"
	kafkaSSLCertVar   = "KAFKA_SSL_CERT_FILE"
	kafkaSSLCAVar     = "KAFKA_SSL_CA_FILE"
	kafkaSSLVerifyVar = "KAFKA_INSECURE_SKIP_VERIFY"
	kafkaTimeoutVar   = "KAFKA_OPERATION_TIMEOUT_SECONDS"

	kafkaCAVar     = "KAFKA_ISSUER_CA_NAME"
	kafkaCATypeVar = "KAFKA_ISSUER_CA_KIND"

	kafkaDefaultTimeout = int64(10)
)

// KafkaConfig are the options to creating a new ClusterAdmin client
type KafkaConfig struct {
	BrokerURI             string
	UseSSL                bool
	TLSConfig             *tls.Config
	SSLKeyFile            string
	SSLCertFile           string
	SSLCAFile             string
	SSLInsecureSkipVerify bool

	IssueCA     string
	IssueCAKind string

	OperationTimeout int64
}

// EnvConfig is from when this was used in a standalone controller
func EnvConfig() *KafkaConfig {
	return &KafkaConfig{
		BrokerURI:             os.Getenv(kafkaHostVar),
		UseSSL:                parseBool(os.Getenv(kafkaUseSSLVar)),
		SSLKeyFile:            os.Getenv(kafkaSSLKeyVar),
		SSLCertFile:           os.Getenv(kafkaSSLCertVar),
		SSLCAFile:             os.Getenv(kafkaSSLCAVar),
		SSLInsecureSkipVerify: parseBool(os.Getenv(kafkaSSLVerifyVar)),
		IssueCA:               os.Getenv(kafkaCAVar),
		IssueCAKind:           os.Getenv(kafkaCATypeVar),
		OperationTimeout:      getOperationTimeout(),
	}
}

// ClusterConfig creates connection options from a KafkaCluster CR
func ClusterConfig(client client.Client, cluster *v1alpha1.KafkaCluster) (*KafkaConfig, error) {
	conf := &KafkaConfig{}
	conf.BrokerURI = generateKafkaAddress(cluster)
	conf.OperationTimeout = kafkaDefaultTimeout

	if cluster.Spec.ListenersConfig.SSLSecrets != nil {
		var err error
		tlsKeys := &corev1.Secret{}
		err = client.Get(context.TODO(), types.NamespacedName{Namespace: cluster.Namespace, Name: cluster.Spec.ListenersConfig.SSLSecrets.TLSSecretName}, tlsKeys)
		if err != nil {
			return conf, err
		}
		clientCert := tlsKeys.Data["clientCert"]
		clientKey := tlsKeys.Data["clientKey"]
		caCert := tlsKeys.Data["caCert"]
		x509ClientCert, err := tls.X509KeyPair(clientCert, clientKey)
		if err != nil {
			return conf, err
		}

		rootCAs := x509.NewCertPool()
		rootCAs.AppendCertsFromPEM(caCert)
		t := &tls.Config{
			Certificates: []tls.Certificate{x509ClientCert},
			RootCAs:      rootCAs,
		}

		conf.UseSSL = true
		conf.TLSConfig = t
		conf.IssueCA = fmt.Sprintf(pki.BrokerIssuerTemplate, cluster.Name)
		conf.IssueCAKind = "ClusterIssuer"
	}

	return conf, nil
}

func getOperationTimeout() int64 {
	var timeout int64
	var err error
	reqTimeout := os.Getenv(kafkaTimeoutVar)
	if reqTimeout == "" {
		log.Info(fmt.Sprint(kafkaTimeoutVar, " is not set. Assuming defaults."))
		timeout = kafkaDefaultTimeout
	} else if timeout, err = strconv.ParseInt(reqTimeout, 10, 64); err != nil {
		log.Info(fmt.Sprint(reqTimeout, " is not a valid integer for ", kafkaTimeoutVar, " - using default ", kafkaDefaultTimeout))
		timeout = kafkaDefaultTimeout
	}
	return timeout
}

// parseBool is a no errors ParseBool
func parseBool(str string) bool {
	if strings.ToLower(str) == "true" {
		return true
	}
	// if empty or anything else return false
	return false
}

func generateKafkaAddress(cluster *v1alpha1.KafkaCluster) string {
	if cluster.Spec.HeadlessServiceEnabled {
		return fmt.Sprintf("%s.%s:%d", fmt.Sprintf(kafka.HeadlessServiceTemplate, cluster.Name), cluster.Namespace, cluster.Spec.ListenersConfig.InternalListeners[0].ContainerPort)
	}
	return fmt.Sprintf("%s.%s.svc.cluster.local:%d", fmt.Sprintf(kafka.AllBrokerServiceTemplate, cluster.Name), cluster.Namespace, cluster.Spec.ListenersConfig.InternalListeners[0].ContainerPort)
}
