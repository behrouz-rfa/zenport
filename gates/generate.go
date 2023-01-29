package gates

//go:generate buf generate
//go:generate mockery --quiet --dir ./gatespb -r --all --inpackage --case underscore
//go:generate mockery --quiet --dir ./internal -r --all --inpackage --case underscore

//go:generate swagger generate client -q -f ./internal/rest/api.swagger.json -c gatesclient -m gatesclient/models --with-flatten=remove-unused
