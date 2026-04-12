---
mode: agent
description: Построить графики для результатов бенчмарков хэшей
---

Построй 3 графика зависимости времени выполнения от количества узлов — для HRW, WRH и Consistent.

Данные бери из файлов:
- `hashes/bench_results_100k.txt`
- `hashes/bench_results_1M.txt`
- `hashes/bench_results_10M.txt`

Обнови файл `hashes/plot/plot.go` — поменяй значения `x` и `y` на соответствующие данные из бенча.

После обновления запусти:
```
make hashes_plot
```

