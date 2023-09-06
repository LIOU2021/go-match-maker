# Introduction
这是一个在模拟游戏配对撮合或是交友软体配对的一个project
- demo
    ```go
    go run ./cmd/main.go
    // or
    go run ./_example/main.go
    ```
- 當用戶加入時的roomId屬於當前hub中沒有的，那麼hub會自動增加該room到hub中
# todo
- 新增初始化檢查緩存有無資料的bug
- 測試腳本執行buff比較小的狀態
- config新增
    - 指定房間要取幾個
    - 房間內的成員要取幾個