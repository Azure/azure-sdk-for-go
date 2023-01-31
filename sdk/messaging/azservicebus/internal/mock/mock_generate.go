// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:generate mockgen -source ../namespace.go -package mock -copyright_file ./testdata/copyright.txt -destination mock_namespace.go NamespaceWithNewAMQPLinks,NamespaceForAMQPLinks

//go:generate mockgen -source ../amqpwrap/amqpwrap.go -package mock -copyright_file ./testdata/copyright.txt -destination mock_amqp.go

//go:generate mockgen -source ../amqpwrap/rpc.go -package mock -copyright_file ./testdata/copyright.txt -destination mock_rpc.go

package mock
