### Introduction

The notifications-sendgrid-receiver is designed to parse posts from sendgrid and convert them to a form that the [notifications service](https://github.com/cloudfoundry-incubator/notifications) can understand.  It currently expects an multipart post that has several parts as described in the [sendgrid documentation](https://sendgrid.com/docs/API_Reference/Webhooks/parse.html)

That's about it.  Its written in go.  There are convenience scripts to deploy to bosh-lite in the bin directory.


### Getting started

This guide assumes you have bosh-lite set up on your workstation. If not, [read this](https://github.com/cloudfoundry/bosh-lite)

Clone the app into your go workspace. Before you can boot the app you need to configure a UAA user so you can properly set the UAA_CLIENT_ID and UAA_CLIENT_SECRET fields of the `./bin/env/development` config file.


The app can be booted by running:
	
	bin/run
	
The test suite is run by:

	bin/unit
