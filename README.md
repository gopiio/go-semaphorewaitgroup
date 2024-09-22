# Go Semaphore Waitgroup

Add sempahore to Go Standard Libray Waitgroup to control the concurrency

```bash
go get github.com/gopiio/go-semaphorewaitgroup
```

## How to use

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

maxConcurrent := int64(2)
swg := NewSemaphoreWaitGroup(ctx, maxConcurrent)
```
