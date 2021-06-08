## Description
- 由於比較不熟且重點練習的部分在快取機制，故先使用滿足 interface 的 `fakedb` 實作來滿足服務運行時取得的資料

## Test
`make test`

## See coverage
`make see-coverage`

## To be improved
- cache layer 重複的部分太多，應進一步抽象化
- cache layer 等 async job 應該用 retry/timeout 的方式來測，測試結果應該會比直接用 `time.Sleep` 來等還要可預測
- 在有快取的情況下、且 cache invalid 時 call async job，也可能會有 cache stampede 的問題，此部分也應再進一步避免
- 加入真實的 DB 實作，並使用 `sqlmock` 來針對一般的情境添加測試

## Reference
- [Multiple Lock Based on Input](https://medium.com/@kf99916/multiple-lock-based-on-input-in-golang-74931a3c8230)