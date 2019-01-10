if `go version | cut -d ' ' -f 3 | cut -d '.' -f 1,2` != "go1.10" ]; then
    cd ./exts/callerstack && go mod tidy && cd ../../
fi
