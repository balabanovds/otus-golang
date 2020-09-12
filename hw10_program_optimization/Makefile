test_file_name = "hw10_program_optimization.test"
cpu_pprof_name = "cpu00.pprof"
mem_pprof_name = "mem00.pprof"

bench:
	go test -bench . -benchmem -benchtime=10x -run ^$$

profile:
	go test -bench ^BenchmarkGetDomainStat$$ -benchmem \
		-benchtime=10x -run ^$$ -cpuprofile $(cpu_pprof_name) -memprofile $(mem_pprof_name)

pprof-cpu: profile
	go tool pprof $(test_file_name) $(cpu_pprof_name)


pprof-mem: profile
	go tool pprof $(test_file_name) $(mem_pprof_name)

test: lint
	go test -v -count=1 -race -gcflags=-l -timeout=30s ./...

test-performance: lint
	go test -v -count=1 -timeout=30s -tags bench .

lint:
	golangci-lint run ./...
