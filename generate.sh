#!/bin/bash
protoc proto/chat.proto --go_out=plugins=grpc:.
