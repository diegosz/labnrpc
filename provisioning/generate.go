package provisioning

//go:generate buf generate

//go:generate mockery --quiet --dir ./provisioningpb -r --all --inpackage --case underscore
//go:generate mockery --quiet --dir ./internal -r --all --inpackage --case underscore

//go:generate swagger generate client -q -f ./internal/rest/api.swagger.json -c provisioningclient -m provisioningclient/models --with-flatten=remove-unused
