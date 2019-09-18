```
make start
```

```
export FN_API_URL=http://caddy:2020
```

```
fn ls app
fn create app pipeline-demo
```

```
cd simple-random/
fn deploy --app pipeline-demo --local
fn invoke pipeline-demo simple-random
```

```
fn config function pipeline-demo simple-random "times" "10"
fn inspect function pipeline-demo simple-random
echo '{"times": 7}' | fn invoke pipeline-demo simple-random --content-type application/json
```

This is sequential A -> B -> C

----------

```
cd polishcalc-solver/
fn deploy --app pipeline-demo --local

echo -n '{"equation": "- * / 15 - 7 + 1 1 3 + 2 + 1 5"}' | fn invoke pipeline-demo polishcalc-solver --content-type application/json
```

```
cd polishcalc-accum/
fn deploy --app pipeline-demo --local

echo -n '{"equations": ["239", "* 2 + 9 512", "+ 3 2", "- * / 15 - 7 + 1 1 3 + 2 + 1 5"]}' | fn invoke pipeline-demo polishcalc-accum --content-type application/json
```

This is parallel, A -> (B, C, D)