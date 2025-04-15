#!/bin/bash

set -e

mkdir -p ./tmp

TIMESTAMP=$(date "+%Y-%m-%d_%H-%M-%S")

STATIC_LOG="./tmp/statictest-$TIMESTAMP.log"
GOPHERMART_LOG="./tmp/gophermarttest-$TIMESTAMP.log"

: > "$STATIC_LOG"
: > "$GOPHERMART_LOG"

# Build
echo "🚀 Build project..." | tee -a "$GOPHERMART_LOG"
go build -o ./cmd/gophermart/gophermart ./cmd/gophermart/*.go 2>&1 | tee -a "$GOPHERMART_LOG"

echo "🔍 Start go vet с statictest..." | tee -a "$STATIC_LOG"
set +e
go vet -vettool=$(which statictest) ./... > "$STATIC_LOG" 2>&1
VET_EXIT_CODE=$?
set -e

# Stop if static error
if [[ -s "$STATIC_LOG" || $VET_EXIT_CODE -ne 0 ]]; then
    echo "❌ Статический анализ не прошел! Тесты не будут запущены."
    echo "📄 См. лог: $STATIC_LOG"
    exit 1
fi

# 🧪 Запуск тестов
echo "🧪 Start praktikum tests..." | tee -a "$GOPHERMART_LOG"
gophermarttest \
  -test.v -test.run=^TestGophermart$ \
  -gophermart-binary-path=cmd/gophermart/gophermart \
  -gophermart-host=localhost \
  -gophermart-port=8080 \
  -gophermart-database-uri="postgresql://postgres:postgres@localhost:5432/gophermart" \
  -accrual-binary-path=cmd/accrual/accrual_darwin_arm64 \
  -accrual-host=localhost \
  -accrual-port=$(random unused-port) \
  -accrual-database-uri="postgresql://postgres:postgres@localhost:5432/gophermart" 2>&1 | tee -a "$GOPHERMART_LOG"

# ✅ Вывод путей к логам
echo "📁 Logs saved:"
echo "   Statictest logs:        $STATIC_LOG"
echo "   Praktikum tests logs:   $GOPHERMART_LOG"
