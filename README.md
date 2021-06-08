## Test
`make test`

## See coverage
`make see-coverage`

## To be improved
- 加入真實的 DB 實作，並使用 `sqlmock` 來針對一班的情境添加測試
- cache layer 重複的部分太多，應進一步抽象化
- cache layer 等 async job 應該用 retry/timeout 的方式來測，測試結果應該會比直接用 `time.Sleep` 來等還要可預測
- cache layer 的 async update cache 會有 cache stampede 的問題，需再思考如何實作