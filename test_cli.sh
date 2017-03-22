#!/bin/bash

echo "Loading data for testing..."
./skeleton-cli -db test.db SET test test0
./skeleton-cli -db test.db SET test test1
./skeleton-cli -db test.db SET test test2
./skeleton-cli -db test.db SET test test3
./skeleton-cli -db test.db SET test test4
./skeleton-cli -db test.db SET test test5
echo ""

echo "Key:"
./skeleton-cli -db test.db GET test
echo ""

echo "Value:"
./skeleton-cli -db test.db GET test VALUE
echo ""

echo "Snapshots:"
./skeleton-cli -db test.db GET test SNAPSHOTS
echo ""
