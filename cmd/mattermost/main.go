// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package main

import (
	"os"

	"github.com/mattermost/mattermost-server/v5/cmd/mattermost/commands"

	// Plugins
	"github.com/mattermost/mattermost-server/v5/model/gitlab"

	// Enterprise Imports
	"github.com/mattermost/mattermost-server/v5/imports"

	// Enterprise Deps
	"github.com/gorilla/handlers"
	"github.com/hako/durafmt"
	"github.com/hashicorp/memberlist"
	"github.com/mattermost/gosaml2"
	"github.com/mattermost/ldap"
	"github.com/mattermost/rsc/qr"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tylerb/graceful"
	"gopkg.in/olivere/elastic.v6"
)

func main() {
	if err := commands.Run(os.Args[1:]); err != nil {
		os.Exit(1)
	}
}
