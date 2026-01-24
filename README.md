# Техническое задание: Windows System Metrics Exporter

### 1. Цель
Создание легковесного агента на Go для сбора системных метрик ОС Windows и их экспорта в формате, совместимом с Prometheus.

### 2. Требования к реализации
* **Периодичность:** Сбор данных строго по `time.Ticker`.
* **Параллелизм:** Каждая категория метрик (CPU, Memory, Disk) должна опрашиваться в отдельной горутине с использованием `sync.WaitGroup` для контроля завершения цикла опроса.
* **Формат:** Данные должны отдаваться через HTTP-эндпоинт `/metrics`.
* **Завершение:** Реализация Graceful Shutdown при получении сигналов `os.Interrupt` или `syscall.SIGTERM`.

### 3. Стек и зависимости
* **Язык:** Go 1.21+
* **Сбор метрик (Windows API wrapper):** [gopsutil/v3](https://github.com/shirou/gopsutil)
* **Prometheus SDK:** [client_golang](https://github.com/prometheus/client_golang)
* **Низкоуровневый доступ к ОС:** [golang.org/x/sys/windows](https://pkg.go.dev/golang.org/x/sys/windows)

### 4. Список метрик (минимальный набор)
| Название метрики | Тип Prometheus | Источник данных |
| :--- | :--- | :--- |
| `win_cpu_usage_percent` | Gauge | `cpu.Percent` |
| `win_mem_usage_bytes` | Gauge | `mem.VirtualMemory` |
| `win_disk_free_bytes` | Gauge | `disk.Usage("C:")` |
| `win_net_bytes_sent_total` | Counter | `net.IOCounters` |

### 5. Ожидаемый результат
Исполняемый файл `.exe`, который при запуске открывает порт (например, `:8080`) и отображает актуальные системные показатели в текстовом виде при обращении к `/metrics`.