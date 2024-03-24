# marketplace-10k-rps

Make new migration script

```migrate create -ext sql -dir migration -seq init```

Migrate database

up local : ```../migrate -database postgres://postgres:root@localhost:5432/socmed?sslmode=disable -path ./db/migrations up```

down : ```migrate -database postgres://postgres:root@localhost:5432/marketplace?sslmode=disable -path ./db/migrations down```

# K6 LOAD TEST RESULT

```
 50 vus
     http_req_duration..............: avg=52.57ms  min=523.81µs med=11.45ms max=1.1s    p(90)=99.08ms  p(95)=242.08ms
       { expected_response:true }...: avg=84.9ms   min=3.37ms   med=24.57ms max=1.1s    p(90)=191.81ms p(95)=391.93ms
checks.........................: 90.64% ✓ 2074       ✗ 214 

100 vus
     http_req_duration..............: avg=75.11ms  min=532.28µs med=21.13ms max=1.56s   p(90)=155.4ms  p(95)=466.75ms
       { expected_response:true }...: avg=120.59ms min=3.62ms   med=49.09ms max=1.56s   p(90)=271.51ms p(95)=681.19ms
checks.........................: 91.49% ✓ 4161      ✗ 387

200 vus
     http_req_duration..............: avg=198.84ms min=490.39µs med=40.62ms max=5.99s    p(90)=402.82ms p(95)=944.68ms
       { expected_response:true }...: avg=307.84ms min=2.93ms   med=77.43ms max=5.99s    p(90)=876.41ms p(95)=1.51s
    checks.........................: 91.11% ✓ 8758       ✗ 854

300 vus
     checks.........................: 93.09% ✓ 11806      ✗ 875
     data_received..................: 3.3 MB 146 kB/s
     data_sent......................: 4.4 MB 198 kB/s
     http_req_blocked...............: avg=3.7ms    min=3.18µs   med=3.82µs   max=215.26ms p(90)=5.64µs   p(95)=27.37µs
     http_req_connecting............: avg=3.62ms   min=0s       med=0s       max=207.1ms  p(90)=0s       p(95)=0s
     http_req_duration..............: avg=332.05ms min=505.24µs med=41.53ms  max=10.43s   p(90)=934.38ms p(95)=1.74s
       { expected_response:true }...: avg=529.7ms  min=3.09ms   med=196.38ms max=8.88s    p(90)=1.7s     p(95)=2.17s

600 vus
checks.........................: 92.61% ✓ 23516      ✗ 1874
     http_req_duration..............: avg=732.39ms min=572.97µs med=156.37ms max=11.97s   p(90)=2.05s   p(95)=3.47s
       { expected_response:true }...: avg=1.23s    min=3.27ms   med=786.42ms max=9.06s    p(90)=2.94s   p(95)=4.19s
```