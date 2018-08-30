# Copyright © 2018 by PACE Telematics GmbH. All rights reserved.
# Created at 2018/08/24 by Vincent Landgraf

JSONAPITEST="http/jsonapi/generator/internal"

jsonapi:
	pace service generate rest --pkg poi \
		--path $(JSONAPITEST)/poi/open-api_test.go \
		--source $(JSONAPITEST)/poi/open-api.json
	pace service generate rest --pkg fueling \
		--path $(JSONAPITEST)/fueling/open-api_test.go \
		--source $(JSONAPITEST)/fueling/open-api.json
	pace service generate rest --pkg pay \
		--path $(JSONAPITEST)/pay/open-api_test.go \
		--source $(JSONAPITEST)/pay/open-api.json
