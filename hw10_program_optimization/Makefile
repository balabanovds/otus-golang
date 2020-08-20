bench:
	go test -bench . -benchmem -benchtime=10x -run ^$$

profile:
	go test -bench ^BenchmarkGetDomainStat$$ -benchmem \
		-benchtime=10x -run ^$$ -cpuprofile cpu00.pprof -memprofile mem00.pprof


profile-mem:
	go test -bench ^BenchmarkGetDomainStat$$ -benchmem -benchtime=10x -run ^$$ -memprofile mem00.pprof
	go tool pprof -alloc_space mem00.pprof
