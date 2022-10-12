# Benchmark grpc forker

```shell
ghz -c 500 -n 200000 --insecure --call grpc.health.v1.Health.Check 0.0.0.0:9090
```

### with forker

```shell
Summary:
  Count:	200000
  Total:	14.38 s
  Slowest:	142.50 ms
  Fastest:	0.14 ms
  Average:	21.14 ms
  Requests/sec:	13910.00

Response time histogram:
  0.136   [1]     |
  14.372  [61253] |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  28.609  [92588] |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  42.845  [36131] |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  57.081  [7548]  |∎∎∎
  71.318  [1670]  |∎
  85.554  [599]   |
  99.790  [100]   |
  114.026 [39]    |
  128.263 [12]    |
  142.499 [59]    |

Latency distribution:
  10 % in 6.91 ms
  25 % in 12.65 ms
  50 % in 19.31 ms
  75 % in 27.64 ms
  90 % in 36.87 ms
  95 % in 42.87 ms
  99 % in 59.01 ms
```
