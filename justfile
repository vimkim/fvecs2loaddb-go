test:
    go run . -fvecs data/vector_256dim_3row_seed0.fvecs -table tbl_test -out out/loaddb_test.object

gen-256-75000:
    go run . -fvecs data/vector_256dim_75000row_seed0.fvecs -table tbl_256_75000 -out out/loaddb_256_75000.object

gen-256-150000:
    go run . -fvecs data/vector_256dim_150000row_seed0.fvecs -table tbl_256_150000 -out out/loaddb_256_150000.object

gen-256-300000:
    go run . -fvecs data/vector_256dim_300000row_seed0.fvecs -table tbl_256_300000 -out out/loaddb_256_300000.object

gen-768-all: gen-768-75000 gen-768-150000 gen-768-300000

gen-768-75000:
    go run . -fvecs data/vector_768dim_75000row_seed0.fvecs -table tbl_768_75000 -out out/loaddb_768_75000.object

gen-768-150000:
    go run . -fvecs data/vector_768dim_150000row_seed0.fvecs -table tbl_768_150000 -out out/loaddb_768_150000.object

gen-768-300000:
    go run . -fvecs data/vector_768dim_300000row_seed0.fvecs -table tbl_768_300000 -out out/loaddb_768_300000.object

gen-1536-75000:
    go run . -fvecs data/vector_1536dim_75000row_seed0.fvecs -table tbl_1536_75000 -out out/loaddb_1536_75000.object

gen-1536-150000:
    go run . -fvecs data/vector_1536dim_150000row_seed0.fvecs -table tbl_1536_150000 -out out/loaddb_1536_150000.object

gen-1536-300000:
    go run . -fvecs data/vector_1536dim_300000row_seed0.fvecs -table tbl_1536_300000 -out out/loaddb_1536_300000.object
