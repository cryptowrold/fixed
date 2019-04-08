**Summary**

A fixed place numeric library with overflow check in Go designed for performance.

All numbers have a fixed 8 decimal places, and the maximum permitted value is + 9999999999,
or just under 10 billion.

The library is safe for concurrent use. It has built-in support for binary and json marshalling.

It is ideally suited for high performance trading financial systems. All common math operations are completed with 0 allocs.

**Performance**

<pre>
goos: darwin
goarch: amd64
pkg: github.com/cryptowrold/fixed
BenchmarkAddFixed-4         	2000000000	         0.93 ns/op	       0 B/op	       0 allocs/op
BenchmarkAddDecimal-4       	 5000000	       338 ns/op	     176 B/op	       8 allocs/op
BenchmarkAddBigInt-4        	100000000	        20.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkAddBigFloat-4      	10000000	       114 ns/op	      48 B/op	       1 allocs/op
BenchmarkMulFixed-4         	200000000	         6.66 ns/op	       0 B/op	       0 allocs/op
BenchmarkMulDecimal-4       	10000000	       105 ns/op	      80 B/op	       2 allocs/op
BenchmarkMulBigInt-4        	100000000	        24.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkMulBigFloat-4      	30000000	        52.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkDivFixed-4         	200000000	         6.76 ns/op	       0 B/op	       0 allocs/op
BenchmarkDivDecimal-4       	 1000000	      1105 ns/op	     568 B/op	      21 allocs/op
BenchmarkDivBigInt-4        	20000000	        63.8 ns/op	       8 B/op	       1 allocs/op
BenchmarkDivBigFloat-4      	10000000	       153 ns/op	      24 B/op	       2 allocs/op
BenchmarkCmpFixed-4         	2000000000	         0.51 ns/op	       0 B/op	       0 allocs/op
BenchmarkCmpDecimal-4       	100000000	        11.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkCmpBigInt-4        	200000000	         7.68 ns/op	       0 B/op	       0 allocs/op
BenchmarkCmpBigFloat-4      	200000000	         7.24 ns/op	       0 B/op	       0 allocs/op
BenchmarkStringFixed-4      	20000000	        71.9 ns/op	      32 B/op	       1 allocs/op
BenchmarkStringNFixed-4     	20000000	        72.6 ns/op	      32 B/op	       1 allocs/op
BenchmarkStringDecimal-4    	 5000000	       308 ns/op	      64 B/op	       5 allocs/op
BenchmarkStringBigInt-4     	10000000	       171 ns/op	      24 B/op	       2 allocs/op
BenchmarkStringBigFloat-4   	 3000000	       571 ns/op	     192 B/op	       8 allocs/op
BenchmarkWriteTo-4          	30000000	        52.4 ns/op	      18 B/op	       0 allocs/op
</pre>

The "decimal" above is the common [shopspring decimal](https://github.com/shopspring/decimal) library