#!/bin/bash

echo "Loading data for testing..."
./skeleton_cli -db test.db SET test test0
./skeleton_cli -db test.db SET test test1
./skeleton_cli -db test.db SET test test2
./skeleton_cli -db test.db SET test test3
./skeleton_cli -db test.db SET test test4
./skeleton_cli -db test.db SET test test5
echo ""

echo "Key:"
./skeleton_cli -db test.db GET test
echo ""

echo "Value:"
./skeleton_cli -db test.db GET test VALUE
echo ""

echo "Snapshots:"
./skeleton_cli -db test.db GET test SNAPSHOTS
echo ""
