# Benchmark echo forker

```shell
ab -n 200000 -c 500 http://localhost:8080/
```

### with forker

```shell
Server Software:        
Server Hostname:        localhost
Server Port:            8080

Document Path:          /
Document Length:        8 bytes

Concurrency Level:      500
Time taken for tests:   31.660 seconds
Complete requests:      200000
Failed requests:        0
Total transferred:      24800000 bytes
HTML transferred:       1600000 bytes
Requests per second:    6317.10 [#/sec] (mean)
Time per request:       79.150 [ms] (mean)
Time per request:       0.158 [ms] (mean, across all concurrent requests)
Transfer rate:          764.96 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0   39   5.5     39      60
Processing:     8   40   6.1     40      82
Waiting:        0   29   6.3     29      64
Total:         44   79   3.0     79     136

Percentage of the requests served within a certain time (ms)
  50%     79
  66%     79
  75%     79
  80%     80
  90%     81
  95%     83
  98%     88
  99%     92
 100%    136 (longest request)
```

### without forker

```shell
Server Software:        
Server Hostname:        localhost
Server Port:            8080

Document Path:          /
Document Length:        8 bytes

Concurrency Level:      500
Time taken for tests:   34.230 seconds
Complete requests:      200000
Failed requests:        0
Total transferred:      24800000 bytes
HTML transferred:       1600000 bytes
Requests per second:    5842.85 [#/sec] (mean)
Time per request:       85.575 [ms] (mean)
Time per request:       0.171 [ms] (mean, across all concurrent requests)
Transfer rate:          707.53 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0   39   5.8     39      64
Processing:     9   46   8.4     45      98
Waiting:        0   33   8.9     31      81
Total:         46   85   6.9     84     153

Percentage of the requests served within a certain time (ms)
  50%     84
  66%     86
  75%     88
  80%     90
  90%     95
  95%     99
  98%    104
  99%    108
 100%    153 (longest request)

```